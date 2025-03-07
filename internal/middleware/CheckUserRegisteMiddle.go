package middleware

import (
	"context"
	"net/http"
	"oosa_rewild/internal/config"
	"oosa_rewild/internal/models"
	"oosa_rewild/internal/structs"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func CheckRegisterMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		headerUserId := c.GetHeader("X-User-Id")
		if headerUserId == "" {
			return
		}
		if headerUserId != "" && headerUserId == "guest" {
			return
		}

		var headerUser structs.UserBindByHeader
		err := c.BindHeader(&headerUser)
		if err != nil || headerUser.Id == "" {
			return
		}
		var user *models.Users
		err = config.DB.Collection("Users").FindOne(context.TODO(), bson.D{{Key: "users_source_id", Value: headerUser.Id}}).Decode(user)

		if err != nil && err == mongo.ErrNoDocuments {
			user, err = saveUserInfo(c, &headerUser)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": "AUTH01-USER: " + err.Error()})
				c.Abort()
				return
			}
		}

		needUpdate := false
		if user.UsersAvatar == "" {
			user.UsersAvatar = headerUser.Avatar
			needUpdate = true
		}
		if user.UsersUsername == "" {
			user.UsersUsername = headerUser.User
			needUpdate = true
		}
		if needUpdate {
			config.DB.Collection("Users").UpdateByID(c, user.UsersId, bson.D{{Key: "$set", Value: bson.D{{Key: "users_avatar", Value: user.UsersAvatar}, {Key: "users_username", Value: user.UsersUsername}}}})
		}

		c.Set("user", user)

	}
}

var defaultFriendAutoAdd = 0
var defaultFriendTakeMestatus = false

func saveUserInfo(c context.Context, user *structs.UserBindByHeader) (*models.Users, error) {
	var User models.Users

	insert := models.Users{
		UsersSource:                           3,
		UsersSourceId:                         user.Id,
		UsersName:                             user.Name,
		UsersEmail:                            user.Email,
		UsersUsername:                         user.User,
		UsersObject:                           user.User,
		UsersAvatar:                           user.Avatar,
		UsersSettingLanguage:                  user.Language,
		UsersSettingIsVisibleFriends:          1,
		UsersSettingIsVisibleStatistics:       1,
		UsersSettingVisibilityActivitySummary: 1,
		UsersSettingFriendAutoAdd:             &defaultFriendAutoAdd,
		UsersIsSubscribed:                     false,
		UsersIsBusiness:                       false,
		UsersTakeMeStatus:                     &defaultFriendTakeMestatus,
		UsersCreatedAt:                        primitive.NewDateTimeFromTime(time.Now()),
	}
	result, _ := config.DB.Collection("Users").InsertOne(c, insert)
	config.DB.Collection("Users").FindOne(c, bson.D{{Key: "_id", Value: result.InsertedID}}).Decode(&User)
	return &User, nil
}
