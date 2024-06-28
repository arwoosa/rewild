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
	userDetail := helpers.GetAuthUser(c)
	userFilter := c.Query("user")
	var results []models.Events

	if userFilter != "" {
		if userFilter == userDetail.UsersId.Hex() {
			helpers.ResponseBadRequestError(c, "Cannot use own user ID")
			return
		}

		agg := mongo.Pipeline{
			bson.D{{
				Key: "$match", Value: bson.M{
					"$or": []bson.M{
						{"event_participants_user": userDetail.UsersId},
						{"event_participants_user": helpers.StringToPrimitiveObjId(userFilter)},
					},
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
				Key: "$unwind", Value: bson.M{"path": "$Events", "preserveNullAndEmptyArrays": true},
			}},
			bson.D{{
				Key: "$replaceRoot", Value: bson.M{"newRoot": "$Events"},
			}},
			bson.D{{
				Key: "$match", Value: bson.M{"events_date": bson.M{"$lte": primitive.NewDateTimeFromTime(time.Now())}},
			}},
		}

		cursor, err := config.DB.Collection("EventParticipants").Aggregate(context.TODO(), agg)

		if err != nil {
			helpers.ResponseBadRequestError(c, err.Error())
			return
		}
		cursor.All(context.TODO(), &results)
	} else {
		filter := bson.M{
			"events_date": bson.M{"$lte": primitive.NewDateTimeFromTime(time.Now())},
		}
		cursor, err := config.DB.Collection("Events").Find(context.TODO(), filter)
		cursor.All(context.TODO(), &results)

		if err != nil {
			return
		}
	}

	if len(results) == 0 {
		helpers.ResponseNoData(c, "No Data")
		return
	}
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
