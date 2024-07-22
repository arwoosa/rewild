package helpers

import (
	"context"
	"oosa_rewild/internal/config"
	"oosa_rewild/internal/models"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func NotificationsCreate(c *gin.Context, notifCode string, notifType string, userId primitive.ObjectID, message string, identifier primitive.ObjectID) {
	userDetail := GetAuthUser(c)
	insert := models.Notifications{
		NotificationsCode:       notifCode,
		NotificationsType:       notifType,
		NotificationsUser:       userId,
		NotificationsMessage:    message,
		NotificationsIdentifier: identifier,
		NotificationsCreatedAt:  primitive.NewDateTimeFromTime(time.Now()),
		NotificationsCreatedBy:  userDetail.UsersId,
	}
	config.DB.Collection("Notifications").InsertOne(context.TODO(), insert)
}
