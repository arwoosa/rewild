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
	EventsDate             string   `json:"events_date" form:"events_date" validate:"required,datetime=2006-01-02T15:04:05Z07:00"`
	EventsDateEnd          string   `json:"events_date_end" form:"events_date_end" validate:"required,datetime=2006-01-02T15:04:05Z07:00"`
	EventsDeadline         string   `json:"events_deadline" form:"events_deadline" validate:"required,datetime=2006-01-02T15:04:05Z07:00"`
	EventsName             string   `json:"events_name" form:"events_name" validate:"required"`
	EventsPlace            string   `json:"events_place" form:"events_place" validate:"required_without=EventsRewilding"`
	EventsPhotoCover       string   `json:"events_photo_cover" form:"events_photo_cover"`
	EventsRewilding        string   `json:"events_rewilding" form:"events_rewilding" validate:"required_without=EventsPlace"`
	EventsType             string   `json:"events_type" form:"events_type" validate:"required"`
	EventsParticipantLimit int      `json:"events_participant_limit" form:"events_participant_limit"`
	EventsPaymentRequired  int      `json:"events_payment_required" form:"events_payment_required" default:"0"`
	EventsPaymentFee       float64  `json:"events_payment_fee" form:"events_payment_fee"`
	EventsRequiresApproval int      `json:"events_requires_approval" form:"events_requires_approval" default:"0"`
	EventsLat              float64  `json:"events_lat" form:"events_lat" validate:"required_without_all=EventsPlace EventsRewilding"`
	EventsLng              float64  `json:"events_lng" form:"events_lng" validate:"required_without_all=EventsPlace EventsRewilding"`
	EventsCountryCode      string   `json:"events_country_code" form:"events_country_code"`
	EventsMeetingPointLat  float64  `json:"events_meeting_point_lat" form:"events_meeting_point_lat" validate:"required"`
	EventsMeetingPointLng  float64  `json:"events_meeting_point_lng" form:"events_meeting_point_lng" validate:"required"`
	EventsMeetingPointName string   `json:"events_meeting_point_name" form:"events_meeting_point_name" validate:"required"`
	EventsParticipants     []string `json:"events_participants" form:"events_participants"`
}

