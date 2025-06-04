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

	// Lookup(EventPolaroids) 計算該用戶在每個活動的拍立得數量
	lookupPolaroids := bson.D{{Key: "$lookup", Value: bson.M{
		"as":   "EventPolaroids",
		"from": "EventPolaroids",
		"let":  bson.M{"event_id": "$_id", "user_id": userId},
		"pipeline": []bson.M{
			{"$match": bson.M{
				"$expr": bson.M{
					"$and": bson.A{
						bson.M{"$eq": bson.A{"$event_polaroids_event", "$$event_id"}},
						bson.M{"$eq": bson.A{"$event_polaroids_created_by", "$$user_id"}},
					}},
			}},
			{"$count": "total"},
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
			"all_polaroid_counts": bson.M{"$push": bson.M{
				"$ifNull": bson.A{
					bson.M{"$arrayElemAt": bson.A{"$EventPolaroids.total", 0}}, // 取EventPolaroids陣列第一筆的total
					0, // EventPolaroids為空陣列，回傳0
				},
			}},
			"all_event_ends": bson.M{"$push": bson.M{
				"polaroid_count": bson.M{
					"$ifNull": bson.A{
						bson.M{"$arrayElemAt": bson.A{"$EventPolaroids.total", 0}},
						0,
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

	// 第一個 $addFields 階段，計算 achievement_star_status 和 achievement_latest_can_upload_time
	addPrimaryComputedFields := bson.D{{
		Key: "$addFields", Value: bson.M{
			// 所有event皆有拍立得上傳(count > 0) => yellow, 反之則為白色
			"achievement_star_status": bson.M{
				"$cond": bson.M{
					"if": bson.M{
						"$allElementsTrue": bson.M{
							"$map": bson.M{
								"input": "$all_polaroid_counts",
								"as":    "count",
								"in": bson.M{
									"$gt": bson.A{"$$count", 0},
								},
							},
						},
					},
					"then": "yellow",
					"else": "white",
				},
			},
			// 從所有活動中，篩出「未上傳拍立得（polaroid_count = 0）且活動已結束（event_end < 現在日期）」的活動
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
										bson.M{"$eq": bson.A{"$$event.polaroid_count", 0}},
										bson.M{"$lte": []interface{}{
											bson.M{"$dateTrunc": bson.M{"date": "$$event.event_end", "unit": "day"}},
											bson.M{"$dateTrunc": bson.M{"date": currentTime, "unit": "day"}},
										}},
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
		},
	}}

	// 第二個 $addFields 階段，計算 achievement_shine
	addShineField := bson.D{{
		Key: "$addFields", Value: bson.M{
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
		lookupPolaroids,    // 計算該用戶在每個活動的拍立得數量
		unwindParticipants,
		filterStage,  // 根據API參數 篩選 country
		groupStage,   // 依rewilding group event
		rewildLookup, // lookup rewilding的資料
		rewildUnwind,
		addPrimaryComputedFields,
		addShineField,
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
