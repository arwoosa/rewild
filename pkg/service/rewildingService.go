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
	area := helpers.GooglePlacesV1GetArea(places.AddressComponents, "administrative_area_level_1")
	elevation := helpers.GoogleMapsElevation(c, places.Location.Latitude, places.Location.Longitude)

	rewilding, inDB := helpers.GetRewildByPlaceId(places.Id)
	rewildingId := rewilding.RewildingID

	if !inDB {
		upsert := bson.D{{Key: "$set", Value: bson.D{
			{Key: "rewilding_area", Value: area},
			{Key: "rewilding_name", Value: places.DisplayName.Text},
			{Key: "rewilding_lat", Value: helpers.FloatToDecimal128(places.Location.Latitude)},
			{Key: "rewilding_lng", Value: helpers.FloatToDecimal128(places.Location.Longitude)},
			{Key: "rewilding_place_id", Value: places.Id},
			{Key: "rewilding_elevation", Value: elevation.Elevation},
		}}}

		filters := bson.D{
			{Key: "rewilding_place_id", Value: places.Id},
		}
		opts := options.Update().SetUpsert(true)
		result, _ := config.DB.Collection("Rewilding").UpdateOne(context.TODO(), filters, upsert, opts)
		rewildingId = result.UpsertedID.(primitive.ObjectID)
	}

	return rewildingId
}
