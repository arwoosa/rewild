package repository

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"oosa_rewild/internal/config"
	"oosa_rewild/internal/helpers"
	"oosa_rewild/internal/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/api/places/v1"
)

type RewildingSearchRepository struct{}
type RewildingSearchRequest struct {
	RewildingType string  `json:"rewilding_type"`
	RewildingCity string  `json:"rewilding_city"`
	RewildingArea string  `json:"rewilding_area"`
	RewildingName string  `json:"rewilding_name" validate:"required"`
	RewildingLat  float64 `json:"rewilding_lat"  validate:"required"`
	RewildingLng  float64 `json:"rewilding_lng"  validate:"required"`
}

// @Summary Rewilding
// @Description Retrieve all rewilding
// @ID rewilding
// @Produce json
// @Tags Rewilding
// @Success 200 {object} []models.Rewilding
// @Failure 400 {object} structs.Message
// @Router /rewilding [get]
func (r RewildingSearchRepository) Retrieve(c *gin.Context) {
	reqLat := c.Query("lat")
	reqLng := c.Query("lng")
	reqSearch := c.Query("search")

	if (reqLat == "" || reqLng == "") && reqSearch == "" {
		// Search by Lat and Lng
		c.JSON(http.StatusBadRequest, gin.H{"message": "Unable to do a search as latitude or longitude is not passed in"})
		return
	}

	lat, err := strconv.ParseFloat(reqLat, 64)
	if err != nil && reqLat != "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid latitude value"})
		return
	}

	lng, err := strconv.ParseFloat(reqLng, 64)
	if err != nil && reqLng != "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid longitude value"})
		return
	}

	placesService := helpers.GooglePlacesInitialise(c)
	if placesService == nil {
		log.Fatalf("error %s", err)
	}

	var includedTypes []string

	includedTypes = append(includedTypes, "hiking_area", "national_park")

	if reqSearch != "" {
		// Search by keyword
		req := places.GoogleMapsPlacesV1SearchTextRequest{
			TextQuery:      reqSearch,
			MaxResultCount: 10,
			LanguageCode:   "zh-TW",
		}

		placeReq := placesService.Places.SearchText(&req)
		placeReq.Header().Add("X-Goog-FieldMask", "places.id,places.types,places.displayName,places.formattedAddress,places.location,places.rating,places.userRatingCount")
		places, errPlace := placeReq.Do()
		if errPlace != nil {
			log.Fatalf("error %s", errPlace)
		}
		c.JSON(http.StatusOK, places)
	} else {
		req := places.GoogleMapsPlacesV1SearchNearbyRequest{
			IncludedTypes:  includedTypes,
			MaxResultCount: 10,
			LocationRestriction: &places.GoogleMapsPlacesV1SearchNearbyRequestLocationRestriction{
				Circle: &places.GoogleMapsPlacesV1Circle{
					Center: &places.GoogleTypeLatLng{
						Latitude:  lat,
						Longitude: lng,
					},
					Radius: 5000,
				},
			},
			LanguageCode: "zh-TW",
		}

		placeReq := placesService.Places.SearchNearby(&req)
		placeReq.Header().Add("X-Goog-FieldMask", "places.id,places.types,places.displayName,places.formattedAddress,places.location,places.rating,places.userRatingCount")
		places, errPlace := placeReq.Do()
		if errPlace != nil {
			log.Fatalf("error %s", errPlace)
		}
		c.JSON(http.StatusOK, places)
	}

	/*var results []models.Rewilding

	cursor, err := config.DB.Collection("Rewilding").Find(context.TODO(), bson.D{})
	if err != nil {
		panic(err)
	}
	cursor.All(context.TODO(), &results)
	c.JSON(200, results)*/
}

