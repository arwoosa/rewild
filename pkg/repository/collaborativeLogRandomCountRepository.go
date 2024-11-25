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

type CollaborativeLogRandomCountRepository struct{}
type CollaborativeLogRandomCountRequest struct {
	EventsQuestionnaireLink string `json:"events_questionnaire_link" validate:"required"`
}

func (r CollaborativeLogRandomCountRepository) Read(c *gin.Context) {
	userDetail := helpers.GetAuthUser(c)
	var Events models.Events
	var EventParticipants models.EventParticipants
	err := CollaborativeLogRepository{}.ReadOne(c, &Events)
	if err != nil {
		return
	}

	filter := bson.D{
		{Key: "event_participants_event", Value: Events.EventsId},
		{Key: "event_participants_user", Value: userDetail.UsersId},
	}
	errParticipant := config.DB.Collection("EventParticipants").FindOne(context.TODO(), filter).Decode(&EventParticipants)
	if errParticipant != nil {
		helpers.ResultEmpty(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"random_count": EventParticipants.EventParticipantsRandomCount})
}

func (r CollaborativeLogRandomCountRepository) Update(c *gin.Context) {
	userDetail := helpers.GetAuthUser(c)
	var Events models.Events
	var EventParticipants models.EventParticipants
	err := CollaborativeLogRepository{}.ReadOne(c, &Events)
	if err != nil {
		return
	}

	filter := bson.D{
		{Key: "event_participants_event", Value: Events.EventsId},
		{Key: "event_participants_user", Value: userDetail.UsersId},
	}
	errParticipant := config.DB.Collection("EventParticipants").FindOne(context.TODO(), filter).Decode(&EventParticipants)
	if errParticipant != nil {
		helpers.ResultEmpty(c, err)
		return
	}

	if EventParticipants.EventParticipantsRandomCount < 3 {
		filterUpd := bson.D{{Key: "_id", Value: EventParticipants.EventParticipantsId}}
		newCount := EventParticipants.EventParticipantsRandomCount + 1

		eventUpd := bson.D{{Key: "$set", Value: map[string]interface{}{
			"event_participants_random_count": newCount,
		}}}
		config.DB.Collection("EventParticipants").UpdateOne(context.TODO(), filterUpd, eventUpd)
		c.JSON(http.StatusOK, gin.H{"random_count": newCount})
	} else {
		helpers.ResponseBadRequestError(c, "Maximum random count of 3 reached")
	}
}
