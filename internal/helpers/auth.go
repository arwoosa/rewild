package helpers

import (
	"context"
	"oosa_rewild/internal/config"
	"oosa_rewild/internal/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetAuthUserByCheckHeaders(c *gin.Context) models.Users {
	user := GetAuthUser(c)
	if !user.UsersId.IsZero() {
		return user
	}
	headerUserId := c.GetHeader("X-User-Id")
	if headerUserId == "" || headerUserId == "guest" {
		return user
	}
	err := config.DB.Collection("Users").FindOne(context.TODO(), bson.D{{Key: "users_source_id", Value: headerUserId}}).Decode(&user)
	if err != nil {
		return models.Users{}
	}
	return user
}

func GetAuthUser(c *gin.Context) models.Users {
	user, exists := c.Get("user")
	if exists {
		userDetail := user.(*models.Users)
		return *userDetail
	} else {
		return models.Users{}
	}
}

func FindUserSourceId(userIds []primitive.ObjectID) (map[primitive.ObjectID]string, error) {
	collection := config.DB.Collection("Users")

	var usersDoc []models.Users
	cursor, err := collection.Find(context.TODO(), bson.M{"_id": bson.M{"$in": userIds}})
	if err != nil {
		return nil, err
	}
	err = cursor.All(context.TODO(), &usersDoc)
	if err != nil {
		return nil, err
	}
	result := map[primitive.ObjectID]string{}
	for _, u := range usersDoc {
		result[u.UsersId] = u.UsersSourceId
	}
	return result, nil
}
