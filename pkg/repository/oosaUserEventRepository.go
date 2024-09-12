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

type OosaUserEventRepository struct{}

func (r OosaUserEventRepository) Retrieve(c *gin.Context) {
	isPast := c.Query("past")
	var results []models.Events
	userDetail := helpers.GetAuthUser(c)
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
		"EventParticipants.event_participants_user": userDetail.UsersId,
	}

	if isPast != "" {
		match["events_date"] = bson.M{"$lt": currentTime}
	} else {
		match["events_date"] = bson.M{"$gte": currentTime}
	}

	filterStage := bson.D{{Key: "$match", Value: match}}

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
		bson.D{{Key: "$limit", Value: 5}},
	}

	cursor, err := config.DB.Collection("Events").Aggregate(context.TODO(), pipeline)
	cursor.All(context.TODO(), &results)

	if err != nil {
		return
	}

	if len(results) == 0 {
		helpers.ResponseNoData(c, "No Data")
		return
	}

	results = EventRepository{}.RetrieveParticipantDetails(results)
	c.JSON(http.StatusOK, results)
}
