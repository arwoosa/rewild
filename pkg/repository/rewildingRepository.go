package repository

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"oosa_rewild/internal/config"
	"oosa_rewild/internal/helpers"
	"oosa_rewild/internal/middleware"
	"oosa_rewild/internal/models"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/api/places/v1"
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

type RewildingAutocomplete struct {
	PlaceId  string   `json:"place_id"`
	Name     string   `json:"name"`
	Location string   `json:"location"`
	Type     []string `json:"type"`
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

	if owner == "true" {
		middleware.CheckIfAuth(c)
		userDetail := helpers.GetAuthUser(c)
		agg = append(agg, bson.D{{
			Key: "$match", Value: bson.D{
				{Key: "rewilding_created_by", Value: userDetail.UsersId},
				{Key: "rewilding_deleted_at", Value: bson.M{"$exists": false}},
			},
		}})
	} else {
		agg = append(agg, bson.D{{
			Key: "$match", Value: bson.D{
				{Key: "rewilding_deleted_at", Value: bson.M{"$exists": false}},
			},
		}})
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

func (r RewildingRepository) GetOneRewilding(stringId string, Rewilding *models.Rewilding) error {
	id, _ := primitive.ObjectIDFromHex(stringId)
	err := config.DB.Collection("Rewilding").FindOne(context.TODO(), bson.D{{Key: "_id", Value: id}}).Decode(&Rewilding)
	return err
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

	if len(payload.RewildingPocketList) > 0 {
		var PocketLists []models.PocketLists
		PocketListErr := []string{}
		orFilter := []bson.M{}
		for _, v := range payload.RewildingPocketList {
			orFilter = append(orFilter, bson.M{"_id": helpers.StringToPrimitiveObjId(v)})
		}
		filter := bson.D{
			{Key: "$or", Value: orFilter},
		}

		cursor, _ := config.DB.Collection("PocketLists").Find(context.TODO(), filter)
		cursor.All(context.TODO(), &PocketLists)

		if len(PocketLists) != len(payload.RewildingPocketList) {
			helpers.ResponseError(c, "Invalid pocket list")
			return
		}

		for _, v := range PocketLists {
			limit := int(config.APP_LIMIT.PocketListItems)
			if v.PocketListsCount >= limit {
				errMessage := "Cannot add to " + v.PocketListsName + " has reached limit of " + strconv.Itoa(limit)
				PocketListErr = append(PocketListErr, errMessage)
			}
		}

		if len(PocketListErr) > 0 {
			helpers.ResponseError(c, strings.Join(PocketListErr, ", "))
			return
		}
	}

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
			meta, _ := LinkRepository{}.GetMeta(c, referenceInformation)
			insRefLink := models.RewildingReferenceLinks{
				RewildingReferenceLinksLink:          meta.Url,
				RewildingReferenceLinksTitle:         meta.Title,
				RewildingReferenceLinksDescription:   meta.Description,
				RewildingReferenceLinksOGTitle:       meta.OGTitle,
				RewildingReferenceLinksOGDescription: meta.OGDescription,
				RewildingReferenceLinksOGImage:       meta.OGImage,
				RewildingReferenceLinksOGSiteName:    meta.OGTitle,
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

func (r RewildingRepository) Delete(c *gin.Context) {
	userDetail := helpers.GetAuthUser(c)
	id, _ := primitive.ObjectIDFromHex(c.Param("id"))
	var Rewilding models.Rewilding
	filter := bson.D{
		{Key: "_id", Value: id},
		{Key: "rewilding_created_by", Value: userDetail.UsersId},
		{Key: "rewilding_deleted_at", Value: bson.M{"$exists": false}},
	}
	err := config.DB.Collection("Rewilding").FindOne(context.TODO(), filter).Decode(&Rewilding)
	fmt.Println(userDetail.UsersId)
	if err != nil {
		helpers.ResultEmpty(c, err)
		return
	}

	updFilter := bson.D{{Key: "_id", Value: id}}
	currentTime := primitive.NewDateTimeFromTime(time.Now())
	Rewilding.RewildingDeletedAt = &currentTime
	Rewilding.RewildingDeletedBy = &userDetail.UsersId
	upd := bson.D{{Key: "$set", Value: Rewilding}}
	config.DB.Collection("Rewilding").UpdateOne(context.TODO(), updFilter, upd)

	helpers.ResponseSuccessMessage(c, "Rewilding deleted")
}

func (r RewildingRepository) Options(c *gin.Context) {
	RefRewildingTypes := helpers.RefRewildingTypes()
	c.JSON(http.StatusOK, gin.H{"rewilding_types": RefRewildingTypes})
}

func (r RewildingRepository) Action(c *gin.Context) {
	action := c.Param("action")

	switch action {
	case "rewilding:searchText":
		r.SearchText(c)
		return
	case "rewilding:searchNearby":
		r.SearchNearby(c)
		return
	case "rewilding:autocomplete":
		r.Autocomplete(c)
		return
	}
}

func (r RewildingRepository) SearchText(c *gin.Context) {
	reqSearch := c.Query("keyword")
	reqLanguage := c.Query("language")
	reqRectLowLat := c.Query("rectangle_low_lat")
	reqRectLowLng := c.Query("rectangle_low_lng")
	reqRectHightLat := c.Query("rectangle_hight_lat")
	reqRectHightLng := c.Query("rectangle_hight_lng")

	placesService := helpers.GooglePlacesInitialise(c)
	if placesService == nil {
		return
	}

	languageCode := "zh-TW"

	if reqLanguage != "" {
		languageCode = reqLanguage
	}

	req := places.GoogleMapsPlacesV1SearchTextRequest{
		TextQuery:      reqSearch,
		MaxResultCount: 10,
		LanguageCode:   languageCode,
	}

	if reqRectLowLat != "" && reqRectLowLng != "" && reqRectHightLat != "" && reqRectHightLng != "" {
		req.LocationBias = &places.GoogleMapsPlacesV1SearchTextRequestLocationBias{
			Rectangle: &places.GoogleGeoTypeViewport{
				Low: &places.GoogleTypeLatLng{
					Latitude:  helpers.StringToFloat(reqRectLowLat),
					Longitude: helpers.StringToFloat(reqRectLowLng),
				},
				High: &places.GoogleTypeLatLng{
					Latitude:  helpers.StringToFloat(reqRectHightLat),
					Longitude: helpers.StringToFloat(reqRectHightLng),
				},
			},
		}
	}

	placeReq := placesService.Places.SearchText(&req)
	placeReq.Header().Add("X-Goog-FieldMask", "places.id,places.types,places.displayName,places.formattedAddress,places.location,places.rating,places.userRatingCount")
	places, errPlace := placeReq.Do()
	if errPlace != nil {
		log.Fatalf("error %s", errPlace)
	}

	Rewilding := r.GooglePlaceToRewildingList(c, places.Places)
	c.JSON(http.StatusOK, Rewilding)
}

func (r RewildingRepository) SearchNearby(c *gin.Context) {
	reqType := c.Query("type")
	reqLat := c.Query("lat")
	reqLng := c.Query("lng")
	reqRadius := c.Query("radius")
	reqLanguage := c.Query("language")

	placesService := helpers.GooglePlacesInitialise(c)
	if placesService == nil {
		return
	}

	radius := 5000.00
	languageCode := "zh-TW"
	includedTypes := []string{"national_park", "hiking_area", "campground", "camping_cabin", "park", "playground"}

	if reqRadius != "" {
		radius = helpers.StringToFloat(reqRadius)
	}
	if reqLanguage != "" {
		languageCode = reqLanguage
	}
	if reqType != "" {
		includedTypes = strings.Split(reqType, ",")
	}

	req := places.GoogleMapsPlacesV1SearchNearbyRequest{
		IncludedTypes:  includedTypes,
		MaxResultCount: 10,
		LocationRestriction: &places.GoogleMapsPlacesV1SearchNearbyRequestLocationRestriction{
			Circle: &places.GoogleMapsPlacesV1Circle{
				Center: &places.GoogleTypeLatLng{
					Latitude:  helpers.StringToFloat(reqLat),
					Longitude: helpers.StringToFloat(reqLng),
				},
				Radius: radius,
			},
		},
		LanguageCode: languageCode,
	}

	placeReq := placesService.Places.SearchNearby(&req)
	placeReq.Header().Add("X-Goog-FieldMask", "places.id,places.types,places.displayName,places.formattedAddress,places.location,places.rating,places.userRatingCount")
	places, errPlace := placeReq.Do()
	if errPlace != nil {
		log.Fatalf("error %s", errPlace)
	}

	Rewilding := r.GooglePlaceToRewildingList(c, places.Places)
	c.JSON(http.StatusOK, Rewilding)
}

func (r RewildingRepository) Autocomplete(c *gin.Context) {
	var Autocomplete []RewildingAutocomplete
	reqInput := c.Query("input")
	reqLanguage := c.Query("language")

	placesService := helpers.GooglePlacesInitialise(c)
	if placesService == nil {
		return
	}

	languageCode := "zh-TW"
	if reqLanguage != "" {
		languageCode = reqLanguage
	}

	req := places.GoogleMapsPlacesV1AutocompletePlacesRequest{
		Input:        reqInput,
		LanguageCode: languageCode,
	}

	placeReq := placesService.Places.Autocomplete(&req)
	places, errPlace := placeReq.Do()
	if errPlace != nil {
		log.Fatalf("error %s", errPlace)
	}

	for _, v := range places.Suggestions {
		entry := RewildingAutocomplete{
			PlaceId:  v.PlacePrediction.PlaceId,
			Name:     v.PlacePrediction.StructuredFormat.MainText.Text,
			Location: v.PlacePrediction.StructuredFormat.SecondaryText.Text,
			Type:     v.PlacePrediction.Types,
		}
		Autocomplete = append(Autocomplete, entry)
	}

	c.JSON(http.StatusOK, Autocomplete)
}

func (r RewildingRepository) Places(c *gin.Context) {
	placesId := c.Param("placesId")
	placeRewilding, gplaces := r.GooglePlaceToRewilding(c, placesId)

	if gplaces != nil {
		c.JSON(200, placeRewilding)
	}
}

func (r RewildingRepository) GooglePlaceToRewildingList(c *gin.Context, placeObject []*places.GoogleMapsPlacesV1Place) []models.Rewilding {
	var Rewilding []models.Rewilding
	for _, v := range placeObject {
		placeRewilding, gplaces := r.GooglePlaceToRewilding(c, v.Id)
		if gplaces != nil {
			Rewilding = append(Rewilding, placeRewilding)
		}
	}
	return Rewilding
}

func (r RewildingRepository) GooglePlaceToRewilding(c *gin.Context, placeId string) (models.Rewilding, *places.GoogleMapsPlacesV1Place) {
	gplaces := helpers.GooglePlaceById(c, placeId)

	if gplaces != nil {
		location := helpers.GooglePlacesV1ToLocationArray(gplaces.AddressComponents)
		area, _ := helpers.GooglePlacesV1GetArea(gplaces.AddressComponents, "administrative_area_level_1")
		_, countryCode := helpers.GooglePlacesV1GetArea(gplaces.AddressComponents, "country")

		elevation := helpers.GoogleMapsElevation(c, gplaces.Location.Latitude, gplaces.Location.Longitude)
		RewildingPhotos := helpers.RewildGooglePhotos(c, gplaces.Photos)
		applyOfficial := false

		return models.Rewilding{
			RewildingArea:          area,
			RewildingLocation:      location,
			RewildingCountryCode:   countryCode,
			RewildingName:          gplaces.DisplayName.Text,
			RewildingLat:           gplaces.Location.Latitude,
			RewildingLng:           gplaces.Location.Longitude,
			RewildingPlaceId:       gplaces.Id,
			RewildingElevation:     elevation.Elevation,
			RewildingPhotos:        RewildingPhotos,
			RewildingApplyOfficial: &applyOfficial,
		}, gplaces
	}
	return models.Rewilding{}, gplaces
}