func (r RewildingSearchRepository) Create(c *gin.Context) {
	userDetail := helpers.GetAuthUser(c)
	var payload RewildingSearchRequest
	validateError := helpers.Validate(c, &payload)
	if validateError != nil {
		return
	}

	lat, _ := primitive.ParseDecimal128(fmt.Sprint(payload.RewildingLat))
	lng, _ := primitive.ParseDecimal128(fmt.Sprint(payload.RewildingLng))

	geocode := helpers.GoogleMapsGeocode(c, payload.RewildingLat, payload.RewildingLng)
	elevation := helpers.GoogleMapsElevation(c, payload.RewildingLat, payload.RewildingLng)

	area := helpers.GooglePlacesGetArea(geocode.AddressComponents, "administrative_area_level_1")

	/*Cache References*/
	var RewildingTypeData models.RefRewildingTypes
	typeId, _ := primitive.ObjectIDFromHex(payload.RewildingType)
	config.DB.Collection("RefRewildingTypes").FindOne(context.TODO(), bson.D{{Key: "_id", Value: typeId}}).Decode(&RewildingTypeData)

	insert := models.Rewilding{
		RewildingType:      helpers.StringToPrimitiveObjId(payload.RewildingType),
		RewildingTypeData:  RewildingTypeData,
		RewildingCity:      payload.RewildingCity,
		RewildingArea:      area,
		RewildingName:      payload.RewildingName,
		RewildingLat:       lat,
		RewildingLng:       lng,
		RewildingElevation: elevation.Elevation,
		RewildingCreatedBy: userDetail.UsersId,
		RewildingCreatedAt: primitive.NewDateTimeFromTime(time.Now()),
	}

	result, _ := config.DB.Collection("Rewilding").InsertOne(context.TODO(), insert)

	var Rewilding models.Rewilding
	err := config.DB.Collection("Rewilding").FindOne(context.TODO(), bson.D{{Key: "_id", Value: result.InsertedID}}).Decode(&Rewilding)

	if err != nil {
		helpers.ResultEmpty(c, err)
		return
	}

	c.JSON(200, Rewilding)
}

func (r RewildingSearchRepository) Read(c *gin.Context) {
	id, _ := primitive.ObjectIDFromHex(c.Param("id"))
	var Rewilding models.Rewilding
	err := config.DB.Collection("Rewilding").FindOne(context.TODO(), bson.D{{Key: "_id", Value: id}}).Decode(&Rewilding)

	if err != nil {
		helpers.ResultEmpty(c, err)
		return
	}

	c.JSON(200, Rewilding)
}

func (r RewildingSearchRepository) Update(c *gin.Context) {
	id, _ := primitive.ObjectIDFromHex(c.Param("id"))
	var payload RewildingSearchRequest
	helpers.Validate(c, &payload)

	lat, _ := primitive.ParseDecimal128(fmt.Sprint(payload.RewildingLat))
	lng, _ := primitive.ParseDecimal128(fmt.Sprint(payload.RewildingLng))

	/*Cache References*/
	var RewildingTypeData models.RefRewildingTypes
	typeId, _ := primitive.ObjectIDFromHex(payload.RewildingType)
	config.DB.Collection("RefRewildingTypes").FindOne(context.TODO(), bson.D{{Key: "_id", Value: typeId}}).Decode(&RewildingTypeData)

	/*Find one*/
	var Rewilding models.Rewilding
	config.DB.Collection("Rewilding").FindOne(context.TODO(), bson.D{{Key: "_id", Value: id}}).Decode(&Rewilding)

	Rewilding.RewildingType = helpers.StringToPrimitiveObjId(payload.RewildingType)
	Rewilding.RewildingTypeData = RewildingTypeData
	Rewilding.RewildingCity = payload.RewildingCity
	Rewilding.RewildingArea = payload.RewildingArea
	Rewilding.RewildingName = payload.RewildingName
	Rewilding.RewildingLat = lat
	Rewilding.RewildingLng = lng

	result, _ := config.DB.Collection("Rewilding").ReplaceOne(context.TODO(), bson.D{{Key: "_id", Value: Rewilding.RewildingID}}, Rewilding)
	fmt.Println(result)
}
