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
)

type EventScheduleRepository struct{}
type EventScheduleRequest struct {
	EventSchedulesDatetime    string `json:"event_schedules_datetime" validate:"required,datetime=2006-01-02 15:04:05"`
	EventSchedulesDescription string `json:"event_schedules_description" validate:"required"`
	// #TODO: Datetime: RFC-3339 2024
	// RFC3339     = "2006-01-02T15:04:05Z07:00" 2024-03-24T00:00:00+08:00
}

func (r EventScheduleRepository) Retrieve(c *gin.Context) {
	err := EventRepository{}.ReadOne(c, &models.Events{})
	if err != nil {
		return
	}

	var results []models.EventSchedules
	cursor, err := config.DB.Collection("EventSchedules").Find(context.TODO(), bson.D{})
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

func (r EventScheduleRepository) Create(c *gin.Context) {
	err := EventRepository{}.ReadOne(c, &models.Events{})
	if err != nil {
		return
	}

	userDetail := helpers.GetAuthUser(c)
	var payload EventScheduleRequest
	validateError := helpers.Validate(c, &payload)
	if validateError != nil {
		return
	}

	insert := models.EventSchedules{
		EventSchedulesEvent:     helpers.StringToPrimitiveObjId(c.Param("id")),
		EventSchedulesCreatedBy: userDetail.UsersId,
		EventSchedulesCreatedAt: primitive.NewDateTimeFromTime(time.Now()),
	}

	r.ProcessData(c, &insert, payload)

	result, err := config.DB.Collection("EventSchedules").InsertOne(context.TODO(), insert)
	if err != nil {
		fmt.Println("ERROR", err.Error())
		return
	}

	var EventSchedules models.EventSchedules
	config.DB.Collection("EventSchedules").FindOne(context.TODO(), bson.D{{Key: "_id", Value: result.InsertedID}}).Decode(&EventSchedules)
	c.JSON(http.StatusOK, EventSchedules)
}

func (r EventScheduleRepository) Read(c *gin.Context) {
	err := EventRepository{}.ReadOne(c, &models.Events{})
	if err != nil {
		return
	}

	var EventSchedules models.EventSchedules
	errMb := r.ReadOne(c, &EventSchedules)
	if errMb == nil {
		c.JSON(http.StatusOK, EventSchedules)
	}
}

func (r EventScheduleRepository) ReadOne(c *gin.Context, EventSchedules *models.EventSchedules) error {
	eventId := helpers.StringToPrimitiveObjId(c.Param("id"))
	EventSchedulesId := helpers.StringToPrimitiveObjId(c.Param("scheduleId"))
	filter := bson.D{{Key: "_id", Value: EventSchedulesId}, {Key: "event_schedules_event", Value: eventId}}
	err := config.DB.Collection("EventSchedules").FindOne(context.TODO(), filter).Decode(&EventSchedules)
	if err != nil {
		helpers.ResultEmpty(c, err)
	}
	return err
}

func (r EventScheduleRepository) Update(c *gin.Context) {
	err := EventRepository{}.ReadOne(c, &models.Events{})
	if err != nil {
		return
	}

	var payload EventScheduleRequest
	validateError := helpers.Validate(c, &payload)
	if validateError != nil {
		return
	}

	var EventSchedules models.EventSchedules
	errMb := r.ReadOne(c, &EventSchedules)
	if errMb == nil {
		r.ProcessData(c, &EventSchedules, payload)
		filters := bson.D{{Key: "_id", Value: EventSchedules.EventSchedulesId}, {Key: "event_schedules_event", Value: EventSchedules.EventSchedulesEvent}}
		upd := bson.D{{Key: "$set", Value: EventSchedules}}
		config.DB.Collection("EventSchedules").UpdateOne(context.TODO(), filters, upd)
		c.JSON(http.StatusOK, EventSchedules)
	}
}

func (r EventScheduleRepository) Delete(c *gin.Context) {
	err := EventRepository{}.ReadOne(c, &models.Events{})
	if err != nil {
		return
	}

	var EventSchedules models.EventSchedules
	errMb := r.ReadOne(c, &EventSchedules)
	if errMb == nil {
		filters := bson.D{{Key: "_id", Value: EventSchedules.EventSchedulesId}}
		config.DB.Collection("EventSchedules").DeleteOne(context.TODO(), filters)
		helpers.ResultMessageSuccess(c, "Schedule deleted")
	}
}

func (r EventScheduleRepository) ProcessData(c *gin.Context, EventSchedules *models.EventSchedules, payload EventScheduleRequest) {
	time := helpers.StringToDateTime(payload.EventSchedulesDatetime)
	EventSchedules.EventSchedulesDatetime = primitive.NewDateTimeFromTime(time)
	EventSchedules.EventSchedulesDescription = payload.EventSchedulesDescription
}
