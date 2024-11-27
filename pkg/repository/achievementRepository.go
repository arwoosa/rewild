package repository

import (
	"context"
	"net/http"
	"oosa_rewild/internal/config"
	"oosa_rewild/internal/helpers"
	"oosa_rewild/internal/models"
	"strconv"
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
	starType := c.Query("star_type")
	country := c.Query("country")
	otherUserId := c.Query("user_id")
	otherUserObjId := helpers.StringToPrimitiveObjId(otherUserId)

	currentTime := primitive.NewDateTimeFromTime(time.Now())
	lookupStage := bson.D{{Key: "$lookup", Value: bson.M{
		"as":   "EventParticipants",
		"from": "EventParticipants",
		"let":  bson.M{"event_id": "$_id"},
		"pipeline": []bson.M{
			{"$match": bson.M{
				"$expr": bson.M{
					"$eq": bson.A{"$event_participants_event", "$$event_id"}},
			}},
			{"$match": bson.M{"event_participants_user": userId}},
		},
	}}}

	lookupOtherUser := bson.D{{
		Key: "$lookup", Value: bson.M{
			"as":   "OtherParticipant",
			"from": "EventParticipants",
			"let":  bson.M{"event_id": "$_id"},
			"pipeline": []bson.M{
				{"$match": bson.M{
					"$expr": bson.M{
						"$eq": bson.A{"$event_participants_event", "$$event_id"}},
				}},
				{"$match": bson.M{"event_participants_user": otherUserObjId}},
			},
		},
	}}

	unwindStage := bson.D{
		{Key: "$unwind", Value: "$EventParticipants"},
	}

	match := bson.M{
		"EventParticipants.event_participants_achievement_eligible": true,
		"events_date": bson.M{"$lt": currentTime},
	}

	if achievementType != "" {
		match["events_rewilding_achievement_type"] = achievementType
	}

	if starType != "" {
		starTypeInt, _ := strconv.Atoi(starType)
		match["EventParticipants.event_participants_star_type"] = starTypeInt
	}

	if country != "" {
		match["events_country_code"] = country
	}

	filterStage := bson.D{{Key: "$match", Value: match}}

	pipeline := mongo.Pipeline{
		lookupStage,
	}

	if !helpers.MongoZeroID(userId) && otherUserId != "" {
		pipeline = append(pipeline, lookupOtherUser, bson.D{
			{Key: "$unwind", Value: "$OtherParticipant"},
		})
	}

	groupStage := bson.D{
		{Key: "$group", Value: bson.M{
			"_id":                        "$events_rewilding",
			"rewilding_count":            bson.M{"$sum": 1},
			"rewilding_star_type":        bson.M{"$min": "$EventParticipants.event_participants_star_type"},
			"rewilding_star_unlocked_at": bson.M{"$min": "$EventParticipants.event_participants_achievement_unlocked_at"},
		}},
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
					bson.M{"rewilding_star_type": "$rewilding_star_type"},
					bson.M{"rewilding_star_unlocked_at": "$rewilding_star_unlocked_at"},
				},
			},
		},
	}}

	pipeline = append(pipeline, unwindStage, filterStage, groupStage, rewildLookup, rewildUnwind, replaceRoot)
	cursor, err := config.DB.Collection("Events").Aggregate(context.TODO(), pipeline)
	cursor.All(context.TODO(), results)

	return err
}

func (t AchievementRepository) Places(c *gin.Context) {
	userDetail := helpers.GetAuthUser(c)
	achievementType := c.Query("achievement_type")
	var AchievementPlaces []models.AchievementPlaces

	filter := bson.M{}
	if achievementType != "" {
		filter["ref_achievement_places_type"] = achievementType
	}

	matchStage := bson.D{{Key: "$match", Value: filter}}

	pipeline := mongo.Pipeline{
		matchStage,
		bson.D{{Key: "$lookup", Value: bson.M{
			"as":   "Events",
			"from": "Events",
			"let":  bson.M{"ref_achievement_places_id": "$_id"},
			"pipeline": []bson.M{
				{
					"$match": bson.M{
						"$expr": bson.M{
							"$eq": bson.A{"$events_rewilding_achievement_type_id", "$$ref_achievement_places_id"},
						},
					},
				},
				{"$match": bson.M{"events_rewilding_achievement_eligible": true}},
				{
					"$lookup": bson.M{
						"as":   "Participant",
						"from": "EventParticipants",
						"let":  bson.M{"event_id": "$_id"},
						"pipeline": []bson.M{
							{"$match": bson.M{
								"$expr": bson.M{
									"$eq": bson.A{"$event_participants_event", "$$event_id"}},
							}},
							{"$match": bson.M{"event_participants_user": userDetail.UsersId}},
						},
					},
				},
				{
					"$set": bson.M{"ref_achievement_places_count": bson.M{"$size": "$Participant"}},
				},
				{
					"$group": bson.M{
						"_id":   "events_rewilding_achievement_type_id",
						"count": bson.M{"$sum": "$ref_achievement_places_count"},
					},
				},
			},
		}}},
		bson.D{
			{Key: "$unwind", Value: bson.M{
				"path":                       "$Events",
				"preserveNullAndEmptyArrays": true,
			}},
		},
		bson.D{
			{Key: "$set", Value: bson.M{
				"ref_achievement_places_count": "$Events.count"},
			},
		},
	}

	cursor, _ := config.DB.Collection("RefAchievementPlaces").Aggregate(context.TODO(), pipeline)
	cursor.All(context.TODO(), &AchievementPlaces)

	if len(AchievementPlaces) == 0 {
		helpers.ResponseNoData(c, "No Data")
		return
	}
	c.JSON(200, AchievementPlaces)
}
