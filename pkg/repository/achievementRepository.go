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

// deprecated
func (t AchievementRepository) Retrieve(c *gin.Context) {
	c.Header("Deprecation", "true")
	c.Header("Warning", "299 - 'This API is deprecated, please use /user/{userId}/achievement'")
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

// userId: 被查看者
func (t AchievementRepository) GetAchievementsByUserIdV2(c *gin.Context, userId primitive.ObjectID, results *[]models.AchievementRewildingV2) error {
	country := c.Query("country")
	currentTime := primitive.NewDateTimeFromTime(time.Now())

	// filter 已刪除的 event
	deletedFilterStage := bson.D{{Key: "$match", Value: bson.M{
		"events_deleted": bson.M{"$ne": 1},
	}}}

	// Lookup(EventParticipants) 留下 被查看者 有參與的 event
	lookupParticipants := bson.D{{Key: "$lookup", Value: bson.M{
		"as":   "EventParticipants",
		"from": "EventParticipants",
		"let":  bson.M{"event_id": "$_id"},
		"pipeline": []bson.M{
			{"$match": bson.M{
				"$expr": bson.M{
					"$eq": bson.A{"$event_participants_event", "$$event_id"},
				},
			}},
			{"$match": bson.M{"event_participants_user": userId}},
			{"$match": bson.M{"event_participants_status": GetEventParticipantStatus("ACCEPTED")}},
		},
	}}}

	unwindParticipants := bson.D{{Key: "$unwind", Value: "$EventParticipants"}}

	match := bson.M{}
	if country != "" {
		match["events_country_code"] = country
	}
	filterStage := bson.D{{Key: "$match", Value: match}}

	groupStage := bson.D{
		{Key: "$group", Value: bson.M{
			"_id":             "$events_rewilding",
			"rewilding_count": bson.M{"$sum": 1},
			"all_star_types": bson.M{"$push": bson.M{
				"$cond": bson.M{
					"if": bson.M{
						"$and": bson.A{
							bson.M{"$ne": bson.A{"$EventParticipants.event_participants_star_type", nil}},
							bson.M{"$ne": bson.A{
								bson.M{"$type": "$EventParticipants.event_participants_star_type"},
								"missing",
							}},
						},
					},
					"then": "$EventParticipants.event_participants_star_type",
					"else": nil,
				},
			}},
			"all_event_ends": bson.M{"$push": bson.M{
				"star_type": bson.M{
					"$cond": bson.M{
						"if": bson.M{
							"$and": bson.A{
								bson.M{"$ne": bson.A{"$EventParticipants.event_participants_star_type", nil}},
								bson.M{"$ne": bson.A{
									bson.M{"$type": "$EventParticipants.event_participants_star_type"},
									"missing",
								}},
							},
						},
						"then": "$EventParticipants.event_participants_star_type",
						"else": nil,
					},
				},
				"event_end": "$events_date_end",
			}},
		}},
	}

	rewildLookup := bson.D{{Key: "$lookup", Value: bson.M{
		"from":         "Rewilding",
		"localField":   "_id",
		"foreignField": "_id",
		"as":           "RewildingDetail",
	}}}

	rewildUnwind := bson.D{{Key: "$unwind", Value: "$RewildingDetail"}}

	addComputedFields := bson.D{{
		Key: "$addFields", Value: bson.M{
			// 所有event皆有star_type=都上傳拍立得 => yellow, 反之則為白色
			"achievement_star_status": bson.M{
				"$cond": bson.M{
					"if": bson.M{
						"$allElementsTrue": bson.M{
							"$map": bson.M{
								"input": "$all_star_types",
								"as":    "type",
								"in": bson.M{
									"$ne": bson.A{"$$type", nil},
								},
							},
						},
					},
					"then": "yellow",
					"else": "white",
				},
			},
			// 從所有活動中，篩出「未得星星（star_type 為 nil）且活動已結束（event_end < 現在時間）」的活動
			// 找出其中結束時間最晚的一筆。
			"achievement_latest_can_upload_time": bson.M{
				"$let": bson.M{
					"vars": bson.M{
						"uploadables": bson.M{
							"$filter": bson.M{
								"input": "$all_event_ends",
								"as":    "event",
								"cond": bson.M{
									"$and": bson.A{
										bson.M{"$eq": bson.A{"$$event.star_type", nil}},
										bson.M{"$lt": bson.A{"$$event.event_end", currentTime}},
									},
								},
							},
						},
					},
					"in": bson.M{
						"$let": bson.M{
							"vars": bson.M{
								"latest": bson.M{"$max": "$$uploadables.event_end"},
							},
							"in": "$$latest",
						},
					},
				},
			},
			"achievement_shine": bson.M{
				"$cond": bson.M{
					"if":   bson.M{"$ne": bson.A{"$achievement_latest_can_upload_time", nil}},
					"then": true,
					"else": false,
				},
			},
		},
	}}

	replaceRoot := bson.D{{
		Key: "$replaceRoot", Value: bson.M{
			"newRoot": bson.M{
				"$mergeObjects": bson.A{
					"$RewildingDetail",
					bson.M{"rewilding_count": "$rewilding_count"},
					bson.M{"achievement_star_status": "$achievement_star_status"},
					bson.M{"achievement_latest_can_upload_time": "$achievement_latest_can_upload_time"},
					bson.M{"achievement_shine": "$achievement_shine"},
				},
			},
		},
	}}

	pipeline := mongo.Pipeline{
		deletedFilterStage, // 過濾已刪除的 event
		lookupParticipants, // 留下 擁有者有參加的 event
		unwindParticipants,
		filterStage,  // 根據API參數 篩選 country, star_type, achievementType，後兩個是舊規格
		groupStage,   // 依rewilding group event
		rewildLookup, // lookup rewilding的資料
		rewildUnwind,
		addComputedFields, // 轉換邏輯：achievement_latest_can_upload_time及achievement_star_status
		replaceRoot,
	}

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
