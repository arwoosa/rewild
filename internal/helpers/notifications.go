package helpers

import (
	"context"
	"oosa_rewild/internal/config"
	"oosa_rewild/internal/models"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	NOTIFICATION_INVITATION              = "INVITATION"
	NOTIFICATION_JOINING_NEW             = "JOINING_NEW"
	NOTIFICATION_JOINING_REQUEST         = "JOINING_REQUEST"
	NOTIFICATION_FRIEND_REQUEST          = "FRIEND_REQUEST"
	NOTIFICATION_FRIEND_REQUEST_ACCEPTED = "FRIEND_REQUEST_ACCEPTED"
	NOTIFICATION_BADGE_NEW               = "BADGE_NEW"
	NOTIFICATION_STARMAP_NEW             = "STARMAP_NEW"
	NOTIFICATION_STARMAP_PROGRESS        = "STARMAP_PROGRESS"
	NOTIFICATION_UPDATE_POLICY           = "UPDATE_POLICY"
	NOTIFICATION_EVENT_INFO              = "EVENT_INFO"
	NOTIFICATION_EVENT_COUNTDOWN         = "EVENT_COUNTDOWN"
	NOTIFICATION_EVENT_MEMBER_DELETED    = "EVENT_MEMBER_DELETED"
	NOTIFICATION_EVENT_JOIN_DENIED       = "EVENT_JOIN_DENIED"
	NOTIFICATION_COLOG_PHOTO_UPLOADED    = "COLOG_PHOTO_UPLOADED"
	NOTIFICATION_COLOG_REMIND            = "COLOG_REMIND"
)

func NotificationsCreate(c *gin.Context, notifCode string, userId primitive.ObjectID, message models.NotificationMessage, identifier primitive.ObjectID) {
	userDetail := GetAuthUser(c)
	insert := models.Notifications{
		NotificationsCode:       notifCode,
		NotificationsUser:       userId,
		NotificationsMessage:    message,
		NotificationsIdentifier: identifier,
		NotificationsCreatedAt:  primitive.NewDateTimeFromTime(time.Now()),
		NotificationsCreatedBy:  userDetail.UsersId,
	}
	config.DB.Collection("Notifications").InsertOne(context.TODO(), insert)
}

func NotificationFormatEvent(Events models.Events) map[string]any {
	return map[string]any{
		"events_id":   Events.EventsId,
		"events_name": Events.EventsName,
	}
}

func NotificationFormatUser(Users models.Users) map[string]any {
	return map[string]any{
		"users_id":   Users.UsersId,
		"users_name": Users.UsersName,
	}
}

func NotificationFormatBadges(Badges models.Badges) map[string]any {
	return map[string]any{
		"badges_id":   Badges.BadgesId,
		"badges_code": Badges.BadgesCode,
		"badges_name": Badges.BadgesName,
	}
}
