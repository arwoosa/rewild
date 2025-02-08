package repository

import (
	"context"
	"fmt"
	"net/http"
	"oosa_rewild/internal/config"
	"oosa_rewild/internal/helpers"
	"oosa_rewild/internal/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type EventInvitationRepository struct{}
type EventInvitationRequest struct {
	EventParticipantsStatus int64 `json:"event_participants_status" validate:"required"`
}

func (r EventInvitationRepository) Retrieve(c *gin.Context) {
	userDetail := helpers.GetAuthUser(c)
	applied := c.Query("applied")
	fmt.Print("EventInvitationRepository: Retrieve", applied)
	var EventParticipants []models.EventParticipants

	matchParam := bson.D{{
		Key: "$match", Value: bson.M{
			"event_participants_user":   userDetail.UsersId,
			"event_participants_status": GetEventParticipantStatus("PENDING"),
		},
	}}

	if applied == "true" {
		matchParam = bson.D{{
			Key: "$match", Value: bson.M{
				"event_participants_status": GetEventParticipantStatus("APPLIED"),
			},
		}}
	}

	agg := mongo.Pipeline{
		matchParam,
		bson.D{{
			Key: "$lookup", Value: bson.M{
				"from":         "Users",
				"localField":   "event_participants_created_by",
				"foreignField": "_id",
				"as":           "event_participants_invited_by",
			},
		}},
		bson.D{{
			Key: "$unwind", Value: "$event_participants_invited_by",
		}},
	}

	if applied == "true" {
		agg = append(agg,
			bson.D{{
				Key: "$lookup", Value: bson.M{
					"from":         "Events",
					"localField":   "event_participants_event",
					"foreignField": "_id",
					"as":           "event_participants_event_detail",
				},
			}},
			bson.D{{
				Key: "$unwind", Value: "$event_participants_event_detail",
			}},
			bson.D{{
				Key: "$match", Value: bson.M{
					"event_participants_event_detail.events_created_by": userDetail.UsersId,
				},
			}})
	}

	cursor, err := config.DB.Collection("EventParticipants").Aggregate(context.TODO(), agg)
	cursor.All(context.TODO(), &EventParticipants)

	if err != nil {
		return
	}

	if len(EventParticipants) == 0 {
		helpers.ResponseNoData(c, "No Data")
		return
	}
	c.JSON(http.StatusOK, EventParticipants)
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
	applied := c.Query("applied")
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

	if applied == "true" {
		filter = bson.D{
			{Key: "_id", Value: id},
			{Key: "event_participants_status", Value: GetEventParticipantStatus("APPLIED")},
		}
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

	var Event models.Events
	eventFilter := bson.D{
		{Key: "_id", Value: results.EventParticipantsEvent},
	}
	if err := config.DB.Collection("Events").FindOne(context.TODO(), eventFilter).Decode(&Event); err != nil {
		helpers.ResponseError(c, "Event not found")
		return
	}

	if payload.EventParticipantsStatus == GetEventParticipantStatus("ACCEPTED") {
		countFilter := bson.D{
			{Key: "event_participants_event", Value: results.EventParticipantsEvent},
			{Key: "event_participants_status", Value: GetEventParticipantStatus("ACCEPTED")},
		}
		acceptedCount, err := config.DB.Collection("EventParticipants").CountDocuments(context.TODO(), countFilter)
		if err != nil {
			helpers.ResponseError(c, "Error checking accepted participants count")
			return
		}
		if *Event.EventsParticipantLimit != 0 && int(acceptedCount+1) > *Event.EventsParticipantLimit {
			limitMsg := "This event can only be attended by " + strconv.Itoa(*Event.EventsParticipantLimit) + " participants"
			helpers.ResponseBadRequestError(c, limitMsg)
			return
		}
	}

	ActiveParticipants := EventParticipantsRepository{}.ActiveParticipants(results.EventParticipantsEvent)

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
	EventRepository{}.HandleBadges(c, id)
	c.JSON(http.StatusOK, results)
}
