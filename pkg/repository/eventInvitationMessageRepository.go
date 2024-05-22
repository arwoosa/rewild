package repository

import (
	"context"
	"net/http"
	"oosa_rewild/internal/config"
	"oosa_rewild/internal/helpers"
	"oosa_rewild/internal/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

type EventInvitationMessageRepository struct{}

type EventInvitationMessageRequest struct {
	EventsInvitationMessage string `json:"events_invitation_message" validate:"required"`
}

func (r EventInvitationMessageRepository) Update(c *gin.Context) {
	var Event models.Events
	err := EventRepository{}.ReadOne(c, &Event)
	if err != nil {
		return
	}

	//userDetail := helpers.GetAuthUser(c)
	var payload EventInvitationMessageRequest
	validateError := helpers.Validate(c, &payload)
	if validateError != nil {
		return
	}

	Event.EventsInvitationMessage = payload.EventsInvitationMessage

	filters := bson.D{{Key: "_id", Value: Event.EventsId}}
	upd := bson.D{{Key: "$set", Value: Event}}
	config.DB.Collection("Events").UpdateOne(context.TODO(), filters, upd)
	c.JSON(http.StatusOK, Event)
}
