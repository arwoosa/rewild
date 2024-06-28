package repository

import (
	"context"
	"fmt"
	"net/http"
	"oosa_rewild/internal/config"
	"oosa_rewild/internal/helpers"
	"oosa_rewild/internal/models"
	"oosa_rewild/pkg/service"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type EventRepository struct{}
type EventRequest struct {
	EventsDate      string `json:"events_date" validate:"required,datetime=2006-01-02 15:04:05"`
	EventsDateEnd   string `json:"events_date_end" validate:"required,datetime=2006-01-02 15:04:05"`
	EventsDeadline  string `json:"events_deadline" validate:"required,datetime=2006-01-02 15:04:05"`
	EventsName      string `json:"events_name" validate:"required"`
	EventsPlace     string `json:"events_place" validate:"required_without=EventsRewilding"`
	EventsRewilding string `json:"events_rewilding" validate:"required_without=EventsPlace"`
	// EventsType             string  `json:"events_type" validate:"required"`
	EventsParticipantLimit int     `json:"events_participant_limit" validate:"required"`
	EventsPaymentRequired  int     `json:"events_payment_required" validate:"required"`
	EventsPaymentFee       float64 `json:"events_payment_fee" validate:"required"`
	EventsRequiresApproval int     `json:"events_requires_approval" validate:"required"`
	EventsLat              float64 `json:"events_lat" validate:"required_without_all=EventsPlace EventsRewilding"`
	EventsLng              float64 `json:"events_lng" validate:"required_without_all=EventsPlace EventsRewilding"`
	EventsMeetingPointLat  float64 `json:"events_meeting_point_lat" validate:"required"`
	EventsMeetingPointLng  float64 `json:"events_meeting_point_lng" validate:"required"`
}

func (r EventRepository) Retrieve(c *gin.Context) {
	var results []models.Events
	filter := bson.M{
		"events_date": bson.M{"$gte": primitive.NewDateTimeFromTime(time.Now())},
	}
	cursor, err := config.DB.Collection("Events").Find(context.TODO(), filter)
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

func (r EventRepository) Create(c *gin.Context) {
	userDetail := helpers.GetAuthUser(c)
	var payload EventRequest
	validateError := helpers.Validate(c, &payload)
	if validateError != nil {
		return
	}

	lat := helpers.FloatToDecimal128(payload.EventsLat)
	lng := helpers.FloatToDecimal128(payload.EventsLng)

	meetingLat := helpers.FloatToDecimal128(payload.EventsMeetingPointLat)
	meetingLng := helpers.FloatToDecimal128(payload.EventsMeetingPointLng)

	if payload.EventsPlace != "" {
		rewildingId := service.GoogleToRewilding(c, payload.EventsPlace)
		if helpers.MongoZeroID(rewildingId) {
			return
		}
		payload.EventsRewilding = rewildingId.Hex()
	}

	if payload.EventsRewilding != "" {
		rewildingId := helpers.StringToPrimitiveObjId(payload.EventsRewilding)
		rewildingData, err := service.GetRewildingById(rewildingId)
		if err != nil {
			helpers.ResultEmpty(c, err)
			return
		}
		lat = rewildingData.RewildingLat
		lng = rewildingData.RewildingLng
		if rewildingData.RewildingPlaceId != "" {
			payload.EventsPlace = rewildingData.RewildingPlaceId
		}
	}

	insert := models.Events{
		EventsDate:      helpers.StringToPrimitiveDateTime(payload.EventsDate),
		EventsDateEnd:   helpers.StringToPrimitiveDateTime(payload.EventsDateEnd),
		EventsDeadline:  helpers.StringToPrimitiveDateTime(payload.EventsDeadline),
		EventsName:      payload.EventsName,
		EventsRewilding: helpers.StringToPrimitiveObjId(payload.EventsRewilding),
		EventsPlace:     payload.EventsPlace,
		//EventsType:             "",
		EventsParticipantLimit: payload.EventsParticipantLimit,
		EventsPaymentRequired:  payload.EventsPaymentRequired,
		EventsPaymentFee:       payload.EventsPaymentFee,
		EventsRequiresApproval: payload.EventsRequiresApproval,
		EventsLat:              lat,
		EventsLng:              lng,
		EventsMeetingPointLat:  meetingLat,
		EventsMeetingPointLng:  meetingLng,
		EventsCreatedBy:        userDetail.UsersId,
		EventsCreatedAt:        primitive.NewDateTimeFromTime(time.Now()),
	}

	result, err := config.DB.Collection("Events").InsertOne(context.TODO(), insert)
	if err != nil {
		fmt.Println("ERROR", err.Error())
		return
	}

	var Events models.Events
	config.DB.Collection("Events").FindOne(context.TODO(), bson.D{{Key: "_id", Value: result.InsertedID}}).Decode(&Events)
	c.JSON(http.StatusOK, Events)
}

func (r EventRepository) Read(c *gin.Context) {
	var Events models.Events
	err := r.ReadOne(c, &Events)
	if err == nil {
		c.JSON(http.StatusOK, Events)
	}
}

func (r EventRepository) ReadOne(c *gin.Context, Events *models.Events) error {
	id, _ := primitive.ObjectIDFromHex(c.Param("id"))
	filter := bson.D{{Key: "_id", Value: id}}
	err := config.DB.Collection("Events").FindOne(context.TODO(), filter).Decode(&Events)
	if err != nil {
		helpers.ResultEmpty(c, err)
	}
	return err
}

func (r EventRepository) Update(c *gin.Context) {
	userDetail := helpers.GetAuthUser(c)
	var Events models.Events
	var payload EventRequest
	validateError := helpers.Validate(c, &payload)
	if validateError != nil {
		return
	}

	err := r.ReadOne(c, &Events)
	if err == nil {
		Events.EventsUpdatedBy = userDetail.UsersId
		Events.EventsUpdatedAt = primitive.NewDateTimeFromTime(time.Now())

		r.ProcessData(c, &Events, payload)
		filters := bson.D{{Key: "_id", Value: Events.EventsId}}
		upd := bson.D{{Key: "$set", Value: Events}}
		config.DB.Collection("Events").UpdateOne(context.TODO(), filters, upd)
		c.JSON(http.StatusOK, Events)
	}
}

func (r EventRepository) ProcessData(c *gin.Context, Events *models.Events, payload EventRequest) {
	lat := helpers.FloatToDecimal128(payload.EventsLat)
	lng := helpers.FloatToDecimal128(payload.EventsLng)

	Events.EventsDate = helpers.StringToPrimitiveDateTime(payload.EventsDate)
	Events.EventsDeadline = helpers.StringToPrimitiveDateTime(payload.EventsDeadline)
	Events.EventsName = payload.EventsName
	Events.EventsRewilding = helpers.StringToPrimitiveObjId(payload.EventsRewilding)
	Events.EventsPlace = payload.EventsPlace
	Events.EventsPaymentRequired = payload.EventsPaymentRequired
	Events.EventsPaymentFee = payload.EventsPaymentFee
	Events.EventsRequiresApproval = payload.EventsRequiresApproval
	Events.EventsLat = lat
	Events.EventsLng = lng
}
