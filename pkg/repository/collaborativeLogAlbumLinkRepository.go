package repository

import (
	"context"
	"errors"
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
	EventAlbumLinkVisibility int64  `json:"event_album_link_visibility"`
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

	userDetail := helpers.GetAuthUser(c)
	var payload CollaborativeLogAlbumLinkRequest
	validateError := helpers.Validate(c, &payload)
	if validateError != nil {
		return
	}

	insert := models.EventAlbumLink{
		EventAlbumLinkEvent:     Events.EventsId,
		EventAlbumLinkCreatedBy: userDetail.UsersId,
		EventAlbumLinkCreatedAt: primitive.NewDateTimeFromTime(time.Now()),
	}

	errProcess := r.ProcessData(&insert, payload)
	if errProcess != nil {
		helpers.ResponseBadRequestError(c, errProcess.Error())
		return
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

func (r CollaborativeLogAlbumLinkRepository) Update(c *gin.Context) {
	var Events models.Events
	err := CollaborativeLogRepository{}.ReadOne(c, &Events)
	if err != nil {
		return
	}

	// userDetail := helpers.GetAuthUser(c)
	var payload CollaborativeLogAlbumLinkRequest
	validateError := helpers.Validate(c, &payload)
	if validateError != nil {
		return
	}

	var EventAlbumLink models.EventAlbumLink
	albumLinkId := helpers.StringToPrimitiveObjId(c.Param("albumLinkId"))
	filter := bson.D{{Key: "_id", Value: albumLinkId}, {Key: "event_album_link_event", Value: Events.EventsId}}
	errAlbumLink := config.DB.Collection("EventAlbumLink").FindOne(context.TODO(), filter).Decode(&EventAlbumLink)
	if errAlbumLink != nil {
		helpers.ResultEmpty(c, err)
		return
	}

	r.ProcessData(&EventAlbumLink, payload)

	filters := bson.D{{Key: "_id", Value: EventAlbumLink.EventAlbumLinkId}, {Key: "event_album_link_event", Value: EventAlbumLink.EventAlbumLinkEvent}}
	upd := bson.D{{Key: "$set", Value: EventAlbumLink}}
	config.DB.Collection("EventAlbumLink").UpdateOne(context.TODO(), filters, upd)

	c.JSON(200, EventAlbumLink)
}

func (r CollaborativeLogAlbumLinkRepository) ProcessData(EventAlbumLink *models.EventAlbumLink, payload CollaborativeLogAlbumLinkRequest) error {
	EventAlbumLink.EventAlbumLinkAlbumUrl = payload.EventAlbumLinkAlbumUrl
	EventAlbumLink.EventAlbumLinkVisibility = &payload.EventAlbumLinkVisibility

	allowedString := []string{"photos.google.com", "icloud.com", "flickr.com", "mega.com", "mega.nz"}
	url, err := helpers.GetDomain(payload.EventAlbumLinkAlbumUrl)

	if err != nil {
		return err
	}

	if !helpers.StringInSlice(url, allowedString) {
		return errors.New("目前支援的有: Google Photo, iCloud, flickr, Mega Passed domain: " + url)
	}

	return nil
}
