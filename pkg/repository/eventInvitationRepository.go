package repository

import (
	"context"
	"fmt"
	"net/http"
	"oosa_rewild/internal/config"
	"oosa_rewild/internal/helpers"
	"oosa_rewild/internal/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type EventInvitationRepository struct{}
type EventInvitationRequest struct {
	EventParticipantsStatus int64 `json:"event_participants_status" validate:"required"`
}

func (r EventInvitationRepository) Retrieve(c *gin.Context) {
	userDetail := helpers.GetAuthUser(c)
	fmt.Print("EventInvitationRepository: Retrieve")
	var results []models.EventParticipants
	filter := bson.D{
		{Key: "event_participants_user", Value: userDetail.UsersId},
		{Key: "event_participants_status", Value: GetEventParticipantStatus("PENDING")},
	}

	cursor, err := config.DB.Collection("EventParticipants").Find(context.TODO(), filter)
	cursor.All(context.TODO(), &results)

	if err != nil {
		return
	}

	if len(results) == 0 {
		helpers.ResponseNoData(c, "No Data")
		return
	}
	c.JSON(http.StatusOK, results)
}

func (r EventInvitationRepository) Read(c *gin.Context) {
	id, _ := primitive.ObjectIDFromHex(c.Param("id"))
	userDetail := helpers.GetAuthUser(c)
	var results models.EventParticipants
	filter := bson.D{
		{Key: "_id", Value: id},
		{Key: "event_participants_user", Value: userDetail.UsersId},
		{Key: "event_participants_status", Value: GetEventParticipantStatus("PENDING")},
	}
	err := config.DB.Collection("EventParticipants").FindOne(context.TODO(), filter).Decode(&results)
	if err != nil {
		helpers.ResultEmpty(c, err)
	}

	c.JSON(http.StatusOK, results)
}

func (r EventInvitationRepository) Update(c *gin.Context) {
	var payload EventInvitationRequest
	validateError := helpers.Validate(c, &payload)
	if validateError != nil {
		return
	}

	id, _ := primitive.ObjectIDFromHex(c.Param("id"))
	userDetail := helpers.GetAuthUser(c)
	var results models.EventParticipants
	filter := bson.D{
		{Key: "_id", Value: id},
		{Key: "event_participants_user", Value: userDetail.UsersId},
		{Key: "event_participants_status", Value: GetEventParticipantStatus("PENDING")},
	}
	err := config.DB.Collection("EventParticipants").FindOne(context.TODO(), filter).Decode(&results)
	if err != nil {
		helpers.ResultEmpty(c, err)
		return
	}

	if results.EventParticipantsStatus == 1 {
		helpers.ResultMessageSuccess(c, "This invitation has already been accepted")
		return
	}

	if results.EventParticipantsStatus == 2 {
		helpers.ResultMessageSuccess(c, "This invitation has already been rejected")
		return
	}

	if payload.EventParticipantsStatus != 1 && payload.EventParticipantsStatus != 2 {
		helpers.ResultMessageError(c, "Unsupported status")
		return
	}

	ActiveParticipants := EventParticipantsRepository{}.ActiveParticipants(results.EventParticipantsEvent)

	var Event models.Events
	eventFilter := bson.D{
		{Key: "_id", Value: results.EventParticipantsEvent},
	}
	config.DB.Collection("Events").FindOne(context.TODO(), eventFilter).Decode(&Event)

	for _, v := range ActiveParticipants {
		NotificationMessage := models.NotificationMessage{
			Message: "{0}有新的夥伴加入!",
			Data:    []map[string]interface{}{helpers.NotificationFormatEvent(Event)},
		}
		helpers.NotificationsCreate(c, helpers.NOTIFICATION_JOINING_NEW, v.EventParticipantsUser, NotificationMessage, results.EventParticipantsEvent)
	}

	results.EventParticipantsStatus = payload.EventParticipantsStatus
	upd := bson.D{{Key: "$set", Value: results}}
	config.DB.Collection("EventParticipants").UpdateOne(context.TODO(), filter, upd)
	EventRepository{}.HandleParticipation(c, userDetail.UsersId, id)
	c.JSON(http.StatusOK, results)
}
