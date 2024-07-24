package repository

import (
	"context"
	"fmt"
	"net/http"
	"oosa_rewild/internal/config"
	"oosa_rewild/internal/helpers"
	"oosa_rewild/internal/models"
	"oosa_rewild/pkg/openweather"
	"oosa_rewild/pkg/service"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type EventRepository struct{}
type EventRequest struct {
	EventsDate      string `json:"events_date" validate:"required,datetime=2006-01-02T15:04:05Z07:00"`
	EventsDateEnd   string `json:"events_date_end" validate:"required,datetime=2006-01-02T15:04:05Z07:00"`
	EventsDeadline  string `json:"events_deadline" validate:"required,datetime=2006-01-02T15:04:05Z07:00"`
	EventsName      string `json:"events_name" validate:"required"`
	EventsPlace     string `json:"events_place" validate:"required_without=EventsRewilding"`
	EventsRewilding string `json:"events_rewilding" validate:"required_without=EventsPlace"`
	// EventsType             string  `json:"events_type" validate:"required"`
	EventsParticipantLimit int     `json:"events_participant_limit" validate:"required"`
	EventsPaymentRequired  int     `json:"events_payment_required" default:"0"`
	EventsPaymentFee       float64 `json:"events_payment_fee" validate:"required"`
	EventsRequiresApproval int     `json:"events_requires_approval" default:"0"`
	EventsLat              float64 `json:"events_lat" validate:"required_without_all=EventsPlace EventsRewilding"`
	EventsLng              float64 `json:"events_lng" validate:"required_without_all=EventsPlace EventsRewilding"`
	EventsMeetingPointLat  float64 `json:"events_meeting_point_lat" validate:"required"`
	EventsMeetingPointLng  float64 `json:"events_meeting_point_lng" validate:"required"`
}

type EventFormDataRequest struct {
	EventsDate      string `form:"events_date" validate:"required,datetime=2006-01-02T15:04:05Z07:00"`
	EventsDateEnd   string `form:"events_date_end" validate:"required,datetime=2006-01-02T15:04:05Z07:00"`
	EventsDeadline  string `form:"events_deadline" validate:"required,datetime=2006-01-02T15:04:05Z07:00"`
	EventsName      string `form:"events_name" validate:"required"`
	EventsPlace     string `form:"events_place" validate:"required_without=EventsRewilding"`
	EventsRewilding string `form:"events_rewilding" validate:"required_without=EventsPlace"`
	// EventsType             string  `form:"events_type" validate:"required"`
	EventsParticipantLimit int     `form:"events_participant_limit" validate:"required"`
	EventsPaymentRequired  int     `form:"events_payment_required" default:"0"`
	EventsPaymentFee       float64 `form:"events_payment_fee" validate:"required"`
	EventsRequiresApproval int     `form:"events_requires_approval" default:"0"`
	EventsLat              float64 `form:"events_lat" validate:"required_without_all=EventsPlace EventsRewilding"`
	EventsLng              float64 `form:"events_lng" validate:"required_without_all=EventsPlace EventsRewilding"`
	EventsMeetingPointLat  float64 `form:"events_meeting_point_lat" validate:"required"`
	EventsMeetingPointLng  float64 `form:"events_meeting_point_lng" validate:"required"`
}

func (r EventRepository) Retrieve(c *gin.Context) {
	var results []models.Events
	agg := mongo.Pipeline{
		bson.D{{
			Key: "$match", Value: bson.M{
				"events_date":    bson.M{"$gte": primitive.NewDateTimeFromTime(time.Now())},
				"events_deleted": bson.M{"$exists": false},
			},
		}},
		bson.D{{
			Key: "$lookup", Value: bson.M{
				"from":         "Users",
				"localField":   "events_created_by",
				"foreignField": "_id",
				"as":           "events_created_by_user",
			},
		}},
		bson.D{{
			Key: "$unwind", Value: "$events_created_by_user",
		}},
	}

	cursor, err := config.DB.Collection("Events").Aggregate(context.TODO(), agg)
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
	var payload EventFormDataRequest
	validateError := helpers.ValidateForm(c, &payload)
	if validateError != nil {
		return
	}

	// lat := helpers.FloatToDecimal128(payload.EventsLat)
	// lng := helpers.FloatToDecimal128(payload.EventsLng)

	// meetingLat := helpers.FloatToDecimal128(payload.EventsMeetingPointLat)
	// meetingLng := helpers.FloatToDecimal128(payload.EventsMeetingPointLng)

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
		lat := rewildingData.RewildingLat
		lng := rewildingData.RewildingLng
		if rewildingData.RewildingPlaceId != "" {
			payload.EventsPlace = rewildingData.RewildingPlaceId
		}
		payload.EventsLat = helpers.Decimal128ToFloat(lat)
		payload.EventsLng = helpers.Decimal128ToFloat(lng)
	}

	openWeather := openweather.OpenWeatherRepository{}.Forecast(c, payload.EventsLat, payload.EventsLng)
	insert := models.Events{
		// EventsDate:      helpers.StringToPrimitiveDateTime(payload.EventsDate),
		// EventsDateEnd:   helpers.StringToPrimitiveDateTime(payload.EventsDateEnd),
		// EventsDeadline:  helpers.StringToPrimitiveDateTime(payload.EventsDeadline),
		// EventsName:      payload.EventsName,
		// EventsRewilding: helpers.StringToPrimitiveObjId(payload.EventsRewilding),
		// EventsPlace: payload.EventsPlace,
		// EventsType:             "",
		// EventsParticipantLimit: payload.EventsParticipantLimit,
		// EventsPaymentRequired:  payload.EventsPaymentRequired,
		// EventsPaymentFee:       payload.EventsPaymentFee,
		// EventsRequiresApproval: &payload.EventsRequiresApproval,
		// EventsLat: lat,
		// EventsLng: lng,
		// EventsMeetingPointLat: meetingLat,
		// EventsMeetingPointLng: meetingLng,
		EventsCityId:    openWeather.City.Id,
		EventsCreatedBy: userDetail.UsersId,
		EventsCreatedAt: primitive.NewDateTimeFromTime(time.Now()),
	}

	file, err := helpers.ValidatePhotoRequest(c, "events_photo", false)

	if file == nil {
		// HAS FILE
		if err != nil {
			return
		}
	} else {
		cloudflare := CloudflareRepository{}
		cloudflareResponse, postErr := cloudflare.Post(c, file)
		if postErr != nil {
			helpers.ResponseBadRequestError(c, postErr.Error())
			return
		}
		fileName := cloudflare.ImageDelivery(cloudflareResponse.Result.Id, "public")
		insert.EventsPhoto = fileName
	}
	r.ProcessDataForm(c, &insert, payload)

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
	var Events []models.Events
	id, _ := primitive.ObjectIDFromHex(c.Param("id"))
	agg := mongo.Pipeline{
		bson.D{{
			Key: "$match", Value: bson.M{"_id": id},
		}},
		bson.D{{
			Key: "$lookup", Value: bson.M{
				"from":         "Users",
				"localField":   "events_created_by",
				"foreignField": "_id",
				"as":           "events_created_by_user",
			},
		}},
		bson.D{{
			Key: "$unwind", Value: "$events_created_by_user",
		}},
		bson.D{{Key: "$limit", Value: 1}},
	}

	cursor, err := config.DB.Collection("Events").Aggregate(context.TODO(), agg)
	cursor.All(context.TODO(), &Events)
	// err := r.ReadOne(c, &Events)
	if err == nil {
		if len(Events) == 0 {
			helpers.ResponseNoData(c, "No data")
		} else {
			c.JSON(http.StatusOK, Events[0])
		}
	}
}

