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
	var results []models.AchievementRewilding
	userDetail := helpers.GetAuthUser(c)

	err := t.GetAchievementsByUserId(c, userDetail.UsersId, &results)
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

func (t AchievementRepository) GetAchievementsByUserId(c *gin.Context, userId primitive.ObjectID, results *[]models.AchievementRewilding) error {
	achievementType := c.Query("achievement_type")

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
		"EventParticipants.event_participants_user": userId,
		"events_date": bson.M{"$lt": currentTime},
	}

	if achievementType != "" {
		match["events_rewilding_achievement_type"] = achievementType
		match["events_rewilding_achievement_eligible"] = true
	}

	filterStage := bson.D{{Key: "$match", Value: match}}
	groupStage := bson.D{
		{Key: "$group", Value: bson.M{
			"_id":             "$events_rewilding",
			"rewilding_count": bson.M{"$sum": 1},
		}},
	}

	pipeline := mongo.Pipeline{
		lookupStage,
		unwindStage,
	}

	rewildLookup := bson.D{{Key: "$lookup", Value: bson.M{
		"from":         "Rewilding",
		"localField":   "_id",
		"foreignField": "_id",
		"as":           "RewildingDetail",
	}}}
	rewildUnwind := bson.D{
		{Key: "$unwind", Value: "$RewildingDetail"},
	}
	replaceRoot := bson.D{{
		Key: "$replaceRoot", Value: bson.M{
			"newRoot": bson.M{
				"$mergeObjects": bson.A{
					"$RewildingDetail",
					bson.M{"rewilding_count": "$rewilding_count"},
				},
			},
		},
	}}

	pipeline = append(pipeline, filterStage, groupStage, rewildLookup, rewildUnwind, replaceRoot)
	cursor, err := config.DB.Collection("Events").Aggregate(context.TODO(), pipeline)
	cursor.All(context.TODO(), results)

	return err
}
