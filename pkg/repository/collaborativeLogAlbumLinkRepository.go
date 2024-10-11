package repository

import (
	"context"
	"fmt"
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

type CollaborativeLogAlbumLinkRepository struct{}
type CollaborativeLogAlbumLinkRequest struct {
	EventAlbumLinkAlbumUrl   string `json:"event_album_link_album_url" validate:"required"`
	EventAlbumLinkVisibility int64  `json:"event_album_link_visibility" validate:"required"`
}

func (r CollaborativeLogAlbumLinkRepository) Retrieve(c *gin.Context) {
	var Events models.Events
	err := CollaborativeLogRepository{}.ReadOne(c, &Events)
	if err != nil {
		return
	}

	var EventAlbumLink []models.EventAlbumLink
	agg := mongo.Pipeline{
		bson.D{{
			Key: "$match", Value: bson.M{
				"event_album_link_event": Events.EventsId,
			},
		}},
		bson.D{{
			Key: "$lookup", Value: bson.M{
				"from":         "Users",
				"localField":   "event_album_link_created_by",
				"foreignField": "_id",
				"as":           "event_album_link_created_by_user",
			},
		}},
		bson.D{{
			Key: "$unwind", Value: "$event_album_link_created_by_user",
		}},
	}
	cursor, err := config.DB.Collection("EventAlbumLink").Aggregate(context.TODO(), agg)
	cursor.All(context.TODO(), &EventAlbumLink)

	if err != nil {
		return
	}

	if len(EventAlbumLink) == 0 {
		helpers.ResponseNoData(c, "No Data")
		return
	}
	c.JSON(http.StatusOK, EventAlbumLink)
}

func (r CollaborativeLogAlbumLinkRepository) Create(c *gin.Context) {
	var Events models.Events
	err := CollaborativeLogRepository{}.ReadOne(c, &Events)
	if err != nil {
		return
	}

	allowedString := []string{"photos.google.com", "icloud.com", "flickr.com", "mega.com", "mega.nz"}

	userDetail := helpers.GetAuthUser(c)
	var payload CollaborativeLogAlbumLinkRequest
	validateError := helpers.Validate(c, &payload)
	if validateError != nil {
		return
	}
	url, err := helpers.GetDomain(payload.EventAlbumLinkAlbumUrl)

	if err != nil {
		helpers.ResponseBadRequestError(c, err.Error())
		return
	}

	if !helpers.StringInSlice(url, allowedString) {
		helpers.ResponseBadRequestError(c, "目前支援的有: Google Photo, iCloud, flickr, Mega Passed domain: "+url)
		return
	}

	insert := models.EventAlbumLink{
		EventAlbumLinkEvent:      Events.EventsId,
		EventAlbumLinkAlbumUrl:   payload.EventAlbumLinkAlbumUrl,
		EventAlbumLinkVisibility: payload.EventAlbumLinkVisibility,
		EventAlbumLinkCreatedBy:  userDetail.UsersId,
		EventAlbumLinkCreatedAt:  primitive.NewDateTimeFromTime(time.Now()),
	}

	result, err := config.DB.Collection("EventAlbumLink").InsertOne(context.TODO(), insert)
	if err != nil {
		fmt.Println("ERROR", err.Error())
		return
	}

	var EventAlbumLink models.EventAlbumLink
	config.DB.Collection("EventAlbumLink").FindOne(context.TODO(), bson.D{{Key: "_id", Value: result.InsertedID}}).Decode(&EventAlbumLink)
	c.JSON(http.StatusOK, EventAlbumLink)
}