func (r EventRepository) Delete(c *gin.Context) {
	var Events models.Events
	err := r.ReadOne(c, &Events)

	// TODO: What are the rules of allowing delete

	isDeleted := 1
	Events.EventsDeleted = &isDeleted
	Events.EventsDeletedAt = primitive.NewDateTimeFromTime(time.Now())

	if err == nil {
		filters := bson.D{{Key: "_id", Value: Events.EventsId}}
		upd := bson.D{{Key: "$set", Value: Events}}
		config.DB.Collection("Events").UpdateOne(context.TODO(), filters, upd)
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
		r.Read(c)
	}
}

func (r EventRepository) ProcessData(c *gin.Context, Events *models.Events, payload EventRequest) {
	eventsLat := payload.EventsLat
	eventsLng := payload.EventsLng
	meetingPointLat := payload.EventsMeetingPointLat
	meetingPointLng := payload.EventsMeetingPointLng

	eventDate := helpers.StringToPrimitiveDateTime(payload.EventsDate)
	eventDateEnd := helpers.StringToPrimitiveDateTime(payload.EventsDateEnd)
	eventDurationSecond := eventDateEnd.Time().Sub(eventDate.Time()).Seconds()

	eventStatisticDistance := (helpers.Haversine(eventsLat, eventsLng, meetingPointLat, meetingPointLng) * 100000) / 70

	Events.EventsDate = eventDate
	Events.EventsDateEnd = eventDateEnd
	Events.EventsDeadline = helpers.StringToPrimitiveDateTime(payload.EventsDeadline)
	Events.EventsName = payload.EventsName
	Events.EventsRewilding = helpers.StringToPrimitiveObjId(payload.EventsRewilding)
	Events.EventsPlace = payload.EventsPlace
	Events.EventsPaymentRequired = payload.EventsPaymentRequired
	Events.EventsPaymentFee = payload.EventsPaymentFee
	Events.EventsRequiresApproval = &payload.EventsRequiresApproval
	Events.EventsMeetingPointLat = helpers.FloatToDecimal128(meetingPointLat)
	Events.EventsMeetingPointLng = helpers.FloatToDecimal128(meetingPointLng)
	Events.EventsLat = helpers.FloatToDecimal128(eventsLat)
	Events.EventsLng = helpers.FloatToDecimal128(eventsLng)
	Events.EventsParticipantLimit = payload.EventsParticipantLimit

	Events.EventsStatisticDistance = helpers.FloatToDecimal128(eventStatisticDistance)
	Events.EventsStatisticTime = eventDurationSecond
}

