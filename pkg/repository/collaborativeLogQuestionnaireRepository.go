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

type CollaborativeLogQuestionnaireRepository struct{}
type CollaborativeLogQuestionnaireRequest struct {
	EventsQuestionnaireLink string `json:"events_questionnaire_link" validate:"required"`
}

func (r CollaborativeLogQuestionnaireRepository) Create(c *gin.Context) {
	var Events models.Events
	err := CollaborativeLogRepository{}.ReadOne(c, &Events)
	if err != nil {
		return
	}

	allowedString := []string{"forms.gle", "docs.google.com", "surveycake.com"}

	var payload CollaborativeLogQuestionnaireRequest
	validateError := helpers.Validate(c, &payload)
	if validateError != nil {
		return
	}
	url, err := helpers.GetDomain(payload.EventsQuestionnaireLink)

	if err != nil {
		helpers.ResponseBadRequestError(c, err.Error())
		return
	}

	if !helpers.StringInSlice(url, allowedString) {
		helpers.ResponseBadRequestError(c, "目前支援的有: Google 表單, SurveyCake. Passed domain: "+url)
		return
	}

	Events.EventsQuestionnaireLink = payload.EventsQuestionnaireLink

	filters := bson.D{{Key: "_id", Value: Events.EventsId}}
	upd := bson.D{{Key: "$set", Value: Events}}
	config.DB.Collection("Events").UpdateOne(context.TODO(), filters, upd)

	c.JSON(http.StatusOK, Events)
}
