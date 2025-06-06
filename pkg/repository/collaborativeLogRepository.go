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

// buildSortGroupExpression 構建排序群組的 MongoDB 表達式
// 根據拍立得數量和事件日期決定排序分組
func buildSortGroupExpression(currentTime primitive.DateTime) bson.M {
	return bson.M{
		"$switch": bson.M{
			"branches": []bson.M{
				// A. 已開放上傳且尚未上傳任何拍立得
				{
					"case": bson.M{
						"$and": []interface{}{
							bson.M{"$lte": []interface{}{
								bson.M{"$dateTrunc": bson.M{"date": "$events_date_end", "unit": "day"}},
								bson.M{"$dateTrunc": bson.M{"date": currentTime, "unit": "day"}},
							}},
							bson.M{"$eq": []interface{}{
								"$user_polaroid_count",
								0,
							}},
						},
					},
					"then": 0,
				},
				// B. 尚未開放上傳
				{
					"case": bson.M{
						"$gt": []interface{}{
							bson.M{"$dateTrunc": bson.M{"date": "$events_date_end", "unit": "day"}},
							bson.M{"$dateTrunc": bson.M{"date": currentTime, "unit": "day"}},
						},
					},
					"then": 1,
				},
				// C. 已上傳任何數量的拍立得
				{
					"case": bson.M{
						"$gt": []interface{}{
							"$user_polaroid_count",
							0, // 修改為大於0
						},
					},
					"then": 2,
				},
			},
			"default": 3, // 其他情況
		},
	}
}

// buildCanViewDetailsExpression 構建是否可查看詳情的 MongoDB 表達式
func buildCanViewDetailsExpression(currentTime primitive.DateTime) bson.M {
	return bson.M{
		"$lte": []interface{}{
			bson.M{"$dateTrunc": bson.M{"date": "$events_date_end", "unit": "day"}},
			bson.M{"$dateTrunc": bson.M{"date": currentTime, "unit": "day"}},
		},
	}
}

// buildUploadPolaroidsExpression 構建是否可上傳拍立得的 MongoDB 表達式
func buildUploadPolaroidsExpression(currentTime primitive.DateTime) bson.M {
	return bson.M{
		"$and": []interface{}{
			bson.M{"$lte": []interface{}{
				bson.M{"$dateTrunc": bson.M{"date": "$events_date_end", "unit": "day"}},
				bson.M{"$dateTrunc": bson.M{"date": currentTime, "unit": "day"}},
			}},
			bson.M{"$eq": []interface{}{
				"$user_polaroid_count",
				0, // 只有拍立得數量為0時才顯示上傳按鈕
			}},
		},
	}
}

// buildEffectiveDateFields 構建有效日期欄位的 MongoDB 表達式
func buildEffectiveDateFields() bson.M {
	return bson.M{
		"effective_primary_date": bson.M{
			"$cond": bson.M{
				"if":   bson.M{"$eq": bson.A{"$sort_group", 2}},
				"then": "$events_date",     // 已上傳的，以行程開始日期優先排序
				"else": "$events_date_end", // 未上傳的，以行程結束日期優先排序
			},
		},
		"effective_secondary_date": bson.M{
			"$cond": bson.M{
				"if":   bson.M{"$eq": bson.A{"$sort_group", 2}},
				"then": "$events_date_end",
				"else": "$events_date",
			},
		},
	}
}

// buildCommonAggregationStages 構建通用的聚合階段
func buildCommonAggregationStages(currentTime primitive.DateTime, uploadPolaroids interface{}) []bson.D {
	return []bson.D{
		// 添加基本欄位
		{{Key: "$addFields", Value: bson.M{
			"can_view_details": buildCanViewDetailsExpression(currentTime),
			"upload_polaroids": uploadPolaroids,
			"sort_group":       buildSortGroupExpression(currentTime),
		}}},
		// 添加事件操作欄位
		{{Key: "$addFields", Value: bson.M{
			"events_actions": bson.M{
				"can_view_details": "$can_view_details",
				"upload_polaroids": "$upload_polaroids",
			},
		}}},
		// 查找事件創建者資訊
		{{Key: "$lookup", Value: bson.M{
			"from":         "Users",
			"localField":   "events_created_by",
			"foreignField": "_id",
			"as":           "events_created_by_user",
		}}},
		{{Key: "$unwind", Value: "$events_created_by_user"}},
		// 添加有效日期欄位
		{{Key: "$addFields", Value: buildEffectiveDateFields()}},
		// 使用 $facet 來分別處理不同組別的排序
		{{Key: "$facet", Value: bson.M{
			"group_0": []bson.M{ // A組：已開放上傳但尚未上傳
				{"$match": bson.M{"sort_group": 0}},
				{"$sort": bson.M{
					"effective_primary_date":   -1, // 大到小
					"effective_secondary_date": -1,
				}},
			},
			"group_1": []bson.M{ // B組：尚未開放上傳
				{"$match": bson.M{"sort_group": 1}},
				{"$sort": bson.M{
					"effective_primary_date":   1, // 小到大
					"effective_secondary_date": 1,
				}},
			},
			"group_2": []bson.M{ // C組：已上傳
				{"$match": bson.M{"sort_group": 2}},
				{"$sort": bson.M{
					"effective_primary_date":   -1, // 大到小
					"effective_secondary_date": -1,
				}},
			},
		}}},
		// 合併結果並保持組別順序
		{{Key: "$project", Value: bson.M{
			"results": bson.M{
				"$concatArrays": bson.A{"$group_0", "$group_1", "$group_2"},
			},
		}}},
		{{Key: "$unwind", Value: "$results"}},
		{{Key: "$replaceRoot", Value: bson.M{"newRoot": "$results"}}},
	}
}

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
			// 添加 lookup EventPolaroids 來計算該用戶的拍立得數量
			bson.D{{Key: "$lookup", Value: bson.M{
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
			}}},
			// 計算 user_polaroid_count 欄位
			bson.D{{Key: "$addFields", Value: bson.M{
				"user_polaroid_count": bson.M{
					"$ifNull": bson.A{
						bson.M{"$arrayElemAt": bson.A{"$EventPolaroids.total", 0}},
						0,
					},
				},
			}}},
		}

		// 新增共用的Stage，查看他人行程不顯示上傳拍立得按鈕
		commonStages := buildCommonAggregationStages(currentTime, false)
		agg = append(agg, commonStages...)

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

		// 添加 lookup EventPolaroids 來計算該用戶的拍立得數量
		lookupPolaroids := bson.D{{Key: "$lookup", Value: bson.M{
			"as":   "EventPolaroids",
			"from": "EventPolaroids",
			"let":  bson.M{"event_id": "$_id", "user_id": authUserId},
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

		unwindStage := bson.D{
			{Key: "$unwind", Value: "$EventParticipants"},
		}

		// 計算 user_polaroid_count 欄位
		addPolaroidCountStage := bson.D{{Key: "$addFields", Value: bson.M{
			"user_polaroid_count": bson.M{
				"$ifNull": bson.A{
					bson.M{"$arrayElemAt": bson.A{"$EventPolaroids.total", 0}},
					0,
				},
			},
		}}}

		agg := mongo.Pipeline{
			bson.D{{
				Key: "$match", Value: filterEvent,
			}},
			lookupStage,
			lookupPolaroids,
			unwindStage,
			addPolaroidCountStage,
		}

		// 新增共用的Stage，查看自己的行程，依照顯示上傳拍立得按鈕的條件來決定
		commonStages := buildCommonAggregationStages(currentTime, buildUploadPolaroidsExpression(currentTime))
		agg = append(agg, commonStages...)

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
