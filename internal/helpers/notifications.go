package helpers

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"oosa_rewild/internal/config"
	"oosa_rewild/internal/models"
	"reflect"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
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

func NotificationAddToContext(c *gin.Context, from primitive.ObjectID, event string, to primitive.ObjectID, data map[string]interface{}) {
	userDocFrom, errfrom := findUserSourceId(from)
	userDocTo, errto := findUserSourceId(to)
	if errfrom != nil {
		return
	}
	if errto != nil {
		return
	}
	newNotifPayload := map[string]interface{}{
		"from":  userDocFrom.UsersSourceId,
		"event": event,
		"to":    []string{userDocTo.UsersSourceId},
		"data":  data,
	}

	existing, exists := c.Get(config.APP.NotificationContextKey)
	if !exists {
		c.Set(config.APP.NotificationContextKey, newNotifPayload)
		return
	}

	switch notif := existing.(type) {
	case []interface{}:
		found := false
		for i, n := range notif {
			if existingNotif, ok := n.(map[string]interface{}); ok {
				if isSameNotification(existingNotif, newNotifPayload) {
					mergeToField(existingNotif, newNotifPayload)
					notif[i] = existingNotif
					found = true
					break
				}
			}
		}
		if !found {
			notif = append(notif, newNotifPayload)
		}
		c.Set(config.APP.NotificationContextKey, notif)
	case map[string]interface{}:
		if isSameNotification(notif, newNotifPayload) {
			mergeToField(notif, newNotifPayload)
			c.Set(config.APP.NotificationContextKey, notif)
		} else {
			c.Set(config.APP.NotificationContextKey, []interface{}{notif, newNotifPayload})
		}
	default:
		c.Set(config.APP.NotificationContextKey, newNotifPayload)
	}
}

func NotificationWriteHeader(c *gin.Context) {
	if notif, exists := c.Get(config.APP.NotificationContextKey); exists {
		var notifSlice []interface{}
		switch n := notif.(type) {
		case []interface{}:
			notifSlice = n
		default:
			notifSlice = []interface{}{n}
		}

		jsonData, err := json.Marshal(notifSlice)
		if err != nil {
			fmt.Println("ERROR marshaling notification:", err)
			return
		}

		encoded := base64.StdEncoding.EncodeToString(jsonData)

		c.Writer.Header().Set(config.APP.NotificationHeaderName, encoded)
	}
}

func findUserSourceId(userId primitive.ObjectID) (*models.Users, error) {
	collection := config.DB.Collection("Users")

	var userDoc models.Users
	err := collection.FindOne(context.TODO(), bson.M{"_id": userId}).Decode(&userDoc)
	if err != nil {
		return nil, err
	}
	return &userDoc, nil
}

func isSameNotification(a, b map[string]interface{}) bool {
	if a["from"] != b["from"] {
		return false
	}
	if a["event"] != b["event"] {
		return false
	}
	return reflect.DeepEqual(a["data"], b["data"])
}

func mergeToField(existingNotif, newNotif map[string]interface{}) {
	var existingTo []string
	if toVal, ok := existingNotif["to"]; ok {
		if arr, ok2 := toVal.([]string); ok2 {
			existingTo = arr
		} else {
			existingTo = []string{}
		}
	}
	var newTo []string
	if toVal, ok := newNotif["to"]; ok {
		if arr, ok2 := toVal.([]string); ok2 {
			newTo = arr
		} else {
			newTo = []string{}
		}
	}
	for _, recipient := range newTo {
		if !stringInSlice(recipient, existingTo) {
			existingTo = append(existingTo, recipient)
		}
	}
	existingNotif["to"] = existingTo
}

func stringInSlice(s string, list []string) bool {
	for _, v := range list {
		if v == s {
			return true
		}
	}
	return false
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
