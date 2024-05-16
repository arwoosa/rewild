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
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PocketListRepository struct{}
type PocketListRequest struct {
	PocketListsName string `json:"pocket_lists_name" validate:"required"`
}

func (r PocketListRepository) Retrieve(c *gin.Context) {
	userDetail := helpers.GetAuthUser(c)
	var results []models.PocketLists

	filters := bson.D{{Key: "pocket_lists_user", Value: userDetail.UsersId}}
	cursor, err := config.DB.Collection("PocketLists").Find(context.TODO(), filters)
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

func (r PocketListRepository) Create(c *gin.Context) {
	userDetail := helpers.GetAuthUser(c)
	var payload PocketListRequest
	validateError := helpers.Validate(c, &payload)
	if validateError != nil {
		return
	}

	insert := models.PocketLists{
		PocketListsUser:      userDetail.UsersId,
		PocketListsName:      payload.PocketListsName,
		PocketListsCount:     0,
		PocketListsCreatedAt: primitive.NewDateTimeFromTime(time.Now()),
	}

	result, err := config.DB.Collection("PocketLists").InsertOne(context.TODO(), insert)
	if err != nil {
		fmt.Println("ERROR", err.Error())
		return
	}

	var PocketLists models.PocketLists
	config.DB.Collection("PocketLists").FindOne(context.TODO(), bson.D{{Key: "_id", Value: result.InsertedID}}).Decode(&PocketLists)
	c.JSON(http.StatusOK, PocketLists)
}

func (r PocketListRepository) Read(c *gin.Context) {
	var PocketLists models.PocketLists
	err := r.ReadOne(c, &PocketLists, "")
	if err == nil {
		c.JSON(http.StatusOK, PocketLists)
	}
}

func (r PocketListRepository) Update(c *gin.Context) {
	// userDetail := user.(*models.Users)
	var payload PocketListRequest
	validateError := helpers.Validate(c, &payload)
	if validateError != nil {
		return
	}
	var PocketLists models.PocketLists
	errRead := r.ReadOne(c, &PocketLists, "")

	if errRead != nil {
		return
	}

	PocketLists.PocketListsName = payload.PocketListsName
	_, errUpd := config.DB.Collection("PocketLists").ReplaceOne(context.TODO(), bson.D{{Key: "_id", Value: PocketLists.PocketListsId}}, PocketLists)
	if errUpd == nil {
		c.JSON(http.StatusOK, PocketLists)
	}
}

func (r PocketListRepository) ReadOne(c *gin.Context, PocketList *models.PocketLists, pocketListId string) error {
	idVal := c.Param("id")
	if pocketListId != "" {
		idVal = pocketListId
	}

	id, _ := primitive.ObjectIDFromHex(idVal)
	/*filter := bson.D{{Key: "_id", Value: id}}
	err := config.DB.Collection("PocketLists").FindOne(context.TODO(), filter).Decode(&PocketList)
	if err != nil {
		helpers.ResultEmpty(c, err)
	}
	return err*/
	return r.ReadById(c, PocketList, id)
}

func (r PocketListRepository) ReadById(c *gin.Context, PocketList *models.PocketLists, id primitive.ObjectID) error {
	fmt.Println("ReadById", id)
	filter := bson.D{{Key: "_id", Value: id}}
	err := config.DB.Collection("PocketLists").FindOne(context.TODO(), filter).Decode(&PocketList)
	if err != nil {
		helpers.ResultEmpty(c, err)
	}
	return err
}

func (r PocketListRepository) Delete(c *gin.Context) {
	var PocketLists models.PocketLists
	err := r.ReadOne(c, &PocketLists, "")
	if err == nil {
		filtersItem := bson.D{{Key: "pocket_list_items_mst", Value: PocketLists.PocketListsId}}
		config.DB.Collection("PocketListItems").DeleteMany(context.TODO(), filtersItem)
		filters := bson.D{{Key: "_id", Value: PocketLists.PocketListsId}}
		config.DB.Collection("PocketLists").DeleteOne(context.TODO(), filters)
		helpers.ResultMessageSuccess(c, "Pocket list deleted")
	}
}

func (r PocketListRepository) UpdateCount(c *gin.Context, pocketListId string) {
	var PocketLists models.PocketLists
	r.ReadOne(c, &PocketLists, pocketListId)

	fmt.Println("UpdateCount:INIT")
	opts := options.Count().SetHint("_id_")
	filter := bson.D{{Key: "pocket_list_items_mst", Value: helpers.StringToPrimitiveObjId(pocketListId)}}
	count, err := config.DB.Collection("PocketListItems").CountDocuments(context.TODO(), filter, opts)
	if err != nil {
		fmt.Println("UpdateCount", err.Error())
	}

	PocketLists.PocketListsCount = int(count)
	config.DB.Collection("PocketLists").ReplaceOne(context.TODO(), bson.D{{Key: "_id", Value: PocketLists.PocketListsId}}, PocketLists)
}
