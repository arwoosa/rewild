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

type CollaborativeLogRepository struct{}

func (r CollaborativeLogRepository) Retrieve(c *gin.Context) {
	authUserId := helpers.GetAuthUser(c).UsersId            // 查詢者id
	userId := helpers.StringToPrimitiveObjId(c.Param("id")) // 被查詢者id
	rewildingId := helpers.StringToPrimitiveObjId(c.Query("rewilding_id"))
	var results []models.Events
	currentTime := primitive.NewDateTimeFromTime(time.Now())

	if authUserId != userId {
		agg := mongo.Pipeline{
			// 撈出共同參與的events
			bson.D{{
				Key: "$match", Value: bson.M{
					"$or": []bson.M{
						{"event_participants_user": authUserId},
						{"event_participants_user": userId},
					},
					"event_participants_status": GetEventParticipantStatus("ACCEPTED"),
				},
			}},
			bson.D{{
				Key: "$group", Value: bson.M{"_id": "$event_participants_event", "count": bson.M{"$sum": 1}},
			}},
			bson.D{{
				Key: "$match", Value: bson.M{"count": 2},
			}},
			bson.D{{
				Key: "$lookup", Value: bson.M{"from": "Events",
					"localField":   "_id",
					"foreignField": "_id",
					"as":           "Events",
				},
			}},
			bson.D{{
				Key: "$unset", Value: bson.A{"_id", "count"},
			}},
			bson.D{{
				Key: "$unwind", Value: bson.M{"path": "$Events"},
			}},
			bson.D{{
				Key: "$replaceRoot", Value: bson.M{"newRoot": "$Events"},
			}},
			// 過濾掉已刪除的events，並且留下該rewilding(地點)的events
			bson.D{{
				Key: "$match", Value: bson.M{
					"events_deleted":   bson.M{"$exists": false},
					"events_rewilding": rewildingId,
				},
			}},
			bson.D{{Key: "$addFields", Value: bson.M{
				"can_view_details": bson.M{
					"$lte": []interface{}{
						bson.M{"$dateTrunc": bson.M{"date": "$events_date_end", "unit": "day"}},
						bson.M{"$dateTrunc": bson.M{"date": currentTime, "unit": "day"}},
					},
				},
				"upload_polaroids": false, // 查看者他人行程，不顯示上傳拍立得按鈕
				"sort_group": bson.M{
					"$switch": bson.M{
						"branches": []bson.M{
							// A. 已開放上傳但尚未上傳
							{
								"case": bson.M{
									"$and": []interface{}{
										bson.M{"$lt": []interface{}{
											bson.M{"$dateTrunc": bson.M{"date": "$events_date_end", "unit": "day"}},
											bson.M{"$dateTrunc": bson.M{"date": currentTime, "unit": "day"}},
										}},
										bson.M{"$eq": []interface{}{bson.M{"$type": "$event_participants_star_type"}, "missing"}},
									},
								},
								"then": 0,
							},
							// B. 尚未開放上傳
							{
								"case": bson.M{
									"$gte": []interface{}{
										bson.M{"$dateTrunc": bson.M{"date": "$events_date_end", "unit": "day"}},
										bson.M{"$dateTrunc": bson.M{"date": currentTime, "unit": "day"}},
									},
								},
								"then": 1,
							},
							// C. 已上傳
							{
								"case": bson.M{
									"$ne": []interface{}{bson.M{"$type": "$event_participants_star_type"}, "missing"},
								},
								"then": 2,
							},
						},
						"default": 3, // 其他情況
					},
				},
			}}},
			bson.D{{Key: "$addFields", Value: bson.M{
				"events_actions": bson.M{
					"can_view_details": "$can_view_details",
					"upload_polaroids": "$upload_polaroids",
				},
			}}},
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
			bson.D{{Key: "$addFields", Value: bson.M{
				"effective_primary_date": bson.M{
					"$cond": bson.M{
						"if":   bson.M{"$eq": bson.A{"$sort_group", 2}},
						"then": "$events_date",     // 已上傳拍立的，以行程開始日期優先排序
						"else": "$events_date_end", // 未上傳拍立的，以行程結束日期優先排序
					},
				},
				"effective_secondary_date": bson.M{
					"$cond": bson.M{
						"if":   bson.M{"$eq": bson.A{"$sort_group", 2}},
						"then": "$events_date_end",
						"else": "$events_date",
					},
				},
			}}},
			bson.D{{Key: "$sort", Value: bson.M{
				"sort_group":               1,
				"effective_primary_date":   -1,
				"effective_secondary_date": -1,
			}}},
		}

		cursor, err := config.DB.Collection("EventParticipants").Aggregate(context.TODO(), agg)

		if err != nil {
			helpers.ResponseBadRequestError(c, err.Error())
			return
		}
		cursor.All(context.TODO(), &results)
	} else { // 原本查自己的query
		filterEvent := bson.M{
			"events_deleted":   bson.M{"$exists": false},
			"events_rewilding": rewildingId,
		}

		lookupStage := bson.D{{Key: "$lookup", Value: bson.M{
			"as":   "EventParticipants",
			"from": "EventParticipants",
			"let":  bson.M{"event_id": "$_id"},
			"pipeline": []bson.M{
				{"$match": bson.M{
					"$expr": bson.M{
						"$eq": bson.A{"$event_participants_event", "$$event_id"}},
				}},
				{"$match": bson.M{
					"event_participants_user":   authUserId,
					"event_participants_status": GetEventParticipantStatus("ACCEPTED"),
				}},
			},
		}}}

		unwindStage := bson.D{
			{Key: "$unwind", Value: "$EventParticipants"},
		}

		agg := mongo.Pipeline{
			bson.D{{
				Key: "$match", Value: filterEvent,
			}},
			lookupStage,
			unwindStage,
			bson.D{{Key: "$addFields", Value: bson.M{
				"can_view_details": bson.M{
					"$lte": []interface{}{
						bson.M{"$dateTrunc": bson.M{"date": "$events_date_end", "unit": "day"}},
						bson.M{"$dateTrunc": bson.M{"date": currentTime, "unit": "day"}},
					},
				},
				"upload_polaroids": bson.M{
					"$and": []interface{}{
						bson.M{"$lte": []interface{}{
							bson.M{"$dateTrunc": bson.M{"date": "$events_date_end", "unit": "day"}},
							bson.M{"$dateTrunc": bson.M{"date": currentTime, "unit": "day"}},
						}},
						bson.M{"$eq": []interface{}{bson.M{"$type": "$EventParticipants.event_participants_star_type"}, "missing"}},
					},
				},
				"sort_group": bson.M{
					"$switch": bson.M{
						"branches": []bson.M{
							{
								"case": bson.M{
									"$and": []interface{}{
										bson.M{"$lt": []interface{}{
											bson.M{"$dateTrunc": bson.M{"date": "$events_date_end", "unit": "day"}},
											bson.M{"$dateTrunc": bson.M{"date": currentTime, "unit": "day"}},
										}},
										bson.M{"$eq": []interface{}{bson.M{"$type": "$EventParticipants.event_participants_star_type"}, "missing"}},
									},
								},
								"then": 0,
							},
							{
								"case": bson.M{
									"$gte": []interface{}{
										bson.M{"$dateTrunc": bson.M{"date": "$events_date_end", "unit": "day"}},
										bson.M{"$dateTrunc": bson.M{"date": currentTime, "unit": "day"}},
									},
								},
								"then": 1,
							},
							{
								"case": bson.M{
									"$ne": []interface{}{bson.M{"$type": "$EventParticipants.event_participants_star_type"}, "missing"},
								},
								"then": 2,
							},
						},
						"default": 3,
					},
				},
			}}},
			bson.D{{Key: "$addFields", Value: bson.M{
				"events_actions": bson.M{
					"can_view_details": "$can_view_details",
					"upload_polaroids": "$upload_polaroids",
				},
			}}},
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
			bson.D{{Key: "$addFields", Value: bson.M{
				"effective_primary_date": bson.M{
					"$cond": bson.M{
						"if":   bson.M{"$eq": bson.A{"$sort_group", 2}},
						"then": "$events_date",
						"else": "$events_date_end",
					},
				},
				"effective_secondary_date": bson.M{
					"$cond": bson.M{
						"if":   bson.M{"$eq": bson.A{"$sort_group", 2}},
						"then": "$events_date_end",
						"else": "$events_date",
					},
				},
			}}},
			bson.D{{Key: "$sort", Value: bson.M{
				"sort_group":               1,
				"effective_primary_date":   -1,
				"effective_secondary_date": -1,
			}}},
		}
		cursor, err := config.DB.Collection("Events").Aggregate(context.TODO(), agg)
		cursor.All(context.TODO(), &results)

		if err != nil {
			return
		}
	}

	if len(results) == 0 {
		helpers.ResponseNoData(c, "No Data")
		return
	}

	results = EventRepository{}.RetrieveParticipantDetails(results)
	c.JSON(http.StatusOK, results)
}

func (r CollaborativeLogRepository) ReadOne(c *gin.Context, Events *models.Events) error {
	id, _ := primitive.ObjectIDFromHex(c.Param("id"))
	filter := bson.D{{Key: "_id", Value: id}}
	err := config.DB.Collection("Events").FindOne(context.TODO(), filter).Decode(&Events)
	if err != nil {
		helpers.ResultEmpty(c, err)
	}
	return err
}
