package repository

import (
	"context"
	"fmt"
	"net/http"
	"oosa_rewild/internal/config"
	"oosa_rewild/internal/helpers"
	"oosa_rewild/internal/models"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type EventAccountingRepository struct{}
type EventAccountingRequest struct {
	EventAccountingMessage string  `json:"event_accounting_message" validate:"required"`
	EventAccountingAmount  float64 `json:"event_accounting_amount" validate:"required,max=999999"`
	EventAccountingPaidBy  string  `json:"event_accounting_paid_by" validate:"required"`
}
type EventAccountingRequestBulkRequestValidate struct {
	Data []EventAccountingRequest `json:"data" validate:"required,dive"`
}

func (r EventAccountingRepository) Retrieve(c *gin.Context) {
	eventId := helpers.StringToPrimitiveObjId(c.Param("id"))
	err := EventRepository{}.ReadOne(c, &models.Events{})
	if err != nil {
		return
	}

	agg := mongo.Pipeline{
		bson.D{{
			Key: "$match", Value: bson.M{"event_accounting_event": eventId},
		}},
		bson.D{{
			Key: "$lookup", Value: bson.M{
				"from":         "Users",
				"localField":   "event_accounting_created_by",
				"foreignField": "_id",
				"as":           "event_accounting_created_by_user",
			},
		}},
		bson.D{{
			Key: "$unwind", Value: "$event_accounting_created_by_user",
		}},
	}

	var results []models.EventAccounting
	cursor, err := config.DB.Collection("EventAccounting").Aggregate(context.TODO(), agg)
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

func (r EventAccountingRepository) Create(c *gin.Context) {
	eventId := helpers.StringToPrimitiveObjId(c.Param("id"))
	err := EventRepository{}.ReadOne(c, &models.Events{})
	if err != nil {
		return
	}

	userDetail := helpers.GetAuthUser(c)
	var payload []EventAccountingRequest

	errBind := c.BindJSON(&payload)
	if errBind != nil {
		helpers.ResponseBadRequestError(c, err.Error())
		return
	}

	payloadValidate := EventAccountingRequestBulkRequestValidate{
		Data: payload,
	}

	validateError := helpers.ValidateObject(c, payloadValidate)
	if validateError != nil {
		return
	}

	lenAccounting := len(payload)
	if lenAccounting > int(config.APP_LIMIT.EventAccountingLimit) {
		errMessage := "Unable to add accounting records. Maximum allowed: " + strconv.Itoa(int(config.APP_LIMIT.EventAccountingLimit))
		helpers.ResponseBadRequestError(c, errMessage)
		return
	}

	var strLenErr []string
	var insertAccounting []interface{}
	for k, v := range payload {
		match, errMessage := helpers.ValidateStringLength(v.EventAccountingMessage, int(config.APP_LIMIT.LengthEventAccountingMessage))
		fmt.Println(match, v.EventAccountingMessage, int(config.APP_LIMIT.LengthEventAccountingMessage))
		if !match {
			strLenErr = append(strLenErr, "Accounting message at index "+strconv.Itoa(k)+" can only contain "+errMessage)
		} else {
			var EventParticipants models.EventParticipants
			filterParticipant := bson.D{
				{Key: "event_participants_event", Value: eventId},
				{Key: "event_participants_user", Value: helpers.StringToPrimitiveObjId(v.EventAccountingPaidBy)},
			}
			filterParticipantErr := config.DB.Collection("EventParticipants").FindOne(context.TODO(), filterParticipant).Decode(&EventParticipants)

			if filterParticipantErr != nil {
				strLenErr = append(strLenErr, "Payee at index "+strconv.Itoa(k)+" is not an event participant")
				continue
			}

			insertAccounting = append(insertAccounting, models.EventAccounting{
				EventAccountingEvent:     eventId,
				EventAccountingMessage:   v.EventAccountingMessage,
				EventAccountingAmount:    v.EventAccountingAmount,
				EventAccountingPaidBy:    helpers.StringToPrimitiveObjId(v.EventAccountingPaidBy),
				EventAccountingCreatedBy: userDetail.UsersId,
				EventAccountingCreatedAt: primitive.NewDateTimeFromTime(time.Now()),
			})
		}
	}

	if len(strLenErr) > 0 {
		helpers.ResponseBadRequestError(c, strings.Join(strLenErr, ", "))
		return
	}

	filterDelete := bson.D{
		{Key: "event_accounting_event", Value: eventId},
	}
	config.DB.Collection("EventAccounting").DeleteMany(context.TODO(), filterDelete)
	config.DB.Collection("EventAccounting").InsertMany(context.TODO(), insertAccounting)

	r.Retrieve(c)
}

func (r EventAccountingRepository) Read(c *gin.Context) {
	err := EventRepository{}.ReadOne(c, &models.Events{})
	if err != nil {
		return
	}

	var EventAccounting models.EventAccounting
	errMb := r.ReadOne(c, &EventAccounting)
	if errMb == nil {
		c.JSON(http.StatusOK, EventAccounting)
	}
}

func (r EventAccountingRepository) ReadOne(c *gin.Context, EventAccounting *models.EventAccounting) error {
	eventId := helpers.StringToPrimitiveObjId(c.Param("id"))
	EventAccountingId := helpers.StringToPrimitiveObjId(c.Param("accountingId"))
	filter := bson.D{{Key: "_id", Value: EventAccountingId}, {Key: "event_accounting_event", Value: eventId}}
	err := config.DB.Collection("EventAccounting").FindOne(context.TODO(), filter).Decode(&EventAccounting)
	if err != nil {
		helpers.ResultEmpty(c, err)
	}
	return err
}

func (r EventAccountingRepository) Update(c *gin.Context) {
	err := EventRepository{}.ReadOne(c, &models.Events{})
	if err != nil {
		return
	}

	var payload EventAccountingRequest
	validateError := helpers.Validate(c, &payload)
	if validateError != nil {
		return
	}

	match, errMessage := helpers.ValidateStringStyle1(payload.EventAccountingMessage, int(config.APP_LIMIT.LengthEventAccountingMessage))
	if !match {
		helpers.ResponseBadRequestError(c, "Name can only contain "+errMessage)
		return
	}

	var EventAccounting models.EventAccounting
	errMb := r.ReadOne(c, &EventAccounting)
	if errMb == nil {
		r.ProcessData(c, &EventAccounting, payload)
		filters := bson.D{{Key: "_id", Value: EventAccounting.EventAccountingId}, {Key: "event_accounting_event", Value: EventAccounting.EventAccountingEvent}}
		upd := bson.D{{Key: "$set", Value: EventAccounting}}
		config.DB.Collection("EventAccounting").UpdateOne(context.TODO(), filters, upd)
		r.Read(c)
		// c.JSON(http.StatusOK, EventAccounting)
	}
}

func (r EventAccountingRepository) Delete(c *gin.Context) {
	err := EventRepository{}.ReadOne(c, &models.Events{})
	if err != nil {
		return
	}

	var EventAccounting models.EventAccounting
	errMb := r.ReadOne(c, &EventAccounting)
	if errMb == nil {
		filters := bson.D{{Key: "_id", Value: EventAccounting.EventAccountingId}}
		config.DB.Collection("EventAccounting").DeleteOne(context.TODO(), filters)
		helpers.ResultMessageSuccess(c, "Accounting record deleted")
	}
}

func (r EventAccountingRepository) ProcessData(c *gin.Context, EventAccounting *models.EventAccounting, payload EventAccountingRequest) {
	paidBy := helpers.StringToPrimitiveObjId(payload.EventAccountingPaidBy)
	EventAccounting.EventAccountingMessage = payload.EventAccountingMessage
	EventAccounting.EventAccountingAmount = payload.EventAccountingAmount
	EventAccounting.EventAccountingPaidBy = paidBy
}
