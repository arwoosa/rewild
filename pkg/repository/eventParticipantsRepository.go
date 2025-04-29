package repository

import (
	"context"
	"fmt"
	"net/http"
	"oosa_rewild/internal/config"
	"oosa_rewild/internal/helpers"
	"oosa_rewild/internal/models"
	"strconv"
	"strings"
	"time"

	"github.com/arwoosa/notifaction/helper"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type EventParticipantsRepository struct{}
type EventParticipantsRequest struct {
	EventParticipantsUser []string `json:"event_participants_user" validate:"required"`
}

func GetEventParticipantStatus(status string) int64 {
	ParticipantStatus := map[string]int64{
		"PENDING":  0,
		"ACCEPTED": 1,
		"REJECTED": 2,
		"APPLIED":  3,
	}
	return ParticipantStatus[status]
}

func GetEventParticipantStatusLabel(status int64) string {
	ParticipantStatus := map[int64]string{
		0: "PENDING",
		1: "ACCEPTED",
		2: "REJECTED",
		3: "APPLIED",
	}
	return ParticipantStatus[status]
}

func (r EventParticipantsRepository) Retrieve(c *gin.Context) {
	eventId := helpers.StringToPrimitiveObjId(c.Param("id"))
	err := EventRepository{}.ReadOne(c, &models.Events{})
	if err != nil {
		return
	}

	agg := mongo.Pipeline{
		bson.D{{
			Key: "$match", Value: bson.M{"event_participants_event": eventId},
		}},
		bson.D{{
			Key: "$lookup", Value: bson.M{
				"from":         "Users",
				"localField":   "event_participants_user",
				"foreignField": "_id",
				"as":           "event_participants_user_detail",
			},
		}},
		bson.D{{
			Key: "$unwind", Value: "$event_participants_user_detail",
		}},
	}

	var results []models.EventParticipants
	cursor, err := config.DB.Collection("EventParticipants").Aggregate(context.TODO(), agg)
	cursor.All(context.TODO(), &results)

	if err != nil {
		return
	}

	if len(results) == 0 {
		helpers.ResponseNoData(c, "No Data")
		return
	}

	for k, v := range results {
		results[k].EventParticipantsStatusLabel = GetEventParticipantStatusLabel(v.EventParticipantsStatus)
	}
	c.JSON(http.StatusOK, results)
}

func (r EventParticipantsRepository) Create(c *gin.Context) {
	var Event models.Events
	err := EventRepository{}.ReadOne(c, &Event)
	if err != nil {
		return
	}

	userDetail := helpers.GetAuthUser(c)
	var payload EventParticipantsRequest
	validateError := helpers.Validate(c, &payload)
	if validateError != nil {
		return
	}

	eventId := helpers.StringToPrimitiveObjId(c.Param("id"))

	if Event.EventsParticipantLimit != nil && *Event.EventsParticipantLimit != 0 {
		countFilter := bson.D{
			{Key: "event_participants_event", Value: eventId},
			{Key: "event_participants_status", Value: GetEventParticipantStatus("ACCEPTED")},
		}
		acceptedCount, err := config.DB.Collection("EventParticipants").CountDocuments(context.TODO(), countFilter)
		if err != nil {
			helpers.ResponseError(c, "Error checking accepted participants count")
			return
		}
		// 若加上這次邀請的數量會超過上限，則回傳錯誤
		if int(acceptedCount)+len(payload.EventParticipantsUser) > *Event.EventsParticipantLimit {
			limitMsg := "This event can only be attended by " + strconv.Itoa(*Event.EventsParticipantLimit) + " participants"
			helpers.ResponseBadRequestError(c, limitMsg)
			return
		}
	}

	var invitedUserMsg []string
	invitedUserId := make([]primitive.ObjectID, 0)
	var foundEvent models.Events
	config.DB.Collection("Events").FindOne(context.TODO(), bson.M{"_id": eventId}).Decode(&foundEvent)
	notifyData := map[string]string{
		"events_name": foundEvent.EventsName,
	}
	var notifyMsg helper.NotifyMsg
	for _, v := range payload.EventParticipantsUser {
		canInvite := true
		uID := v
		invitedUser := helpers.StringToPrimitiveObjId(uID)

		var EventParticipantsCheck models.EventParticipants
		checkParticipant := config.DB.Collection("EventParticipants").FindOne(context.TODO(), bson.D{{Key: "event_participants_event", Value: eventId}, {Key: "event_participants_user", Value: invitedUser}}).Decode(&EventParticipantsCheck)

		if checkParticipant != mongo.ErrNoDocuments {
			invitedUserMsg = append(invitedUserMsg, uID+" has already been invited")
			canInvite = false
		}

		var EventInvitationCheck models.Users
		checkUser := config.DB.Collection("Users").FindOne(context.TODO(), bson.D{{Key: "_id", Value: invitedUser}}).Decode(&EventInvitationCheck)

		if checkUser == mongo.ErrNoDocuments {
			invitedUserMsg = append(invitedUserMsg, uID+" not found")
			canInvite = false
		}

		if canInvite {
			invitedUserId = append(invitedUserId, invitedUser)
			insert := models.EventParticipants{
				EventParticipantsEvent:     eventId,
				EventParticipantsUser:      invitedUser,
				EventParticipantsStatus:    GetEventParticipantStatus("PENDING"),
				EventParticipantsCreatedBy: userDetail.UsersId,
				EventParticipantsCreatedAt: primitive.NewDateTimeFromTime(time.Now()),
				EventParticipantsInvitation: &models.EventParticipantInvitation{
					InvitationMessage:  Event.EventsInvitationMessage,
					InvitationTemplate: Event.EventsInvitationTemplate,
				},
			}

			result, err := config.DB.Collection("EventParticipants").InsertOne(context.TODO(), insert)
			if err != nil {
				fmt.Println("ERROR", err.Error())
				return
			}

			var EventParticipants models.EventParticipants
			config.DB.Collection("EventParticipants").FindOne(context.TODO(), bson.D{{Key: "_id", Value: result.InsertedID}}).Decode(&EventParticipants)

			EventParticipants.EventParticipantsStatusLabel = GetEventParticipantStatusLabel(EventParticipants.EventParticipantsStatus)

			NotificationMessage := models.NotificationMessage{
				Message: "你有一張來自{0}的邀請函",
				Data:    []map[string]interface{}{helpers.NotificationFormatUser(userDetail)},
			}
			helpers.NotificationsCreate(c, helpers.NOTIFICATION_INVITATION, invitedUser, NotificationMessage, EventParticipants.EventParticipantsId)

			if notifyMsg == nil {
				notifyMsg, err = helper.NewNotifyMsg(
					helpers.NOTIFICATION_INVITATION,
					userDetail.UsersId, invitedUser,
					notifyData, helpers.FindUserSourceId)
				if err != nil {
					fmt.Println("new notify msg err: " + err.Error())
				}
			} else {
				notifyMsg.AddTo(EventParticipants.EventParticipantsId)
			}
			EventRepository{}.HandleParticipantFriend(c, eventId)
		}
	}

	if notifyMsg != nil {
		notifyMsg.WriteToHeader(c, config.APP.NotificationHeaderName)
	}

	var EventParticipants []models.EventParticipants
	eventParticipantFilter := bson.D{
		{Key: "event_participants_event", Value: eventId},
		{Key: "event_participants_user", Value: bson.M{"$in": invitedUserId}},
	}
	participantsCursor, _ := config.DB.Collection("EventParticipants").Find(context.TODO(), eventParticipantFilter)
	participantsCursor.All(context.TODO(), &EventParticipants)
	if len(invitedUserId) > 0 {
		c.JSON(http.StatusOK, EventParticipants)
	} else {
		helpers.ResponseError(c, strings.Join(invitedUserMsg, ", "))
	}
}

func (r EventParticipantsRepository) ReadOne(c *gin.Context, EventParticipants *models.EventParticipants) error {
	eventId := helpers.StringToPrimitiveObjId(c.Param("id"))
	participantId := helpers.StringToPrimitiveObjId(c.Param("participantId"))
	filter := bson.D{{Key: "_id", Value: participantId}, {Key: "event_participants_event", Value: eventId}}
	err := config.DB.Collection("EventParticipants").FindOne(context.TODO(), filter).Decode(&EventParticipants)
	if err != nil {
		helpers.ResultEmpty(c, err)
	}
	return err
}

func (r EventParticipantsRepository) Delete(c *gin.Context) {
	err := EventRepository{}.ReadOne(c, &models.Events{})
	if err != nil {
		return
	}

	var EventParticipants models.EventParticipants
	errMb := r.ReadOne(c, &EventParticipants)
	if errMb == nil {
		filters := bson.D{{Key: "_id", Value: EventParticipants.EventParticipantsId}}
		config.DB.Collection("EventParticipants").DeleteOne(context.TODO(), filters)
		helpers.ResultMessageSuccess(c, "User removed from event")
	}
}

func (r EventParticipantsRepository) Read(c *gin.Context) {
	eventId := helpers.StringToPrimitiveObjId(c.Param("id"))

	var eventParticipant models.EventParticipants
	err := r.ReadOne(c, &eventParticipant)
	if err != nil {
		return
	}

	// 獲取事件詳細信息
	var event models.Events
	eventFilter := bson.D{{Key: "_id", Value: eventId}}
	errEvent := config.DB.Collection("Events").FindOne(context.TODO(), eventFilter).Decode(&event)
	if errEvent != nil {
		helpers.ResponseError(c, "Event not found")
		return
	}

	// 獲取參與者詳細信息
	var participant models.Users
	participantFilter := bson.D{{Key: "_id", Value: eventParticipant.EventParticipantsUser}}
	errUser := config.DB.Collection("Users").FindOne(context.TODO(), participantFilter).Decode(&participant)
	if errUser != nil {
		helpers.ResponseError(c, "Participant not found")
		return
	}

	// 計算當前參與者數量
	activeParticipantsFilter := bson.D{
		{Key: "event_participants_event", Value: eventId},
		{Key: "event_participants_status", Value: GetEventParticipantStatus("ACCEPTED")},
	}
	currentCount, _ := config.DB.Collection("EventParticipants").CountDocuments(context.TODO(), activeParticipantsFilter)

	// 構建響應
	response := gin.H{
		"id": eventParticipant.EventParticipantsId.Hex(),
		"event": gin.H{
			"id":            event.EventsId.Hex(),
			"name":          event.EventsName,
			"date":          event.EventsDate.Time().Format("2006-01-02"),
			"current_count": strconv.FormatInt(currentCount, 10),
			"max_count": func() string {
				if event.EventsParticipantLimit != nil {
					return strconv.Itoa(*event.EventsParticipantLimit)
				}
				return "Unlimited"
			}(),
		},
		"participant": gin.H{
			"id":     participant.UsersId.Hex(),
			"name":   participant.UsersName,
			"avatar": participant.UsersAvatar,
		},
		"request_message": func() string {
			if eventParticipant.EventParticipantsInvitation != nil {
				return eventParticipant.EventParticipantsInvitation.InvitationMessage
			}
			return ""
		}(),
		"create_at": eventParticipant.EventParticipantsCreatedAt.Time().Format("2006-01-02 15:04:05"),
	}

	c.JSON(http.StatusOK, response)
}

func (r EventParticipantsRepository) ActiveParticipants(eventId primitive.ObjectID) []models.EventParticipants {
	var ActiveParticipants []models.EventParticipants
	activeParticipantsFilter := bson.D{
		{Key: "event_participants_event", Value: eventId},
		{Key: "event_participants_status", Value: GetEventParticipantStatus("ACCEPTED")},
	}
	cursor, _ := config.DB.Collection("EventParticipants").Find(context.TODO(), activeParticipantsFilter)
	cursor.All(context.TODO(), &ActiveParticipants)
	return ActiveParticipants
}
