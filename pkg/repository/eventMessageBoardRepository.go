package repository

import (
	"context"
	"fmt"
	"net/http"
	"oosa_rewild/internal/config"
	"oosa_rewild/internal/helpers"
	"oosa_rewild/internal/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type EventMessageBoardRepository struct{}
type EventMessageBoardRequest struct {
	EventMessageBoardBaseMessage string `json:"event_message_board_base_message" validate:"required"`
	EventMessageBoardIsPinned    int    `json:"event_message_board_is_pinned"`
}
type EventMessageBoardPinRequest struct {
	EventMessageBoardCategory string `json:"event_message_board_category" validate:"required"`
}

func (r EventMessageBoardRepository) Retrieve(c *gin.Context) {
	eventId := helpers.StringToPrimitiveObjId(c.Param("id"))
	err := EventRepository{}.ReadOne(c, &models.Events{})
	if err != nil {
		return
	}

	match := bson.D{
		{Key: "event_message_board_event", Value: eventId},
		{Key: "event_message_board_base_message", Value: bson.M{"$exists": true}},
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

func (r EventMessageBoardRepository) Create(c *gin.Context) {
	eventId := helpers.StringToPrimitiveObjId(c.Param("id"))
	var Events models.Events
	err := EventRepository{}.ReadOne(c, &Events)
	if err != nil {
		return
	}

	userDetail := helpers.GetAuthUser(c)
	var payload EventMessageBoardRequest
	validateError := helpers.Validate(c, &payload)
	if validateError != nil {
		return
	}

	currentTime := time.Now()
	if currentTime.After(Events.EventsDateEnd.Time()) {
		helpers.ResponseBadRequestError(c, "This event has ended")
		return
	}

	countFilter := bson.D{{Key: "event_message_board_event", Value: eventId}}
	countMessages, _ := config.DB.Collection("EventMessageBoard").CountDocuments(context.TODO(), countFilter)

	if countMessages+1 > config.APP_LIMIT.EventMessageBoardLimit {
		errMessage := "Maximum number of messages encountered (" + strconv.Itoa(int(config.APP_LIMIT.EventMessageBoardLimit)) + ")"
		helpers.ResponseBadRequestError(c, errMessage)
		return
	}

	match, errMessage := helpers.ValidateStringLength(payload.EventMessageBoardBaseMessage, int(config.APP_LIMIT.LengthEventMessageBoardMessage))
	if !match {
		helpers.ResponseBadRequestError(c, "Message board can only contain "+errMessage)
		return
	}

	insert := models.EventMessageBoard{
		EventMessageBoardEvent: eventId,
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

func (r EventMessageBoardRepository) Read(c *gin.Context) {
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

func (r EventMessageBoardRepository) ReadOne(c *gin.Context, EventMessageBoard *models.EventMessageBoard) error {
	eventId := helpers.StringToPrimitiveObjId(c.Param("id"))
	eventMessageBoardId := helpers.StringToPrimitiveObjId(c.Param("messageBoardId"))
	filter := bson.D{{Key: "_id", Value: eventMessageBoardId}, {Key: "event_message_board_event", Value: eventId}}
	err := config.DB.Collection("EventMessageBoard").FindOne(context.TODO(), filter).Decode(&EventMessageBoard)
	if err != nil {
		helpers.ResultEmpty(c, err)
	}
	return err
}

func (r EventMessageBoardRepository) Update(c *gin.Context) {
	err := EventRepository{}.ReadOne(c, &models.Events{})
	if err != nil {
		return
	}

	var payload EventMessageBoardRequest
	validateError := helpers.Validate(c, &payload)
	if validateError != nil {
		return
	}

	match, errMessage := helpers.ValidateStringLength(payload.EventMessageBoardBaseMessage, int(config.APP_LIMIT.LengthEventMessageBoardMessage))
	if !match {
		helpers.ResponseBadRequestError(c, "Message board can only contain "+errMessage)
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

func (r EventMessageBoardRepository) Delete(c *gin.Context) {
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

func (r EventMessageBoardRepository) ProcessData(c *gin.Context, EventMessageBoard *models.EventMessageBoard, payload EventMessageBoardRequest) {
	EventMessageBoard.EventMessageBoardBaseMessage = payload.EventMessageBoardBaseMessage

	isPinned := 0
	if payload.EventMessageBoardIsPinned > 0 {
		isPinned = payload.EventMessageBoardIsPinned
	}
	EventMessageBoard.EventMessageBoardIsPinned = &isPinned
}

func (r EventMessageBoardRepository) Pin(c *gin.Context) {
	err := EventRepository{}.ReadOne(c, &models.Events{})
	if err != nil {
		return
	}

	userDetail := helpers.GetAuthUser(c)

	var payload EventMessageBoardPinRequest
	validateError := helpers.Validate(c, &payload)
	if validateError != nil {
		return
	}

	var EventMessageBoard models.EventMessageBoard
	errMb := r.ReadOne(c, &EventMessageBoard)
	isPinned := 1
	if errMb == nil {
		insert := models.EventMessageBoard{
			EventMessageBoardEvent: EventMessageBoard.EventMessageBoardEvent,
			// EventMessageBoardBaseMessage: ,
			// EventMessageBoardStatus: ,
			EventMessageBoardCategory:     payload.EventMessageBoardCategory,
			EventMessageBoardAnnouncement: EventMessageBoard.EventMessageBoardBaseMessage,
			EventMessageBoardMessageId:    EventMessageBoard.EventMessageBoardId,
			EventMessageBoardCreatedBy:    userDetail.UsersId,
			EventMessageBoardCreatedAt:    primitive.NewDateTimeFromTime(time.Now()),
			EventMessageBoardIsPinned:     &isPinned,
		}

		result, err := config.DB.Collection("EventMessageBoard").InsertOne(context.TODO(), insert)
		if err != nil {
			fmt.Println("ERROR", err.Error())
			return
		}

		var EventMessageBoard models.EventMessageBoard
		config.DB.Collection("EventMessageBoard").FindOne(context.TODO(), bson.D{{Key: "_id", Value: result.InsertedID}}).Decode(&EventMessageBoard)
		c.JSON(http.StatusOK, EventMessageBoard)
	}

}
