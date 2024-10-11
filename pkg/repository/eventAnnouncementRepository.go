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

type EventAnnouncementRepository struct{}
type EventAnnouncementRequest struct {
	EventMessageBoardAnnouncement string `json:"event_message_board_announcement" validate:"required"`
	EventMessageBoardCategory     string `json:"event_message_board_category" validate:"required"`
	EventMessageBoardIsPinned     int    `json:"event_message_board_is_pinned"`
}

func (r EventAnnouncementRepository) Retrieve(c *gin.Context) {
	messageCategory := c.Query("category")
	eventId := helpers.StringToPrimitiveObjId(c.Param("id"))
	err := EventRepository{}.ReadOne(c, &models.Events{})
	if err != nil {
		return
	}

	match := bson.D{
		{Key: "event_message_board_event", Value: eventId},
		{Key: "event_message_board_announcement", Value: bson.M{"$exists": true}},
	}

	if messageCategory != "" {
		match = append(match, bson.E{Key: "event_message_board_category", Value: messageCategory})
	}

	criteria := bson.D{{
		Key: "$match", Value: match,
	}}

	agg := mongo.Pipeline{
		criteria,
		bson.D{{
			Key: "$lookup", Value: bson.M{
				"from":         "Users",
				"localField":   "event_message_board_created_by",
				"foreignField": "_id",
				"as":           "event_message_board_created_by_user",
			},
		}},
		bson.D{{
			Key: "$unwind", Value: "$event_message_board_created_by_user",
		}},
	}

	var results []models.EventMessageBoard
	cursor, err := config.DB.Collection("EventMessageBoard").Aggregate(context.TODO(), agg)
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

func (r EventAnnouncementRepository) Create(c *gin.Context) {
	err := EventRepository{}.ReadOne(c, &models.Events{})
	if err != nil {
		return
	}

	userDetail := helpers.GetAuthUser(c)
	var payload EventAnnouncementRequest
	validateError := helpers.Validate(c, &payload)
	if validateError != nil {
		return
	}

	match, errMessage := helpers.ValidateStringLength(payload.EventMessageBoardAnnouncement, int(config.APP_LIMIT.LengthEventMessageBoardMessage))
	if !match {
		helpers.ResponseBadRequestError(c, "Announcement can only contain "+errMessage)
		return
	}

	insert := models.EventMessageBoard{
		EventMessageBoardEvent: helpers.StringToPrimitiveObjId(c.Param("id")),
		// EventMessageBoardStatus
		// EventMessageBoardCategory
		EventMessageBoardCreatedBy: userDetail.UsersId,
		EventMessageBoardCreatedAt: primitive.NewDateTimeFromTime(time.Now()),
	}

	r.ProcessData(c, &insert, payload)

	result, err := config.DB.Collection("EventMessageBoard").InsertOne(context.TODO(), insert)
	if err != nil {
		fmt.Println("ERROR", err.Error())
		return
	}

	var EventMessageBoard models.EventMessageBoard
	config.DB.Collection("EventMessageBoard").FindOne(context.TODO(), bson.D{{Key: "_id", Value: result.InsertedID}}).Decode(&EventMessageBoard)
	c.JSON(http.StatusOK, EventMessageBoard)
}

func (r EventAnnouncementRepository) Read(c *gin.Context) {
	err := EventRepository{}.ReadOne(c, &models.Events{})
	if err != nil {
		return
	}

	var EventMessageBoard models.EventMessageBoard
	errMb := r.ReadOne(c, &EventMessageBoard)
	if errMb == nil {
		c.JSON(http.StatusOK, EventMessageBoard)
	}
}

func (r EventAnnouncementRepository) ReadOne(c *gin.Context, EventMessageBoard *models.EventMessageBoard) error {
	eventId := helpers.StringToPrimitiveObjId(c.Param("id"))
	eventMessageBoardId := helpers.StringToPrimitiveObjId(c.Param("messageBoardId"))
	filter := bson.D{{Key: "_id", Value: eventMessageBoardId}, {Key: "event_message_board_event", Value: eventId}}
	err := config.DB.Collection("EventMessageBoard").FindOne(context.TODO(), filter).Decode(&EventMessageBoard)
	if err != nil {
		helpers.ResultEmpty(c, err)
	}
	return err
}

func (r EventAnnouncementRepository) Update(c *gin.Context) {
	err := EventRepository{}.ReadOne(c, &models.Events{})
	if err != nil {
		return
	}

	var payload EventAnnouncementRequest
	validateError := helpers.Validate(c, &payload)
	if validateError != nil {
		return
	}

	match, errMessage := helpers.ValidateStringLength(payload.EventMessageBoardAnnouncement, int(config.APP_LIMIT.LengthEventMessageBoardMessage))
	if !match {
		helpers.ResponseBadRequestError(c, "Announcement can only contain "+errMessage)
		return
	}

	var EventMessageBoard models.EventMessageBoard
	errMb := r.ReadOne(c, &EventMessageBoard)
	if errMb == nil {
		r.ProcessData(c, &EventMessageBoard, payload)
		filters := bson.D{{Key: "_id", Value: EventMessageBoard.EventMessageBoardId}, {Key: "event_message_board_event", Value: EventMessageBoard.EventMessageBoardEvent}}
		upd := bson.D{{Key: "$set", Value: EventMessageBoard}}
		config.DB.Collection("EventMessageBoard").UpdateOne(context.TODO(), filters, upd)
		c.JSON(http.StatusOK, EventMessageBoard)
	}
}

func (r EventAnnouncementRepository) Delete(c *gin.Context) {
	err := EventRepository{}.ReadOne(c, &models.Events{})
	if err != nil {
		return
	}

	var EventMessageBoard models.EventMessageBoard
	errMb := r.ReadOne(c, &EventMessageBoard)
	if errMb == nil {
		filters := bson.D{{Key: "_id", Value: EventMessageBoard.EventMessageBoardId}}
		config.DB.Collection("EventMessageBoard").DeleteOne(context.TODO(), filters)
		helpers.ResultMessageSuccess(c, "Message board record deleted")
	}
}

func (r EventAnnouncementRepository) ProcessData(c *gin.Context, EventMessageBoard *models.EventMessageBoard, payload EventAnnouncementRequest) {
	EventMessageBoard.EventMessageBoardAnnouncement = payload.EventMessageBoardAnnouncement
	EventMessageBoard.EventMessageBoardCategory = payload.EventMessageBoardCategory

	isPinned := 0
	if payload.EventMessageBoardIsPinned > 0 {
		isPinned = payload.EventMessageBoardIsPinned
	}
	EventMessageBoard.EventMessageBoardIsPinned = &isPinned
}