func (r EventRepository) ProcessDataForm(c *gin.Context, Events *models.Events, payload EventFormDataRequest) {
	eventsLat := payload.EventsLat
	eventsLng := payload.EventsLng
	meetingPointLat := payload.EventsMeetingPointLat
	meetingPointLng := payload.EventsMeetingPointLng

	eventDate := helpers.StringToPrimitiveDateTime(payload.EventsDate)
	eventDateEnd := helpers.StringToPrimitiveDateTime(payload.EventsDateEnd)
	eventDurationSecond := eventDateEnd.Time().Sub(eventDate.Time()).Seconds()

	eventStatisticDistance := (helpers.Haversine(eventsLat, eventsLng, meetingPointLat, meetingPointLng) * 100000) / 70

	Events.EventsDate = eventDate
	Events.EventsDateEnd = eventDateEnd
	Events.EventsDeadline = helpers.StringToPrimitiveDateTime(payload.EventsDeadline)
	Events.EventsName = payload.EventsName
	Events.EventsRewilding = helpers.StringToPrimitiveObjId(payload.EventsRewilding)
	Events.EventsPlace = payload.EventsPlace
	Events.EventsPaymentRequired = payload.EventsPaymentRequired
	Events.EventsPaymentFee = payload.EventsPaymentFee
	Events.EventsRequiresApproval = &payload.EventsRequiresApproval
	Events.EventsMeetingPointLat = helpers.FloatToDecimal128(meetingPointLat)
	Events.EventsMeetingPointLng = helpers.FloatToDecimal128(meetingPointLng)
	Events.EventsLat = helpers.FloatToDecimal128(eventsLat)
	Events.EventsLng = helpers.FloatToDecimal128(eventsLng)
	Events.EventsParticipantLimit = payload.EventsParticipantLimit

	Events.EventsStatisticDistance = helpers.FloatToDecimal128(eventStatisticDistance)
	Events.EventsStatisticTime = eventDurationSecond
}

func (r EventRepository) Options(c *gin.Context) {
	var RefRewildingTypes []models.RefRewildingTypes
	cursor, err := config.DB.Collection("RefRewildingTypes").Find(context.TODO(), bson.D{})
	if err != nil {
		return
	}
	cursor.All(context.TODO(), &RefRewildingTypes)
	c.JSON(http.StatusOK, gin.H{"rewilding_types": RefRewildingTypes})
}
