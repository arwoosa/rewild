package repository

import (
	"context"
	"fmt"
	"net/http"
	"oosa_rewild/internal/config"
	"oosa_rewild/internal/helpers"
	"oosa_rewild/internal/models"
	"oosa_rewild/pkg/service"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type PocketListItemsRepository struct{}
type PocketListItemsRequest struct {
	PocketListItemsPlaceId   string `json:"pocket_list_items_place_id" validate:"required_without=PocketListItemsRewilding"`
	PocketListItemsPlaceName string `json:"pocket_list_items_place_name" validate:"required_without=PocketListItemsRewilding"`
	PocketListItemsRewilding string `json:"pocket_list_items_rewilding_id" validate:"required_without=PocketListItemsPlaceId"`
}
type PocketListItemsUpdateRequest struct {
	PocketListItemsMst string `json:"pocket_list_items_mst"`
}
type PocketListItemsUpdateBulkRequest struct {
	PocketListItemsId  []string `json:"pocket_list_items_id" validate:"required"`
	PocketListItemsMst string   `json:"pocket_list_items_mst"`
}
type PocketListItemsDeleteBulkRequest struct {
	PocketListItemsId []string `json:"pocket_list_items_id" validate:"required"`
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

	r.RetrievePhoto(&results)
	c.JSON(http.StatusOK, results)
}

func (r PocketListItemsRepository) RetrievePhoto(results *[]models.PocketListItems) {
	for idx, val := range *results {
		if len(val.PocketListItemsRewildingDetail.RewildingPhotos) > 0 {
			for photoIdx, photo := range val.PocketListItemsRewildingDetail.RewildingPhotos {
				if photo.RewildingPhotosPath == "" {
					rewildingId := val.PocketListItemsRewildingDetail.RewildingID
					photoId := photo.RewildingPhotosID
					(*results)[idx].PocketListItemsRewildingDetail.RewildingPhotos[photoIdx].RewildingPhotosPath = config.APP.BaseUrl + "rewilding/" + rewildingId.Hex() + "/photos/" + photoId.Hex()
				}
			}
		}
	}
}

func (r PocketListItemsRepository) Create(c *gin.Context) {
	pocketListId := c.Param("id")
	var payload PocketListItemsRequest
	var PocketLists models.PocketLists
	validateError := helpers.Validate(c, &payload)
	if validateError != nil {
		return
	}

	err := PocketListRepository{}.ReadOne(c, &PocketLists, "")
	if err != nil {
		return
	}

	maxPocketListItems := int(config.APP_LIMIT.PocketListItems)
	if PocketLists.PocketListsCount >= maxPocketListItems {
		helpers.ResponseError(c, "Cannot add to pocket list. Max allowed "+strconv.Itoa(int(maxPocketListItems)))
		return
	}

	var insert models.PocketListItems
	PocketListID := helpers.StringToPrimitiveObjId(pocketListId)
	var PocketListRewildingId primitive.ObjectID
	var PocketListItemsPlaceName string

	if payload.PocketListItemsPlaceId != "" {
		places := helpers.GooglePlaceById(c, payload.PocketListItemsPlaceId)
		if places == nil {
			return
		}

		PocketListRewildingId = service.GoogleToRewilding(c, places.Id)
		PocketListItemsPlaceName = payload.PocketListItemsPlaceName
	} else if payload.PocketListItemsRewilding != "" {
		var Rewilding models.Rewilding
		err := RewildingRepository{}.GetOneRewilding(payload.PocketListItemsRewilding, &Rewilding)
		if err == mongo.ErrNoDocuments {
			helpers.ResponseError(c, "Invalid rewilding ID")
			return
		}

		PocketListRewildingId = Rewilding.RewildingID
		PocketListItemsPlaceName = Rewilding.RewildingName
	}

	checkIfAvailable := r.GetPocketListItemByMstRewildingId(PocketListID, PocketListRewildingId)

	if checkIfAvailable != mongo.ErrNoDocuments {
		helpers.ResponseError(c, PocketListItemsPlaceName+" ID available in pocket list")
		return
	}

	insert = models.PocketListItems{
		PocketListItemsMst:       PocketListID,
		PocketListItemsRewilding: PocketListRewildingId,
		PocketListItemsName:      PocketListItemsPlaceName,
		PocketListItemsCreatedAt: primitive.NewDateTimeFromTime(time.Now()),
	}

	result, err := config.DB.Collection("PocketListItems").InsertOne(context.TODO(), insert)
	if err != nil {
		fmt.Println("ERROR", err.Error())
		return
	}

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

func (r PocketListItemsRepository) UpdateBulk(c *gin.Context) {
	pocketListId := c.Param("id")
	var payload PocketListItemsUpdateBulkRequest
	var PocketLists models.PocketLists
	var PocketListItemsId []primitive.ObjectID
	validateError := helpers.Validate(c, &payload)
	if validateError != nil {
		return
	}

	err := PocketListRepository{}.ReadOne(c, &PocketLists, "")
	if err != nil {
		return
	}

	var newPocketList models.PocketLists
	errNewPocketList := PocketListRepository{}.ReadById(c, &newPocketList, helpers.StringToPrimitiveObjId(payload.PocketListItemsMst))
	if errNewPocketList != nil {
		return
	}

	PocketListId := helpers.StringToPrimitiveObjId(pocketListId)
	NewPocketListId := helpers.StringToPrimitiveObjId(payload.PocketListItemsMst)

	for _, v := range payload.PocketListItemsId {
		PocketListItemsId = append(PocketListItemsId, helpers.StringToPrimitiveObjId(v))
	}

	PocketListItemsOld, filterItemsOld := r.GetPocketListItems(PocketListId, PocketListItemsId)
	PocketListItemsNew, _ := r.GetPocketListItems(NewPocketListId, nil)

	// Make sure not in new pocket list ID
	var PocketListNewRewilding []string
	var PocketListNewError []string
	var PocketListMoved []string
	for _, v := range PocketListItemsNew {
		PocketListNewRewilding = append(PocketListNewRewilding, v.PocketListItemsRewilding.Hex())
	}

	for _, v := range PocketListItemsOld {
		if helpers.StringInSlice(v.PocketListItemsRewilding.Hex(), PocketListNewRewilding) {
			PocketListNewError = append(PocketListNewError, v.PocketListItemsName)
		} else {
			PocketListMoved = append(PocketListMoved, v.PocketListItemsName)
		}
	}

	maxPocketListItems := int(config.APP_LIMIT.PocketListItems)
	if newPocketList.PocketListsCount >= maxPocketListItems {
		helpers.ResponseError(c, "Cannot add to pocket list. Max allowed "+strconv.Itoa(int(maxPocketListItems)))
		return
	}

	if len(PocketListNewError) > 0 {
		grammar := "is"
		if len(PocketListNewError) > 1 {
			grammar = "are"
		}
		helpers.ResponseBadRequestError(c, strings.Join(PocketListNewError, ", ")+" "+grammar+" already in '"+newPocketList.PocketListsName+"' and cannot be moved")
		return
	}

	if len(payload.PocketListItemsId) != len(PocketListItemsOld) {
		helpers.ResponseBadRequestError(c, "Some items staged for moving does not belong to this pocket list")
		return
	}

	upd := bson.D{{Key: "$set", Value: bson.D{{Key: "pocket_list_items_mst", Value: NewPocketListId}}}}
	_, updErr := config.DB.Collection("PocketListItems").UpdateMany(context.TODO(), filterItemsOld, upd)

	if updErr != nil {
		helpers.ResponseBadRequestError(c, updErr.Error())
		return
	}

	PocketListRepository{}.UpdateCount(c, pocketListId)
	PocketListRepository{}.UpdateCount(c, payload.PocketListItemsMst)
	helpers.ResponseSuccessMessage(c, strings.Join(PocketListMoved, ", ")+" moved to '"+newPocketList.PocketListsName+"'")
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

func (r PocketListItemsRepository) DeleteMultiple(c *gin.Context) {
	pocketListId := c.Param("id")
	PocketListId := helpers.StringToPrimitiveObjId(pocketListId)
	var payload PocketListItemsDeleteBulkRequest
	var PocketLists models.PocketLists
	var PocketListItemsId []primitive.ObjectID
	validateError := helpers.Validate(c, &payload)
	if validateError != nil {
		return
	}

	err := PocketListRepository{}.ReadOne(c, &PocketLists, "")
	if err != nil {
		return
	}

	for _, v := range payload.PocketListItemsId {
		PocketListItemsId = append(PocketListItemsId, helpers.StringToPrimitiveObjId(v))
	}

	PocketListItems, filterDelete := r.GetPocketListItems(PocketListId, PocketListItemsId)

	if len(payload.PocketListItemsId) != len(PocketListItems) {
		helpers.ResponseBadRequestError(c, "Some items staged for deletion does not belong to this pocket list")
		return
	}

	config.DB.Collection("PocketListItems").DeleteMany(context.TODO(), filterDelete)
	PocketListRepository{}.UpdateCount(c, pocketListId)
	helpers.ResponseSuccessMessage(c, "Pocket list items removed from "+PocketLists.PocketListsName)
}

func (r PocketListItemsRepository) GetPocketListItems(PocketListId primitive.ObjectID, PocketListItemsId []primitive.ObjectID) ([]models.PocketListItems, primitive.D) {
	var PocketListItems []models.PocketListItems
	filter := bson.D{
		{Key: "pocket_list_items_mst", Value: PocketListId},
	}

	if PocketListItemsId != nil {
		filter = append(filter, primitive.E{Key: "_id", Value: bson.M{"$in": PocketListItemsId}})
	}

	cursor, _ := config.DB.Collection("PocketListItems").Find(context.TODO(), filter)
	cursor.All(context.TODO(), &PocketListItems)
	return PocketListItems, filter
}

func (r PocketListItemsRepository) GetPocketListItemByMstRewildingId(mstId primitive.ObjectID, rewildingId primitive.ObjectID) error {
	var PocketListItems models.PocketListItems
	filter := bson.D{
		{Key: "pocket_list_items_mst", Value: mstId},
		{Key: "pocket_list_items_rewilding", Value: rewildingId},
	}
	err := config.DB.Collection("PocketListItems").FindOne(context.TODO(), filter).Decode(&PocketListItems)
	return err
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
