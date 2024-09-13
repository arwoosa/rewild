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

type AchievementRepository struct{}

func (t AchievementRepository) Retrieve(c *gin.Context) {

	var results []models.EventsCountryCount
	userDetail := helpers.GetAuthUser(c)
	currentTime := primitive.NewDateTimeFromTime(time.Now())
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
		"events_date": bson.M{"$lt": currentTime},
	}

	filterStage := bson.D{{Key: "$match", Value: match}}
	groupStage := bson.D{
		{Key: "$group", Value: bson.M{
			"_id":                  "$events_country_code",
			"events_country_count": bson.M{"$sum": 1},
		}},
	}

	pipeline := mongo.Pipeline{
		lookupStage,
		unwindStage,
		filterStage,
		groupStage,
	}

	cursor, err := config.DB.Collection("Events").Aggregate(context.TODO(), pipeline)
	cursor.All(context.TODO(), &results)

	if err != nil {
		helpers.ResponseError(c, err.Error())
		return
	}

	if len(results) == 0 {
		helpers.ResponseNoData(c, "No Data")
		return
	}

	c.JSON(http.StatusOK, results)
}
