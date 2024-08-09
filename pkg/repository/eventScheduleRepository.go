package repository

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"oosa_rewild/internal/config"
	"oosa_rewild/internal/helpers"
	"oosa_rewild/internal/models"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type EventScheduleRepository struct{}

type EventScheduleRequest struct {
	EventId       string                             `json:"event_id" validate:""`
	EventSchedule []EventScheduleRequestScheduleItem `json:"event_schedule" validate:"required,dive"`
}

type EventScheduleRequestScheduleItem struct {
	EventSchedulesDate string                     `json:"event_schedules_date" validate:"required,datetime=2006-01-02"`
	Schedule           []EventScheduleRequestTime `json:"event_schedules_schedules" validate:"dive"`
	// #TODO: Datetime: RFC-3339 2024
	// RFC3339     = "2006-01-02T15:04:05Z07:00" 2024-03-24T00:00:00+08:00
}
type EventScheduleRequestTime struct {
	EventSchedulesDatetime    string `json:"event_schedules_datetime" validate:"required,datetime=15:04:05"`
	EventSchedulesDescription string `json:"event_schedules_description" validate:"required"`
	// #TODO: Datetime: RFC-3339 2024
	// RFC3339     = "2006-01-02T15:04:05Z07:00" 2024-03-24T00:00:00+08:00
}
type EventScheduleRequestItem struct {
	EventSchedulesDatetime    string `json:"event_schedules_datetime" validate:"required,datetime=2006-01-02T15:04:05Z07:00"`
	EventSchedulesDescription string `json:"event_schedules_description" validate:"required"`
	// #TODO: Datetime: RFC-3339 2024
	// RFC3339     = "2006-01-02T15:04:05Z07:00" 2024-03-24T00:00:00+08:00
}

func (r EventScheduleRepository) Retrieve(c *gin.Context) {
	eventId := helpers.StringToPrimitiveObjId(c.Param("id"))
	var Event models.Events
	Response := EventScheduleRequest{
		EventId:       c.Param("id"),
		EventSchedule: []EventScheduleRequestScheduleItem{},
	}
	config.DB.Collection("Events").FindOne(context.TODO(), bson.D{{Key: "_id", Value: eventId}}).Decode(&Event)

	var results []models.EventSchedules
	filter := bson.D{{Key: "event_schedules_event", Value: eventId}}
	cursor, _ := config.DB.Collection("EventSchedules").Find(context.TODO(), filter)
	cursor.All(context.TODO(), &results)

	days := int(math.Ceil(Event.EventsDateEnd.Time().Sub(Event.EventsDate.Time()).Hours() / 24))

	mappedItem := map[string][]EventScheduleRequestTime{}

	for _, v := range results {
		datetime := v.EventSchedulesDatetime.Time()
		key := datetime.Format("2006-01-02")
		formattedTime := datetime.Format("15:04:05")
		mappedItem[key] = append(mappedItem[key], EventScheduleRequestTime{
			EventSchedulesDatetime:    formattedTime,
			EventSchedulesDescription: v.EventSchedulesDescription,
		})
	}

	for i := 0; i < days; i++ {
		newDate := Event.EventsDate.Time().AddDate(0, 0, i).Format("2006-01-02")

		scheduleItem := mappedItem[newDate]
		if mappedItem[newDate] == nil {
			scheduleItem = make([]EventScheduleRequestTime, 0)
		}

		Response.EventSchedule = append(Response.EventSchedule, EventScheduleRequestScheduleItem{
			EventSchedulesDate: newDate,
			Schedule:           scheduleItem,
		})
	}

	c.JSON(200, Response)
	/*
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
	*/
}

func (r EventScheduleRepository) Create(c *gin.Context) {
	var errList []string
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

	var Event models.Events
	eventId := helpers.StringToPrimitiveObjId(c.Param("id"))
	eventSchedule := []interface{}{}

	filters := bson.M{
		"event_schedules_event": eventId,
	}
	config.DB.Collection("EventSchedules").DeleteMany(context.TODO(), filters)
	config.DB.Collection("Events").FindOne(context.TODO(), bson.D{{Key: "_id", Value: eventId}}).Decode(&Event)

	eventDateStart := Event.EventsDate.Time()
	eventDateEnd := Event.EventsDateEnd.Time()

	for _, v := range payload.EventSchedule {
		scheduleDate := v.EventSchedulesDate

		for _, vSchedule := range v.Schedule {
			fullDatetime := scheduleDate + "T" + vSchedule.EventSchedulesDatetime + "Z"
			fmt.Println(fullDatetime)
			scheduleTime := helpers.StringToDateTime(fullDatetime)

			isValidDate := helpers.TimeIsBetween(scheduleTime, eventDateStart, eventDateEnd)

			if !isValidDate {
				errList = append(errList, scheduleDate+" "+vSchedule.EventSchedulesDatetime)
			}
			fmt.Println("isValidDate: ", isValidDate)
			fmt.Println("scheduleTime: ", scheduleTime)
			fmt.Println("eventDateStart: ", eventDateStart)
			fmt.Println("eventDateEnd: ", eventDateEnd)

			eventSchedule = append(eventSchedule, models.EventSchedules{
				EventSchedulesEvent:       helpers.StringToPrimitiveObjId(c.Param("id")),
				EventSchedulesCreatedBy:   userDetail.UsersId,
				EventSchedulesCreatedAt:   primitive.NewDateTimeFromTime(time.Now()),
				EventSchedulesDatetime:    primitive.NewDateTimeFromTime(scheduleTime),
				EventSchedulesDescription: vSchedule.EventSchedulesDescription,
			})
		}
	}

	if len(errList) > 0 {
		helpers.ResponseBadRequestError(c, strings.Join(errList, ", ")+" are not in the event timeslot")
	} else {
		config.DB.Collection("EventSchedules").InsertMany(context.TODO(), eventSchedule)
		r.Retrieve(c)
	}
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
