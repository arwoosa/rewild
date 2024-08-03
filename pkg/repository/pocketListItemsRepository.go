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
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PocketListItemsRepository struct{}
type PocketListItemsRequest struct {
	PocketListItemsPlaceId   string `json:"pocket_list_items_place_id"`
	PocketListItemsPlaceName string `json:"pocket_list_items_place_name" validate:"required"`
}
type PocketListItemsUpdateRequest struct {
	PocketListItemsMst string `json:"pocket_list_items_mst"`
}

func (r PocketListItemsRepository) Retrieve(c *gin.Context) {
	var results []models.PocketListItems
	rewildingId := helpers.StringToPrimitiveObjId(c.Param("id"))
	agg := mongo.Pipeline{
		bson.D{{
			Key: "$match", Value: bson.M{
				"pocket_list_items_mst": rewildingId,
			},
		}},
		bson.D{{
			Key: "$lookup", Value: bson.M{
				"from":         "Rewilding",
				"localField":   "pocket_list_items_rewilding",
				"foreignField": "_id",
				"as":           "pocket_list_items_rewilding_detail",
			},
		}},
		bson.D{{
			Key: "$unwind", Value: "$pocket_list_items_rewilding_detail",
		}},
	}

	err := PocketListRepository{}.ReadOne(c, &models.PocketLists{}, "")
	if err != nil {
		return
	}

	cursor, err := config.DB.Collection("PocketListItems").Aggregate(context.TODO(), agg)
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

func (r PocketListItemsRepository) Create(c *gin.Context) {
	pocketListId := c.Param("id")
	var payload PocketListItemsRequest
	validateError := helpers.Validate(c, &payload)
	if validateError != nil {
		return
	}

	err := PocketListRepository{}.ReadOne(c, &models.PocketLists{}, "")
	if err != nil {
		return
	}

	places := helpers.GooglePlaceById(c, payload.PocketListItemsPlaceId)
	if places == nil {
		return
	}

	/* Create Rewilding */
	location := helpers.GooglePlacesV1ToLocationArray(places.AddressComponents)
	area, _ := helpers.GooglePlacesV1GetArea(places.AddressComponents, "administrative_area_level_1")
	_, countryCode := helpers.GooglePlacesV1GetArea(places.AddressComponents, "country")

	elevation := helpers.GoogleMapsElevation(c, places.Location.Latitude, places.Location.Longitude)

	rewilding, inDB := helpers.GetRewildByPlaceId(places.Id)
	rewildingId := rewilding.RewildingID

	if !inDB {
		RewildingPhotos, _ := helpers.RewildSaveGooglePhotos(c, places.Photos)

		upsert := bson.D{{Key: "$set", Value: bson.D{
			{Key: "rewilding_area", Value: area},
			{Key: "rewilding_location", Value: location},
			{Key: "rewilding_country_code", Value: countryCode},
			{Key: "rewilding_name", Value: places.DisplayName.Text},
			{Key: "rewilding_lat", Value: places.Location.Latitude},
			{Key: "rewilding_lng", Value: places.Location.Longitude},
			{Key: "rewilding_place_id", Value: places.Id},
			{Key: "rewilding_elevation", Value: elevation.Elevation},
			{Key: "rewilding_photos", Value: RewildingPhotos},
		}}}

		filters := bson.D{
			{Key: "rewilding_place_id", Value: places.Id},
		}
		opts := options.Update().SetUpsert(true)
		result, _ := config.DB.Collection("Rewilding").UpdateOne(context.TODO(), filters, upsert, opts)
		rewildingId = result.UpsertedID.(primitive.ObjectID)
	} else {
		if len(rewilding.RewildingPhotos) == 0 {
			RewildingPhotos, _ := helpers.RewildSaveGooglePhotos(c, places.Photos)

			filters := bson.D{
				{Key: "rewilding_place_id", Value: places.Id},
			}

			rewilding.RewildingPhotos = RewildingPhotos
			upd := bson.D{{Key: "$set", Value: rewilding}}
			config.DB.Collection("Rewilding").UpdateOne(context.TODO(), filters, upd)
		}
	}

	insert := models.PocketListItems{
		PocketListItemsMst:       helpers.StringToPrimitiveObjId(pocketListId),
		PocketListItemsRewilding: rewildingId,
		PocketListItemsName:      payload.PocketListItemsPlaceName,
		PocketListItemsCreatedAt: primitive.NewDateTimeFromTime(time.Now()),
	}

	result, err := config.DB.Collection("PocketListItems").InsertOne(context.TODO(), insert)
	if err != nil {
		fmt.Println("ERROR", err.Error())
		return
	}

	fmt.Println(result.InsertedID.(primitive.ObjectID))
	data, err := helpers.GetPocketListItem(c, result.InsertedID.(primitive.ObjectID))

	PocketListRepository{}.UpdateCount(c, pocketListId)
	if err == nil {
		c.JSON(http.StatusOK, data)
	}
}

func (r PocketListItemsRepository) Read(c *gin.Context) {
	err := PocketListRepository{}.ReadOne(c, &models.PocketLists{}, "")
	if err != nil {
		return
	}

	id, _ := primitive.ObjectIDFromHex(c.Param("itemsId"))
	data, err := helpers.GetPocketListItem(c, id)
	if err == nil {
		c.JSON(http.StatusOK, data)
	}
}

func (r PocketListItemsRepository) Update(c *gin.Context) {
	pocketListId := c.Param("id")
	pocketListItemId := c.Param("itemsId")
	var payload PocketListItemsUpdateRequest
	validateError := helpers.Validate(c, &payload)
	if validateError != nil {
		return
	}

	err := PocketListRepository{}.ReadOne(c, &models.PocketLists{}, "")
	if err != nil {
		return
	}

	id, _ := primitive.ObjectIDFromHex(pocketListItemId)
	PocketListItems, err := helpers.GetPocketListItem(c, id)
	if err == nil {
		if payload.PocketListItemsMst != "" {
			var newPocketList models.PocketLists
			errNewPocketList := PocketListRepository{}.ReadById(c, &newPocketList, helpers.StringToPrimitiveObjId(payload.PocketListItemsMst))
			if errNewPocketList != nil {
				return
			}

			PocketListItems.PocketListItemsMst = helpers.StringToPrimitiveObjId(payload.PocketListItemsMst)
			filters := bson.D{{Key: "_id", Value: PocketListItems.PocketListItemsId}}
			upd := bson.D{{Key: "$set", Value: PocketListItems}}
			config.DB.Collection("PocketListItems").UpdateOne(context.TODO(), filters, upd)

			PocketListRepository{}.UpdateCount(c, pocketListId)
			PocketListRepository{}.UpdateCount(c, payload.PocketListItemsMst)
		}
		helpers.ResultMessageSuccess(c, "Pocket list item record updated")
	}
}

func (r PocketListItemsRepository) Delete(c *gin.Context) {
	pocketListId := c.Param("id")
	err := PocketListRepository{}.ReadOne(c, &models.PocketLists{}, "")
	if err != nil {
		return
	}

	id, _ := primitive.ObjectIDFromHex(c.Param("itemsId"))
	PocketListItems, err := helpers.GetPocketListItem(c, id)
	if err == nil {
		filters := bson.D{{Key: "_id", Value: PocketListItems.PocketListItemsId}}
		config.DB.Collection("PocketListItems").DeleteOne(context.TODO(), filters)
		PocketListRepository{}.UpdateCount(c, pocketListId)
		helpers.ResultMessageSuccess(c, "Pocket list item record deleted")
	}
}

/*
SNIPPET: Joining tables
db.PocketListItems.aggregate([
	{ "$addFields": {
			"pocket_list_id": { "$toObjectId": "$pocket_list_items_mst" },
			"rewilding_id": { "$toObjectId": "$pocket_list_items_rewilding" }
		}
	},
	{	"$lookup": {
			from: 'PocketLists',
			localField: 'pocket_list_id',
			foreignField: '_id',
			as: 'test'
		}
	},
	{	"$lookup": {
			from: 'Rewilding',
			localField: 'rewilding_id',
			foreignField: '_id',
			as: 'rewilding'
		}
	},
	{ "$unwind": "$test" },
	{ "$unwind": "$rewilding" },
	{
    "$set": {
      "pocket_lists_name": "$test.pocket_lists_name",
			"pocket_lists_count": "$test.pocket_lists_count",
			"pocket_list_items_rewilding_name": "$rewilding.rewilding_name",
    }
  },
	{
		"$unset": ["test"]
	}
])
*/