func (r EventRepository) Retrieve(c *gin.Context) {
	var results []models.Events

	eventPeriodBegin := c.Query("event_period_begin")
	eventPeriodEnd := c.Query("event_period_end")
	eventRewilding := c.Query("event_rewilding")
	eventIsPast := c.Query("event_past")

	var eventBegin time.Time
	var eventEnd time.Time

	match := bson.M{
		"events_deleted": bson.M{"$exists": false},
	}

	currentTime := primitive.NewDateTimeFromTime(time.Now())

	if eventPeriodBegin != "" && eventPeriodEnd == "" {
		eventBegin = helpers.StringToDatetime(eventPeriodBegin + " 00:00:00")
		eventEnd = helpers.StringToDatetime(eventPeriodBegin + " 23:59:59")
		match["events_date"] = bson.M{
			"$gte": primitive.NewDateTimeFromTime(eventBegin),
			"$lte": primitive.NewDateTimeFromTime(eventEnd),
		}

	} else if eventPeriodBegin == "" && eventPeriodEnd != "" {
		eventBegin = helpers.StringToDatetime(eventPeriodEnd + " 00:00:00")
		eventEnd = helpers.StringToDatetime(eventPeriodEnd + " 23:59:59")
		match["events_date"] = bson.M{
			"$gte": primitive.NewDateTimeFromTime(eventBegin),
			"$lte": primitive.NewDateTimeFromTime(eventEnd),
		}

	} else if eventPeriodBegin != "" && eventPeriodEnd != "" {
		eventBegin = helpers.StringToDatetime(eventPeriodBegin + " 00:00:00")
		eventEnd = helpers.StringToDatetime(eventPeriodEnd + " 23:59:59")
		match["events_date"] = bson.M{
			"$gte": primitive.NewDateTimeFromTime(eventBegin),
			"$lte": primitive.NewDateTimeFromTime(eventEnd),
		}
	} else if eventIsPast == "1" {
		match["events_date"] = bson.M{"$lt": currentTime}
	} else {
		match["events_date"] = bson.M{"$gte": currentTime}
	}

	if eventRewilding != "" {
		rewildingId := helpers.StringToPrimitiveObjId(eventRewilding)
		match["events_rewilding"] = rewildingId
	}

	filterPeriod := bson.D{{
		Key: "$match", Value: match,
	}}

	agg := mongo.Pipeline{
		filterPeriod,
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
		bson.D{{
			Key: "$lookup", Value: bson.M{
				"from":         "Rewilding",
				"localField":   "events_rewilding",
				"foreignField": "_id",
				"as":           "events_rewilding_detail",
			},
		}},
		bson.D{{
			Key: "$unwind", Value: "$events_rewilding_detail",
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

	results = r.RetrieveParticipantDetails(results)
	c.JSON(http.StatusOK, results)
}

func (r EventRepository) RetrieveParticipantDetails(results []models.Events) []models.Events {
	for k, v := range results {
		var EventParticipantsUser []models.UsersAgg
		var EventParticipants []models.EventParticipants

		agg := mongo.Pipeline{
			bson.D{{
				Key: "$match", Value: bson.M{
					"event_participants_event":  v.EventsId,
					"event_participants_status": 1,
				},
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
		cursor, _ := config.DB.Collection("EventParticipants").Aggregate(context.TODO(), agg)
		cursor.All(context.TODO(), &EventParticipants)

		maxSlice := 3
		noOfParticipants := len(EventParticipants)
		remaining := 0
		if noOfParticipants == 0 {
			EventParticipantsUser = make([]models.UsersAgg, 0)
		} else {
			for kU, vU := range EventParticipants {
				if kU < maxSlice {
					EventParticipantsUser = append(EventParticipantsUser, *vU.EventParticipantsUserDetail)
				}
			}
		}

		if noOfParticipants > maxSlice {
			remaining = noOfParticipants - maxSlice
		}

		EventParticipantsList := models.EventParticipantObj{
			LatestTreeUser: &EventParticipantsUser,
			RemainNumber:   remaining,
		}

		results[k].EventsParticipants = &EventParticipantsList
	}
	return results
}

func (r EventRepository) Create(c *gin.Context) {
	userDetail := helpers.GetAuthUser(c)
	var payload EventRequest
	validateError := helpers.ValidateForm(c, &payload)
	if validateError != nil {
		return
	}

	match, errMessage := helpers.ValidateStringStyle1(payload.EventsName, int(config.APP_LIMIT.LengthEventName))
	if !match {
		helpers.ResponseBadRequestError(c, "Name can only contain "+errMessage)
		return
	}

	eventStartDate := helpers.StringToDateTime(payload.EventsDate)
	isPastEvent := eventStartDate.Before(time.Now())

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
		payload.EventsLat = lat
		payload.EventsLng = lng
		payload.EventsCountryCode = rewildingData.RewildingCountryCode
	}

	openWeather := openweather.OpenWeatherRepository{}.Forecast(c, payload.EventsLat, payload.EventsLng)
	insert := models.Events{
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

		if payload.EventsPhotoCover == "" {
			helpers.ResponseBadRequestError(c, "Please select a cover photo")
			return
		} else {
			if payload.EventsPhotoCover != "DEFAULT_0" && payload.EventsPhotoCover != "DEFAULT_1" {
				helpers.ResponseBadRequestError(c, "Invalid cover photo. Use either DEFAULT_0 or DEFAULT_1")
				return
			}
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

	r.ProcessData(c, &insert, payload)

	result, err := config.DB.Collection("Events").InsertOne(context.TODO(), insert)
	if err != nil {
		fmt.Println("ERROR", err.Error())
		return
	}

	var Events models.Events
	config.DB.Collection("Events").FindOne(context.TODO(), bson.D{{Key: "_id", Value: result.InsertedID}}).Decode(&Events)

	// Create participant records
	insertParticipant := []interface{}{
		models.EventParticipants{
			EventParticipantsEvent:     Events.EventsId,
			EventParticipantsUser:      userDetail.UsersId,
			EventParticipantsStatus:    GetEventParticipantStatus("ACCEPTED"),
			EventParticipantsCreatedBy: userDetail.UsersId,
			EventParticipantsCreatedAt: primitive.NewDateTimeFromTime(time.Now()),
		},
	}

	if isPastEvent && len(payload.EventsParticipants) > 0 {
		// Handle participant
		for _, v := range payload.EventsParticipants {
			insertParticipant = append(insertParticipant, models.EventParticipants{
				EventParticipantsEvent:     Events.EventsId,
				EventParticipantsUser:      helpers.StringToPrimitiveObjId(v),
				EventParticipantsStatus:    GetEventParticipantStatus("ACCEPTED"),
				EventParticipantsCreatedBy: userDetail.UsersId,
				EventParticipantsCreatedAt: primitive.NewDateTimeFromTime(time.Now()),
			})
		}
	}
	config.DB.Collection("EventParticipants").InsertMany(context.TODO(), insertParticipant)

	// Create badge record
	helpers.BadgeEvents(c, Events.EventsId)
	r.HandleParticipation(c, userDetail.UsersId, Events.EventsId)
	c.JSON(http.StatusOK, Events)
}

func (r EventRepository) HandleParticipation(c *gin.Context, userId primitive.ObjectID, eventId primitive.ObjectID) {
	filter := bson.D{{Key: "event_participants_user", Value: userId}}
	count, _ := config.DB.Collection("EventParticipants").CountDocuments(context.TODO(), filter)

	UpdateUser := models.Users{
		UsersEventScheduled: int(count),
	}
	filters := bson.D{{Key: "_id", Value: userId}}
	upd := bson.D{{Key: "$set", Value: UpdateUser}}
	config.DB.Collection("Users").UpdateOne(context.TODO(), filters, upd)

	expAvailable := map[int]int{1: 5, 2: 4, 3: 3, 4: 2, 5: 1}
	expAwarded := expAvailable[int(count)]
	if expAwarded > 0 {
		helpers.ExpAward(c, helpers.EXP_REWILDING, expAwarded, eventId)
	}
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
	validateError := helpers.ValidateForm(c, &payload)
	if validateError != nil {
		return
	}

	err := r.ReadOne(c, &Events)
	if err == nil {
		match, errMessage := helpers.ValidateStringStyle1(payload.EventsName, int(config.APP_LIMIT.LengthEventName))
		if !match {
			helpers.ResponseBadRequestError(c, "Name can only contain "+errMessage)
			return
		}

		Events.EventsUpdatedBy = userDetail.UsersId
		Events.EventsUpdatedAt = primitive.NewDateTimeFromTime(time.Now())

		file, err := helpers.ValidatePhotoRequest(c, "events_photo", false)

		if file == nil {
			// HAS FILE
			if err != nil {
				return
			}

			if payload.EventsPhotoCover == "" {
				helpers.ResponseBadRequestError(c, "Please select a cover photo")
				return
			} else {
				if payload.EventsPhotoCover != "DEFAULT_0" && payload.EventsPhotoCover != "DEFAULT_1" {
					helpers.ResponseBadRequestError(c, "Invalid cover photo. Use either DEFAULT_0 or DEFAULT_1")
					return
				}
			}
			Events.EventsPhoto = ""
		} else {
			cloudflare := CloudflareRepository{}
			cloudflareResponse, postErr := cloudflare.Post(c, file)
			if postErr != nil {
				helpers.ResponseBadRequestError(c, postErr.Error())
				return
			}
			fileName := cloudflare.ImageDelivery(cloudflareResponse.Result.Id, "public")
			Events.EventsPhoto = fileName
		}
		r.ProcessData(c, &Events, payload)
		filters := bson.D{{Key: "_id", Value: Events.EventsId}}
		upd := bson.D{{Key: "$set", Value: Events}}
		config.DB.Collection("Events").UpdateOne(context.TODO(), filters, upd)

		ActiveParticipants := EventParticipantsRepository{}.ActiveParticipants(Events.EventsId)
		for _, v := range ActiveParticipants {
			NotificationMessage := models.NotificationMessage{
				Message: "團主於{0}中更新了重要資訊! 點擊查看",
				Data:    []map[string]interface{}{helpers.NotificationFormatEvent(Events)},
			}
			helpers.NotificationsCreate(c, helpers.NOTIFICATION_EVENT_INFO, v.EventParticipantsUser, NotificationMessage, Events.EventsId)
		}

		r.Read(c)
	}
}

func (r EventRepository) ProcessData(c *gin.Context, Events *models.Events, payload EventRequest) {
	var Rewilding models.Rewilding
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
	Events.EventsType = helpers.StringToPrimitiveObjId(payload.EventsType)
	Events.EventsPaymentRequired = payload.EventsPaymentRequired
	Events.EventsPaymentFee = &payload.EventsPaymentFee
	Events.EventsRequiresApproval = &payload.EventsRequiresApproval
	Events.EventsMeetingPointLat = meetingPointLat
	Events.EventsMeetingPointLng = meetingPointLng
	Events.EventsMeetingPointName = payload.EventsMeetingPointName
	Events.EventsLat = eventsLat
	Events.EventsLng = eventsLng
	Events.EventsParticipantLimit = &payload.EventsParticipantLimit
	Events.EventsCountryCode = payload.EventsCountryCode

	config.DB.Collection("Rewilding").FindOne(context.TODO(), bson.D{{Key: "_id", Value: Events.EventsRewilding}}).Decode(&Rewilding)

	Events.EventsRewildingAchievementType = Rewilding.RewildingAchievementType
	Events.EventsRewildingAchievementTypeID = Rewilding.RewildingAchievementTypeID

	if Events.EventsPhoto == "" && payload.EventsPhotoCover != "" {
		coverImage := config.APP.BaseUrl + "event/cover/"
		if payload.EventsPhotoCover == "DEFAULT_1" {
			var RefRewildingTypes models.RefRewildingTypes
			config.DB.Collection("RefRewildingTypes").FindOne(context.TODO(), bson.D{{Key: "_id", Value: helpers.StringToPrimitiveObjId(payload.EventsType)}}).Decode(&RefRewildingTypes)
			coverImage += RefRewildingTypes.RefRewildingTypesDefaultImage
		} else {
			coverImage += "default_oosa_event.png"
		}
		Events.EventsPhoto = coverImage
	}

	Events.EventsStatisticDistance = eventStatisticDistance
	Events.EventsStatisticTime = eventDurationSecond
}

func (r EventRepository) Options(c *gin.Context) {
	RefRewildingTypes := helpers.RefRewildingTypes()
	c.JSON(http.StatusOK, gin.H{"rewilding_types": RefRewildingTypes})
}
