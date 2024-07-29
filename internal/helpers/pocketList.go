package helpers

import (
	"context"
	"oosa_rewild/internal/config"
	"oosa_rewild/internal/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetPocketList(c *gin.Context, id primitive.ObjectID) models.PocketLists {
	var PocketLists models.PocketLists
	err := config.DB.Collection("PocketLists").FindOne(context.TODO(), bson.D{{Key: "_id", Value: id}}).Decode(&PocketLists)
	ResultEmpty(c, err)
	return PocketLists
}

func GetPocketListItem(c *gin.Context, id primitive.ObjectID) (models.PocketListItems, error) {
	var PocketListItems models.PocketListItems
	err := config.DB.Collection("PocketListItems").FindOne(context.TODO(), bson.D{{Key: "_id", Value: id}}).Decode(&PocketListItems)
	ResultNotFound(c, err, "Pocket list item not found")
	return PocketListItems, err
}
