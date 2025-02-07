package repository

import (
	"context"
	"net/http"
	"oosa_rewild/internal/config"
	"oosa_rewild/internal/helpers"
	"oosa_rewild/internal/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type EventInvitationMessageRepository struct{}

type EventInvitationMessageRequest struct {
	EventsInvitationMessage  string `json:"events_invitation_message" validate:"required"`
	EventsInvitationTemplate string `json:"events_invitation_template" validate:"required"`
}
type EventJoinMessageRequest struct {
	EventParticipantsRequestMessage string `json:"event_participants_request_message" validate:"required"`
}

func (r EventInvitationMessageRepository) Update(c *gin.Context) {
	var Event models.Events
	err := EventRepository{}.ReadOne(c, &Event)
	if err != nil {
		return
	}

	//userDetail := helpers.GetAuthUser(c)
	var payload EventInvitationMessageRequest
	validateError := helpers.Validate(c, &payload)
	if validateError != nil {
		return
	}

	match, errMessage := helpers.ValidateStringLength(payload.EventsInvitationMessage, int(config.APP_LIMIT.LengthEventInvitationMessage))
	if !match {
		helpers.ResponseBadRequestError(c, "Message can only contain "+errMessage)
		return
	}

	Event.EventsInvitationTemplate = payload.EventsInvitationTemplate
	Event.EventsInvitationMessage = payload.EventsInvitationMessage

	filters := bson.D{{Key: "_id", Value: Event.EventsId}}
	upd := bson.D{{Key: "$set", Value: Event}}
	config.DB.Collection("Events").UpdateOne(context.TODO(), filters, upd)
	c.JSON(http.StatusOK, Event)
}

func (r EventInvitationMessageRepository) Join(c *gin.Context) {
	userDetail := helpers.GetAuthUser(c)
	var payload EventJoinMessageRequest
	validateError := helpers.Validate(c, &payload)
	if validateError != nil {
		return
	}

	var Events models.Events
	var EventParticipants models.EventParticipants
	id, _ := primitive.ObjectIDFromHex(c.Param("id"))
	filter := bson.D{{Key: "_id", Value: id}}
	err := config.DB.Collection("Events").FindOne(context.TODO(), filter).Decode(&Events)
	if err != nil {
		helpers.ResponseError(c, "Invalid event")
	}

	match, errMessage := helpers.ValidateStringLength(payload.EventParticipantsRequestMessage, int(config.APP_LIMIT.LengthEventParticipantMessage))
	if !match {
		helpers.ResponseBadRequestError(c, errMessage)
		return
	}

	countFilter := bson.D{{Key: "event_participants_event", Value: id}}
	opts := options.Count().SetHint("_id_")
	countParticipant, _ := config.DB.Collection("EventParticipants").CountDocuments(context.TODO(), countFilter, opts)

	status := GetEventParticipantStatus("ACCEPTED")
	if *Events.EventsRequiresApproval == 1 {
		status = GetEventParticipantStatus("APPLIED")
	}

	if *Events.EventsParticipantLimit != 0 && *Events.EventsParticipantLimit < int(countParticipant+1) {
		eventParticipantLimit := strconv.Itoa(*Events.EventsParticipantLimit)
		helpers.ResponseBadRequestError(c, "This event can only be attended by "+eventParticipantLimit+" participants")
		return
	}

	// Check if user already in this event
	checkParticipantFilter := bson.D{
		{Key: "event_participants_event", Value: id},
		{Key: "event_participants_user", Value: userDetail.UsersId},
	}
	checkParticipantErr := config.DB.Collection("EventParticipants").FindOne(context.TODO(), checkParticipantFilter).Decode(&EventParticipants)
	if checkParticipantErr == nil {
		helpers.ResponseError(c, "You are already in this event")
		return
	}

	insertParticipant := models.EventParticipants{
		EventParticipantsEvent:          id,
		EventParticipantsUser:           userDetail.UsersId,
		EventParticipantsStatus:         status,
		EventParticipantsCreatedBy:      userDetail.UsersId,
		EventParticipantsCreatedAt:      primitive.NewDateTimeFromTime(time.Now()),
		EventParticipantsRequestMessage: payload.EventParticipantsRequestMessage,
	}

	_, inserParticipantErr := config.DB.Collection("EventParticipants").InsertOne(context.TODO(), insertParticipant)

	if inserParticipantErr != nil {
		helpers.ResponseError(c, err.Error())
		return
	}

	if status == GetEventParticipantStatus("ACCEPTED") {
		EventRepository{}.HandleBadges(c, id)
	}

	// 補充缺失功能 建立通知函式，通知訊息給活動主持人，讓主持人知道有使用者申請加入活動
	NotificationMessage := models.NotificationMessage{
		Message: "{0}提出加入{0}的申請!",
		Data: []map[string]interface{}{
			helpers.NotificationFormatUser(userDetail),
			helpers.NotificationFormatEvent(Events),
		},
	}
	helpers.NotificationsCreate(c, helpers.NOTIFICATION_JOINING_REQUEST, userDetail.UsersId, NotificationMessage, Events.EventsCreatedBy)

	Data := map[string]interface{}{
		"events_name": Events.EventsName,
	}
	helpers.NotificationAddToContext(c, Events.EventsCreatedBy, helpers.NOTIFICATION_JOINING_REQUEST, Events.EventsCreatedBy, Data)

	helpers.ResponseSuccessMessage(c, "Join request for event submitted")
}
