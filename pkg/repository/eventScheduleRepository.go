package repository

import (
	"context"
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

type EventScheduleRepository struct{}

type EventScheduleRequest struct {
	Schedule []EventScheduleRequestItem `json:"schedule" validate:"required"`
}
type EventScheduleRequestItem struct {
	EventSchedulesDatetime    string `json:"event_schedules_datetime" validate:"required,datetime=2006-01-02T15:04:05Z07:00"`
	EventSchedulesDescription string `json:"event_schedules_description" validate:"required"`
	// #TODO: Datetime: RFC-3339 2024
	// RFC3339     = "2006-01-02T15:04:05Z07:00" 2024-03-24T00:00:00+08:00
}

func (r EventScheduleRepository) Retrieve(c *gin.Context) {
	eventId := helpers.StringToPrimitiveObjId(c.Param("id"))
	err := EventRepository{}.ReadOne(c, &models.Events{})
	if err != nil {
		return
	}

	agg := mongo.Pipeline{
		bson.D{{
			Key: "$match", Value: bson.M{"event_schedules_event": eventId},
		}},
		bson.D{{
			Key: "$lookup", Value: bson.M{
				"from":         "Users",
				"localField":   "event_schedules_created_by",
				"foreignField": "_id",
				"as":           "event_schedules_created_by_user",
			},
		}},
		bson.D{{
			Key: "$unwind", Value: "$event_schedules_created_by_user",
		}},
	}

	var results []models.EventSchedules
	cursor, err := config.DB.Collection("EventSchedules").Aggregate(context.TODO(), agg)
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

	eventId := helpers.StringToPrimitiveObjId(c.Param("id"))
	eventSchedule := []interface{}{}

	filters := bson.M{
		"event_schedules_event": eventId,
	}

	config.DB.Collection("EventSchedules").DeleteMany(context.TODO(), filters)

	for _, v := range payload.Schedule {
		scheduleTime := helpers.StringToDateTime(v.EventSchedulesDatetime)

		eventSchedule = append(eventSchedule, models.EventSchedules{
			EventSchedulesEvent:       helpers.StringToPrimitiveObjId(c.Param("id")),
			EventSchedulesCreatedBy:   userDetail.UsersId,
			EventSchedulesCreatedAt:   primitive.NewDateTimeFromTime(time.Now()),
			EventSchedulesDatetime:    primitive.NewDateTimeFromTime(scheduleTime),
			EventSchedulesDescription: v.EventSchedulesDescription,
		})
	}

	config.DB.Collection("EventSchedules").InsertMany(context.TODO(), eventSchedule)

	r.Retrieve(c)
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

	var payload EventScheduleRequestItem
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

func (r EventScheduleRepository) ProcessData(c *gin.Context, EventSchedules *models.EventSchedules, payload EventScheduleRequestItem) {
	time := helpers.StringToDateTime(payload.EventSchedulesDatetime)
	EventSchedules.EventSchedulesDatetime = primitive.NewDateTimeFromTime(time)
	EventSchedules.EventSchedulesDescription = payload.EventSchedulesDescription
}
