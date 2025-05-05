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

type UserEventRepository struct{}

func (r UserEventRepository) Retrieve(c *gin.Context) {
	var results []models.Events
	userDetail := helpers.GetAuthUser(c)

	err := r.GetEventByUserId(c, userDetail.UsersId, &results)

	if err != nil {
		return
	}

	if len(results) == 0 {
		helpers.ResponseNoData(c, "No Data")
		return
	}

	type returnData struct {
		No            string `json:"events_invitation_no"`
		models.Events `json:",inline"`
	}

	var responseData []returnData

	for _, event := range results {
		data := returnData{
			Events: event,
		}

		if !event.EventsDeadline.Time().IsZero() {
			data.No = event.EventsDeadline.Time().Format("0102")
		} else {
			data.No = defaultInvitationNo
		}

		responseData = append(responseData, data)
	}

	if len(responseData) == 0 {
		helpers.ResponseNoData(c, "Current login User Don't Have any Event or participant any event")
		return
	}

	c.JSON(http.StatusOK, responseData)
}

func (r UserEventRepository) GetEventByUserId(c *gin.Context, userId primitive.ObjectID, Events *[]models.Events) error {
	isPast := c.Query("past")
	hasPolaroid := c.Query("has_polaroid")
	countryCode := c.Query("country_code")
	currentTime := primitive.NewDateTimeFromTime(time.Now())

	// cursor, err := config.DB.Collection("Events").Find(context.TODO(), filter)
	lookupStage := bson.D{{Key: "$lookup", Value: bson.M{
		"from":         "EventParticipants",
		"localField":   "_id",
		"foreignField": "event_participants_event",
		"as":           "EventParticipants",
	}}}
	unwindStage := bson.D{
		{Key: "$unwind", Value: "$EventParticipants"},
	}

	match := bson.M{
		"EventParticipants.event_participants_user":   userId,
		"EventParticipants.event_participants_status": 1,
		"events_deleted": bson.M{"$exists": false},
	}

	if hasPolaroid == "true" {
		match["EventParticipants.event_participants_polaroid_count"] = bson.M{"$gt": 0}
	} else if hasPolaroid == "false" {
		match["$or"] = []bson.M{
			{"EventParticipants.event_participants_polaroid_count": bson.M{"$exists": false}},
			{"EventParticipants.event_participants_polaroid_count": bson.M{"$eq": 0}},
		}
	}

	if isPast != "" && isPast == "true" {
		match["events_date"] = bson.M{"$lt": currentTime}
	} else {
		match["events_date"] = bson.M{"$gte": currentTime}
	}

	if countryCode != "" {
		match["events_country_code"] = countryCode
	}

	filterStage := bson.D{{Key: "$match", Value: match}}

	paginate := helpers.DataPaginate(c, 5)
	dataFacet := bson.A{}

	if c.Query("page") != "" {
		dataFacet = bson.A{
			paginate[0],
			paginate[1],
		}
	}
	pipeline := mongo.Pipeline{
		lookupStage,
		unwindStage,
		filterStage,
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
		bson.D{{
			Key: "$facet", Value: bson.D{
				{Key: "data", Value: dataFacet},
				{Key: "pagination", Value: bson.A{
					bson.D{{Key: "$count", Value: "total"}},
				}},
			},
		}},
	}

	var EventsPaginated []EventsPaginated

	cursor, err := config.DB.Collection("Events").Aggregate(context.TODO(), pipeline)
	if err != nil {
		helpers.ResponseError(c, "資料庫查詢錯誤: "+err.Error())
		return err
	}

	if err := cursor.All(context.TODO(), &EventsPaginated); err != nil {
		helpers.ResponseError(c, "處理查詢結果錯誤: "+err.Error())
		return err
	}

	if len(EventsPaginated) == 0 {
		helpers.ResponseNoData(c, "查無符合條件的資料")
		return nil
	}

	*Events = EventsPaginated[0].Data
	if len(*Events) > 0 {
		*Events = EventRepository{}.RetrieveParticipantDetails(*Events)
	}

	return nil
}

type Pagination struct {
	Total int64 `bson:"total" json:"total"`
}

type EventsPaginated struct {
	Data       []models.Events `bson:"data" json:"data"`
	Pagination []Pagination    `bson:"pagination" json:"pagination"`
}
