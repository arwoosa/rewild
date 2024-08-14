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
	"go.mongodb.org/mongo-driver/mongo"
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

type RewildingFormDataRequest struct {
	RewildingType                 string   `form:"rewilding_type"`
	RewildingApplyOfficial        bool     `form:"rewilding_apply_official"`
	RewildingReferenceInformation []string `form:"rewilding_reference_information"`
	RewildingPocketList           []string `form:"rewilding_pocket_list"`
	RewildingName                 string   `form:"rewilding_name" validate:"required"`
	RewildingLat                  float64  `form:"rewilding_lat"  validate:"required"`
	RewildingLng                  float64  `form:"rewilding_lng"  validate:"required"`
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

	agg := mongo.Pipeline{
		bson.D{{
			Key: "$lookup", Value: bson.M{
				"from":         "Users",
				"localField":   "rewilding_created_by",
				"foreignField": "_id",
				"as":           "rewilding_created_by_user",
			},
		}},
		bson.D{{
			Key: "$unwind", Value: bson.M{
				"path":                       "$rewilding_created_by_user",
				"preserveNullAndEmptyArrays": true,
			},
		}},
	}

	if owner != "" {
		middleware.CheckIfAuth(c)
		userDetail := helpers.GetAuthUser(c)
		agg = append(agg, bson.D{{
			Key: "$match", Value: bson.M{"rewilding_created_by": userDetail.UsersId},
		}})
		//filter["rewilding_created_by"] = userDetail.UsersId
	}

	cursor, err := config.DB.Collection("Rewilding").Aggregate(context.TODO(), agg)
	if err != nil {
		panic(err)
	}
	cursor.All(context.TODO(), &results)
	if len(results) == 0 {
		helpers.ResponseNoData(c, "No Data")
		return
	}

	for key, v := range results {
		if v.RewildingPhotos == nil {
			results[key].RewildingPhotos = make([]models.RewildingPhotos, 0)
		}
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

	if Rewilding.RewildingPhotos == nil {
		Rewilding.RewildingPhotos = make([]models.RewildingPhotos, 0)
	}

	c.JSON(200, Rewilding)
}

func (r RewildingRepository) Create(c *gin.Context) {
	userDetail := helpers.GetAuthUser(c)
	var payload RewildingFormDataRequest
	validateError := helpers.ValidateForm(c, &payload)
	if validateError != nil {
		return
	}

	lat := payload.RewildingLat
	lng := payload.RewildingLng

	geocode := helpers.GoogleMapsGeocode(c, payload.RewildingLat, payload.RewildingLng)
	elevation := helpers.GoogleMapsElevation(c, payload.RewildingLat, payload.RewildingLng)

	location := helpers.GooglePlacesToLocationArray(geocode.AddressComponents)
	area, _ := helpers.GooglePlacesGetArea(geocode.AddressComponents, "administrative_area_level_1")
	_, countryCode := helpers.GooglePlacesGetArea(geocode.AddressComponents, "country")

	rewildingApplyOfficial := false

	if payload.RewildingApplyOfficial {
		rewildingApplyOfficial = true
	}

	form, _ := c.MultipartForm()
	files := form.File["rewilding_photo[]"]
	var RewildingPhotos []models.RewildingPhotos

	for _, file := range files {
		_, validateErr := helpers.ValidatePhoto(c, file, true)

		if validateErr == nil {
			cloudflare := CloudflareRepository{}
			cloudflareResponse, postErr := cloudflare.Post(c, file)
			if postErr != nil {
				helpers.ResponseBadRequestError(c, postErr.Error())
				return
			}
			fileName := cloudflare.ImageDelivery(cloudflareResponse.Result.Id, "public")
			RwPhoto := models.RewildingPhotos{
				RewildingPhotosID:   primitive.NewObjectID(),
				RewildingPhotosPath: fileName,
			}
			RewildingPhotos = append(RewildingPhotos, RwPhoto)
		}
	}

	insert := models.Rewilding{
		RewildingApplyOfficial: &rewildingApplyOfficial,
		RewildingArea:          area,
		RewildingLocation:      location,
		RewildingCountryCode:   countryCode,
		RewildingName:          payload.RewildingName,
		RewildingLat:           lat,
		RewildingLng:           lng,
		RewildingElevation:     elevation.Elevation,
		RewildingCreatedBy:     userDetail.UsersId,
		RewildingCreatedAt:     primitive.NewDateTimeFromTime(time.Now()),
		RewildingPhotos:        RewildingPhotos,
	}

	var referenceLinks []models.RewildingReferenceLinks
	if len(payload.RewildingReferenceInformation) > 0 {
		for _, referenceInformation := range payload.RewildingReferenceInformation {
			insRefLink := models.RewildingReferenceLinks{
				RewildingReferenceLinksLink: referenceInformation,
			}
			referenceLinks = append(referenceLinks, insRefLink)
		}
		insert.RewildingReferenceLinks = referenceLinks
	}

	result, err := config.DB.Collection("Rewilding").InsertOne(context.TODO(), insert)
	if err != nil {
		fmt.Println("ERROR", err.Error())
		return
	}

	var Rewilding models.Rewilding
	config.DB.Collection("Rewilding").FindOne(context.TODO(), bson.D{{Key: "_id", Value: result.InsertedID}}).Decode(&Rewilding)

	// Add to pocket list
	if len(payload.RewildingPocketList) > 0 {
		for _, pocketListId := range payload.RewildingPocketList {
			insert := models.PocketListItems{
				PocketListItemsMst:       helpers.StringToPrimitiveObjId(pocketListId),
				PocketListItemsRewilding: Rewilding.RewildingID,
				PocketListItemsName:      Rewilding.RewildingName,
			}

			_, err := config.DB.Collection("PocketListItems").InsertOne(context.TODO(), insert)
			if err != nil {
				fmt.Println("ERROR", err.Error())
				return
			} else {
				PocketListRepository{}.UpdateCount(c, pocketListId)
			}
		}
	}

	c.JSON(http.StatusOK, Rewilding)
}

func (r RewildingRepository) Options(c *gin.Context) {
	RefRewildingTypes := helpers.RefRewildingTypes()
	c.JSON(http.StatusOK, gin.H{"rewilding_types": RefRewildingTypes})
}
