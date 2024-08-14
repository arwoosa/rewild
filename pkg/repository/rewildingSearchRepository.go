package repository

import (
	"context"
	"log"
	"net/http"
	"oosa_rewild/internal/config"
	"oosa_rewild/internal/helpers"
	"oosa_rewild/internal/models"
	"strconv"
	"strings"
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
	reqType := c.Query("type")

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

	var RefRewildingTypes models.RefRewildingTypes
	var includedTypes []string
	// isWaterType := false
	if reqType != "" {
		filter := bson.D{{Key: "ref_rewilding_types_key", Value: reqType}}
		err := config.DB.Collection("RefRewildingTypes").FindOne(context.TODO(), filter).Decode(&RefRewildingTypes)
		if err != nil {
			if helpers.MongoZeroID(RefRewildingTypes.RefRewildingTypesId) {
				helpers.ResponseError(c, "Unsupported type")
				return
			}
		}

		if RefRewildingTypes.RefRewildingTypesKey == "water_related" {
			// isWaterType = true
		} else {
			filteredTypes := strings.Split(RefRewildingTypes.RefRewildingTypesGoogle, ",")
			includedTypes = append(includedTypes, filteredTypes...)
		}
	}

	placesService := helpers.GooglePlacesInitialise(c)
	if placesService == nil {
		log.Fatalf("error %s", err)
	}

	if len(includedTypes) == 0 {
		includedTypes = []string{"national_park", "hiking_area", "campground", "camping_cabin", "park", "playground"}
	}

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

	lat := payload.RewildingLat
	lng := payload.RewildingLng

	geocode := helpers.GoogleMapsGeocode(c, payload.RewildingLat, payload.RewildingLng)
	elevation := helpers.GoogleMapsElevation(c, payload.RewildingLat, payload.RewildingLng)

	location := helpers.GooglePlacesToLocationArray(geocode.AddressComponents)
	area, _ := helpers.GooglePlacesGetArea(geocode.AddressComponents, "administrative_area_level_1")
	_, countryCode := helpers.GooglePlacesGetArea(geocode.AddressComponents, "country")
	rewildingApplyOfficial := false

	insert := models.Rewilding{
		RewildingApplyOfficial: &rewildingApplyOfficial,
		RewildingCity:          payload.RewildingCity,
		RewildingArea:          area,
		RewildingLocation:      location,
		RewildingCountryCode:   countryCode,
		RewildingName:          payload.RewildingName,
		RewildingLat:           lat,
		RewildingLng:           lng,
		RewildingElevation:     elevation.Elevation,
		RewildingCreatedBy:     userDetail.UsersId,
		RewildingCreatedAt:     primitive.NewDateTimeFromTime(time.Now()),
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

	lat := payload.RewildingLat
	lng := payload.RewildingLng

	/*Find one*/
	var Rewilding models.Rewilding
	config.DB.Collection("Rewilding").FindOne(context.TODO(), bson.D{{Key: "_id", Value: id}}).Decode(&Rewilding)

	Rewilding.RewildingCity = payload.RewildingCity
	Rewilding.RewildingArea = payload.RewildingArea
	// Rewilding.RewildingLocation = []string{}
	Rewilding.RewildingName = payload.RewildingName
	Rewilding.RewildingLat = lat
	Rewilding.RewildingLng = lng

	config.DB.Collection("Rewilding").ReplaceOne(context.TODO(), bson.D{{Key: "_id", Value: Rewilding.RewildingID}}, Rewilding)
}
