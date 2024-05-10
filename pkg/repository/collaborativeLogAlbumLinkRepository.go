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
	cursor, err := config.DB.Collection("EventAlbumLink").Find(context.TODO(), bson.D{})
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
