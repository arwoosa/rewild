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

	c.JSON(http.StatusOK, results)
}

func (r UserEventRepository) GetEventByUserId(c *gin.Context, userId primitive.ObjectID, Events *[]models.Events) error {
	isPast := c.Query("past")
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
		"EventParticipants.event_participants_user": userId,
		"events_deleted": bson.M{"$exists": false},
	}

	if isPast != "" {
		match["events_date"] = bson.M{"$lt": currentTime}
	} else {
		match["events_date"] = bson.M{"$gte": currentTime}
	}

	if countryCode != "" {
		match["events_country_code"] = countryCode
	}

	filterStage := bson.D{{Key: "$match", Value: match}}

	paginate := helpers.DataPaginate(c, 5)
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
			Key: "$facet", Value: bson.D{
				{Key: "data", Value: bson.A{
					paginate[0],
					paginate[1],
				}},
				{Key: "pagination", Value: bson.A{
					bson.D{{Key: "$count", Value: "total"}},
				}},
			},
		}},
	}

	var EventsPaginated []EventsPaginated

	cursor, err := config.DB.Collection("Events").Aggregate(context.TODO(), pipeline)
	cursor.All(context.TODO(), &EventsPaginated)

	*Events = EventsPaginated[0].Data
	if len(*Events) > 0 {
		*Events = EventRepository{}.RetrieveParticipantDetails(*Events)
	}

	return err
}

type Pagination struct {
	Total int64 `bson:"total" json:"total"`
}

type EventsPaginated struct {
	Data       []models.Events `bson:"data" json:"data"`
	Pagination []Pagination    `bson:"pagination" json:"pagination"`
}
