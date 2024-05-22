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

	cursor, err := config.DB.Collection("Events").Aggregate(context.TODO(), mongo.Pipeline{
		lookupStage,
		unwindStage,
		filterStage,
	})
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

/*

db.Events.aggregate(
    [
		{
			$lookup: {
				from: "EventParticipants",
				localField: "_id",
				foreignField: "event_participants_event",
				as: "EventParticipants"
			}
		},
		{
			$unwind: "$EventParticipants"
		},
		{
			$match : {
				"events_date": { $gte: ISODate("2024-06-01T00:00:00.000Z") },
				"EventParticipants.event_participants_user" : ObjectId("6613fae21977170d4fd608a6")
			}
		}
	]
);

*/
