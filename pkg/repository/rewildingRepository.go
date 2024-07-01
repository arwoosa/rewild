package repository

import (
	"context"
	"fmt"
	"net/http"
	"oosa_rewild/internal/config"
	"oosa_rewild/internal/helpers"
	"oosa_rewild/internal/middleware"
	"oosa_rewild/internal/models"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RewildingRepository struct{}
type RewildingRequest struct {
	RewildingType                 string  `json:"rewilding_type"`
	RewildingApplyOfficial        bool    `json:"rewilding_apply_official"`
	RewildingReferenceInformation string  `json:"rewilding_reference_information"`
	RewildingPocketList           string  `json:"rewilding_pocket_list"`
	RewildingName                 string  `json:"rewilding_name" validate:"required"`
	RewildingLat                  float64 `json:"rewilding_lat"  validate:"required"`
	RewildingLng                  float64 `json:"rewilding_lng"  validate:"required"`
}

// @Summary Rewilding
// @Description Retrieve all rewilding
// @ID rewilding
// @Produce json
// @Tags Rewilding
// @Success 200 {object} []models.Rewilding
// @Failure 400 {object} structs.Message
// @Router /rewilding [get]
func (r RewildingRepository) Retrieve(c *gin.Context) {
	var results []models.Rewilding
	owner := c.Query("owner")

	filter := bson.M{}

	if owner != "" {
		middleware.CheckIfAuth(c)
		userDetail := helpers.GetAuthUser(c)
		filter["rewilding_created_by"] = userDetail.UsersId
	}
	cursor, err := config.DB.Collection("Rewilding").Find(context.TODO(), filter)
	if err != nil {
		panic(err)
	}
	cursor.All(context.TODO(), &results)
	if len(results) == 0 {
		helpers.ResponseNoData(c, "No Data")
		return
	}
	c.JSON(http.StatusOK, results)
}

func (r RewildingRepository) Read(c *gin.Context) {
	id, _ := primitive.ObjectIDFromHex(c.Param("id"))
	var Rewilding models.Rewilding
	err := config.DB.Collection("Rewilding").FindOne(context.TODO(), bson.D{{Key: "_id", Value: id}}).Decode(&Rewilding)

	if err != nil {
		helpers.ResultEmpty(c, err)
		return
	}

	c.JSON(200, Rewilding)
}

func (r RewildingRepository) Create(c *gin.Context) {
	userDetail := helpers.GetAuthUser(c)
	var payload RewildingRequest
	validateError := helpers.Validate(c, &payload)
	if validateError != nil {
		return
	}

	lat, _ := primitive.ParseDecimal128(fmt.Sprint(payload.RewildingLat))
	lng, _ := primitive.ParseDecimal128(fmt.Sprint(payload.RewildingLng))

	rewildingApplyOfficial := false

	if payload.RewildingApplyOfficial {
		rewildingApplyOfficial = true
	}

	insert := models.Rewilding{
		RewildingApplyOfficial: &rewildingApplyOfficial,
		RewildingType:          helpers.StringToPrimitiveObjId(payload.RewildingType),
		RewildingName:          payload.RewildingName,
		RewildingLat:           lat,
		RewildingLng:           lng,
		RewildingCreatedBy:     userDetail.UsersId,
		RewildingCreatedAt:     primitive.NewDateTimeFromTime(time.Now()),
	}

	result, err := config.DB.Collection("Rewilding").InsertOne(context.TODO(), insert)
	if err != nil {
		fmt.Println("ERROR", err.Error())
		return
	}

	var Rewilding models.Rewilding
	config.DB.Collection("Rewilding").FindOne(context.TODO(), bson.D{{Key: "_id", Value: result.InsertedID}}).Decode(&Rewilding)

	// Add to pocket list
	if payload.RewildingPocketList != "" {
		insert := models.PocketListItems{
			PocketListItemsMst:       helpers.StringToPrimitiveObjId(payload.RewildingPocketList),
			PocketListItemsRewilding: Rewilding.RewildingID,
			PocketListItemsName:      Rewilding.RewildingName,
		}

		_, err := config.DB.Collection("PocketListItems").InsertOne(context.TODO(), insert)
		if err != nil {
			fmt.Println("ERROR", err.Error())
			return
		}
		PocketListRepository{}.UpdateCount(c, payload.RewildingPocketList)
	}

	c.JSON(http.StatusOK, Rewilding)
}

func (r RewildingRepository) Options(c *gin.Context) {
	var RefRewildingTypes []models.RefRewildingTypes
	cursor, err := config.DB.Collection("RefRewildingTypes").Find(context.TODO(), bson.D{})
	if err != nil {
		return
	}
	cursor.All(context.TODO(), &RefRewildingTypes)
	c.JSON(http.StatusOK, gin.H{"rewilding_types": RefRewildingTypes})
}
