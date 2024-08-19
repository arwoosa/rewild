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
	"go.mongodb.org/mongo-driver/mongo"
)

type EventParticipantsRepository struct{}
type EventParticipantsRequest struct {
	EventParticipantsUser string `json:"event_participants_user" validate:"required"`
}

func GetEventParticipantStatus(status string) int64 {
	ParticipantStatus := map[string]int64{
		"PENDING":  0,
		"ACCEPTED": 1,
		"REJECTED": 2,
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
	c.JSON(http.StatusOK, results)
}

func (r EventParticipantsRepository) Create(c *gin.Context) {
	err := EventRepository{}.ReadOne(c, &models.Events{})
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
	invitedUser := helpers.StringToPrimitiveObjId(payload.EventParticipantsUser)

	var EventParticipantsCheck models.EventParticipants
	checkParticipant := config.DB.Collection("EventParticipants").FindOne(context.TODO(), bson.D{{Key: "event_participants_event", Value: eventId}, {Key: "event_participants_user", Value: invitedUser}}).Decode(&EventParticipantsCheck)

	if checkParticipant != mongo.ErrNoDocuments {
		helpers.ResultMessageSuccess(c, "This user has already been invited")
		return
	}

	var EventInvitationCheck models.Users
	checkUser := config.DB.Collection("Users").FindOne(context.TODO(), bson.D{{Key: "_id", Value: invitedUser}}).Decode(&EventInvitationCheck)

	if checkUser == mongo.ErrNoDocuments {
		helpers.ResultMessageSuccess(c, "User not found")
		return
	}

	insert := models.EventParticipants{
		EventParticipantsEvent:     eventId,
		EventParticipantsUser:      invitedUser,
		EventParticipantsStatus:    GetEventParticipantStatus("PENDING"),
		EventParticipantsCreatedBy: userDetail.UsersId,
		EventParticipantsCreatedAt: primitive.NewDateTimeFromTime(time.Now()),
	}

	result, err := config.DB.Collection("EventParticipants").InsertOne(context.TODO(), insert)
	if err != nil {
		fmt.Println("ERROR", err.Error())
		return
	}

	var EventParticipants models.EventParticipants
	config.DB.Collection("EventParticipants").FindOne(context.TODO(), bson.D{{Key: "_id", Value: result.InsertedID}}).Decode(&EventParticipants)

	NotificationMessage := models.NotificationMessage{
		Message: "你有一張來自{0}的邀請函",
		Data:    []map[string]interface{}{helpers.NotificationFormatUser(userDetail)},
	}
	helpers.NotificationsCreate(c, helpers.NOTIFICATION_INVITATION, invitedUser, NotificationMessage, EventParticipants.EventParticipantsId)
	c.JSON(http.StatusOK, EventParticipants)
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
