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

type EventReferenceLinksRepository struct{}
type EventReferenceLinksRequest struct {
	EventReferenceLinksLink  string `json:"event_reference_links_link" validate:"required"`
	EventReferenceLinksTitle string `json:"event_reference_links_title" validate:"required"`
}

func (r EventReferenceLinksRepository) Retrieve(c *gin.Context) {
	err := EventRepository{}.ReadOne(c, &models.Events{})
	if err != nil {
		return
	}

	var results []models.EventReferenceLinks
	cursor, err := config.DB.Collection("EventReferenceLinks").Find(context.TODO(), bson.D{})
	cursor.All(context.TODO(), &results)

	if err != nil {
		return
	}

	if len(results) == 0 {
		helpers.ResponseNoData(c, "No Data")
		return
	}
	c.JSON(http.StatusOK, results)
}

func (r EventReferenceLinksRepository) Create(c *gin.Context) {
	err := EventRepository{}.ReadOne(c, &models.Events{})
	if err != nil {
		return
	}

	userDetail := helpers.GetAuthUser(c)
	var payload EventReferenceLinksRequest
	validateError := helpers.Validate(c, &payload)
	if validateError != nil {
		return
	}

	insert := models.EventReferenceLinks{
		EventReferenceLinksEvent:     helpers.StringToPrimitiveObjId(c.Param("id")),
		EventReferenceLinksCreatedBy: userDetail.UsersId,
		EventReferenceLinksCreatedAt: primitive.NewDateTimeFromTime(time.Now()),
	}

	r.ProcessData(c, &insert, payload)

	result, err := config.DB.Collection("EventReferenceLinks").InsertOne(context.TODO(), insert)
	if err != nil {
		fmt.Println("ERROR", err.Error())
		return
	}

	var EventReferenceLinks models.EventReferenceLinks
	config.DB.Collection("EventReferenceLinks").FindOne(context.TODO(), bson.D{{Key: "_id", Value: result.InsertedID}}).Decode(&EventReferenceLinks)
	c.JSON(http.StatusOK, EventReferenceLinks)
}

func (r EventReferenceLinksRepository) Read(c *gin.Context) {
	err := EventRepository{}.ReadOne(c, &models.Events{})
	if err != nil {
		return
	}

	var EventReferenceLinks models.EventReferenceLinks
	errMb := r.ReadOne(c, &EventReferenceLinks)
	if errMb == nil {
		c.JSON(http.StatusOK, EventReferenceLinks)
	}
}

func (r EventReferenceLinksRepository) ReadOne(c *gin.Context, EventReferenceLinks *models.EventReferenceLinks) error {
	eventId := helpers.StringToPrimitiveObjId(c.Param("id"))
	EventReferenceLinksId := helpers.StringToPrimitiveObjId(c.Param("referenceLinkId"))
	filter := bson.D{{Key: "_id", Value: EventReferenceLinksId}, {Key: "event_reference_links_event", Value: eventId}}
	err := config.DB.Collection("EventReferenceLinks").FindOne(context.TODO(), filter).Decode(&EventReferenceLinks)
	if err != nil {
		helpers.ResultEmpty(c, err)
	}
	return err
}

func (r EventReferenceLinksRepository) Update(c *gin.Context) {
	err := EventRepository{}.ReadOne(c, &models.Events{})
	if err != nil {
		return
	}

	var payload EventReferenceLinksRequest
	validateError := helpers.Validate(c, &payload)
	if validateError != nil {
		return
	}

	var EventReferenceLinks models.EventReferenceLinks
	errMb := r.ReadOne(c, &EventReferenceLinks)
	if errMb == nil {
		r.ProcessData(c, &EventReferenceLinks, payload)
		filters := bson.D{{Key: "_id", Value: EventReferenceLinks.EventReferenceLinksId}, {Key: "event_reference_links_event", Value: EventReferenceLinks.EventReferenceLinksEvent}}
		upd := bson.D{{Key: "$set", Value: EventReferenceLinks}}
		config.DB.Collection("EventReferenceLinks").UpdateOne(context.TODO(), filters, upd)
		c.JSON(http.StatusOK, EventReferenceLinks)
	}
}

func (r EventReferenceLinksRepository) Delete(c *gin.Context) {
	err := EventRepository{}.ReadOne(c, &models.Events{})
	if err != nil {
		return
	}

	var EventReferenceLinks models.EventReferenceLinks
	errMb := r.ReadOne(c, &EventReferenceLinks)
	if errMb == nil {
		filters := bson.D{{Key: "_id", Value: EventReferenceLinks.EventReferenceLinksId}}
		config.DB.Collection("EventReferenceLinks").DeleteOne(context.TODO(), filters)
		helpers.ResultMessageSuccess(c, "Reference link deleted")
	}
}

func (r EventReferenceLinksRepository) ProcessData(c *gin.Context, EventReferenceLinks *models.EventReferenceLinks, payload EventReferenceLinksRequest) {
	EventReferenceLinks.EventReferenceLinksLink = payload.EventReferenceLinksLink
	EventReferenceLinks.EventReferenceLinksTitle = payload.EventReferenceLinksTitle
}
