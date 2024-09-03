package repository

import (
	"context"
	"oosa_rewild/internal/config"
	"oosa_rewild/internal/helpers"
	"oosa_rewild/internal/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CollaborativeLogExperienceRepository struct{}
type CollaborativeLogExperienceRequest struct {
	EventsExperience string `json:"events_experience" validate:"required"`
}

func (r CollaborativeLogExperienceRepository) Create(c *gin.Context) {
	var EventParticipants models.EventParticipants

	var payload CollaborativeLogExperienceRequest
	validateError := helpers.Validate(c, &payload)
	if validateError != nil {
		return
	}

	userDetail := helpers.GetAuthUser(c)
	id, _ := primitive.ObjectIDFromHex(c.Param("id"))
	filter := bson.D{
		{Key: "event_participants_event", Value: id},
		{Key: "event_participants_user", Value: userDetail.UsersId},
	}
	err := config.DB.Collection("EventParticipants").FindOne(context.TODO(), filter).Decode(&EventParticipants)
	if err != nil {
		helpers.ResultEmpty(c, err)
	}

	EventParticipants.EventParticipantsExperience = payload.EventsExperience
	updFilter := bson.D{
		{Key: "_id", Value: EventParticipants.EventParticipantsExperience},
	}
	upd := bson.D{{Key: "$set", Value: EventParticipants}}
	config.DB.Collection("EventParticipants").UpdateOne(context.TODO(), updFilter, upd)

	c.JSON(200, EventParticipants)
}
