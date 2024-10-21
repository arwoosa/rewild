package service

import (
	"context"
	"oosa_rewild/internal/config"
	"oosa_rewild/internal/helpers"
	"oosa_rewild/internal/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetRewildingById(rewildingId primitive.ObjectID) (models.Rewilding, error) {
	var Rewilding models.Rewilding
	err := config.DB.Collection("Rewilding").FindOne(context.TODO(), bson.D{{Key: "_id", Value: rewildingId}}).Decode(&Rewilding)
	return Rewilding, err
}

func GoogleToRewilding(c *gin.Context, googlePlaceId string) primitive.ObjectID {
	places := helpers.GooglePlaceById(c, googlePlaceId)
	if places == nil {
		return primitive.ObjectID{}
	}

	/* Create Rewilding */
	location := helpers.GooglePlacesV1ToLocationArray(places.AddressComponents)
	area, _ := helpers.GooglePlacesV1GetArea(places.AddressComponents, "administrative_area_level_1")
	_, countryCode := helpers.GooglePlacesV1GetArea(places.AddressComponents, "country")

	elevation := helpers.GoogleMapsElevation(c, places.Location.Latitude, places.Location.Longitude)

	rewilding, inDB := helpers.GetRewildByPlaceId(places.Id)
	var rewildingId primitive.ObjectID
	if inDB {
		rewildingId = rewilding.RewildingID
	}

	RewildingPhotos, _ := helpers.RewildSaveGooglePhotos(c, places.Photos)

	newRewilding := bson.D{
		{Key: "rewilding_area", Value: area},
		{Key: "rewilding_location", Value: location},
		{Key: "rewilding_country_code", Value: countryCode},
		{Key: "rewilding_name", Value: places.DisplayName.Text},
		{Key: "rewilding_lat", Value: places.Location.Latitude},
		{Key: "rewilding_lng", Value: places.Location.Longitude},
		{Key: "rewilding_place_id", Value: places.Id},
		{Key: "rewilding_elevation", Value: elevation.Elevation},
		{Key: "rewilding_photos", Value: RewildingPhotos},
		{Key: "rewilding_apply_official", Value: true},
	}
	RefAchievementPlaces, RefAchievementPlacesErr := helpers.RewildingAchievementByLatLng(c, places.Location.Latitude, places.Location.Longitude)
	if RefAchievementPlacesErr == nil {
		newRewilding = append(newRewilding,
			bson.E{Key: "rewilding_achievement_type", Value: RefAchievementPlaces.RefAchievementPlacesType},
			bson.E{Key: "rewilding_achievement_type_id", Value: RefAchievementPlaces.RefAchievementPlacesID},
		)
	}
	upsert := bson.D{{Key: "$set", Value: newRewilding}}

	filters := bson.D{
		{Key: "rewilding_place_id", Value: places.Id},
	}
	opts := options.Update().SetUpsert(true)
	result, _ := config.DB.Collection("Rewilding").UpdateOne(context.TODO(), filters, upsert, opts)

	if !inDB {
		rewildingId = result.UpsertedID.(primitive.ObjectID)
	}

	return rewildingId
}
