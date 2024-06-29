package repository

import (
	"context"
	"fmt"
	"net/http"
	"oosa_rewild/internal/config"
	"oosa_rewild/internal/helpers"
	"oosa_rewild/internal/models"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RewildingRegisterRepository struct{}
type RewildingRegisterRequest struct {
	RewildingType                 string  `json:"rewilding_type"`
	RewildingApplyOfficial        bool    `json:"rewilding_apply_official"`
	RewildingReferenceInformation string  `json:"rewilding_reference_information"`
	RewildingPocketList           string  `json:"rewilding_pocket_list"`
	RewildingName                 string  `json:"rewilding_name" validate:"required"`
	RewildingLat                  float64 `json:"rewilding_lat"  validate:"required"`
	RewildingLng                  float64 `json:"rewilding_lng"  validate:"required"`
}

func (r RewildingRegisterRepository) Retrieve(c *gin.Context) {
	userDetail := helpers.GetAuthUser(c)
	var results []models.Rewilding
	filter := bson.D{{Key: "rewilding_created_by", Value: userDetail.UsersId}}
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

func (r RewildingRegisterRepository) Read(c *gin.Context) {
	id, _ := primitive.ObjectIDFromHex(c.Param("id"))
	userDetail := helpers.GetAuthUser(c)
	var Rewilding models.Rewilding
	filter := bson.D{
		{Key: "_id", Value: id},
		{Key: "rewilding_created_by", Value: userDetail.UsersId},
	}
	err := config.DB.Collection("Rewilding").FindOne(context.TODO(), filter).Decode(&Rewilding)

	if err != nil {
		helpers.ResultEmpty(c, err)
		return
	}

	c.JSON(200, Rewilding)
}

func (r RewildingRegisterRepository) Create(c *gin.Context) {
	userDetail := helpers.GetAuthUser(c)
	var payload RewildingRegisterRequest
	validateError := helpers.Validate(c, &payload)
	if validateError != nil {
		return
	}

	insert := models.Rewilding{
		RewildingCreatedBy: userDetail.UsersId,
		RewildingCreatedAt: primitive.NewDateTimeFromTime(time.Now()),
	}

	r.ProcessData(c, &insert, payload)

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

func (r RewildingRegisterRepository) Update(c *gin.Context) {
	var payload RewildingRegisterRequest
	validateError := helpers.Validate(c, &payload)
	if validateError != nil {
		return
	}

	id, _ := primitive.ObjectIDFromHex(c.Param("id"))
	userDetail := helpers.GetAuthUser(c)
	var Rewilding models.Rewilding

	filter := bson.D{
		{Key: "_id", Value: id},
		{Key: "rewilding_created_by", Value: userDetail.UsersId},
	}
	err := config.DB.Collection("Rewilding").FindOne(context.TODO(), filter).Decode(&Rewilding)

	if err != nil {
		helpers.ResultEmpty(c, err)
		return
	}

	r.ProcessData(c, &Rewilding, payload)

	upd := bson.D{{Key: "$set", Value: Rewilding}}
	config.DB.Collection("Rewilding").UpdateOne(context.TODO(), filter, upd)

	c.JSON(200, Rewilding)
}

func (r RewildingRegisterRepository) ProcessData(c *gin.Context, Rewilding *models.Rewilding, payload RewildingRegisterRequest) {
	lat, _ := primitive.ParseDecimal128(fmt.Sprint(payload.RewildingLat))
	lng, _ := primitive.ParseDecimal128(fmt.Sprint(payload.RewildingLng))

	rewildingApplyOfficial := false

	if payload.RewildingApplyOfficial {
		rewildingApplyOfficial = true
	}

	Rewilding.RewildingApplyOfficial = &rewildingApplyOfficial
	Rewilding.RewildingType = helpers.StringToPrimitiveObjId(payload.RewildingType)
	Rewilding.RewildingName = payload.RewildingName
	Rewilding.RewildingLat = lat
	Rewilding.RewildingLng = lng

}

func (r RewildingRegisterRepository) Options(c *gin.Context) {
	var RefRewildingTypes []models.RefRewildingTypes
	cursor, err := config.DB.Collection("RefRewildingTypes").Find(context.TODO(), bson.D{})
	if err != nil {
		return
	}
	cursor.All(context.TODO(), &RefRewildingTypes)
	c.JSON(http.StatusOK, gin.H{"rewilding_types": RefRewildingTypes})
}
