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

type NewsRepository struct{}
type NewsRequest struct {
	NewsDate    string `json:"news_date" validate:"required,datetime=2006-01-02"`
	NewsTitle   string `json:"news_title" validate:"required"`
	NewsContent string `json:"news_content" validate:"required"`
}

func (r NewsRepository) Retrieve(c *gin.Context) {
	var results []models.News
	filters := bson.D{}
	cursor, err := config.DB.Collection("News").Find(context.TODO(), filters)
	if err != nil {
		panic(err)
	}
	cursor.All(context.TODO(), &results)
	if len(results) == 0 {
		helpers.ResponseNoData(c, "No Data")
		return
	}
	c.JSON(http.StatusOK, results)
}

func (r NewsRepository) Create(c *gin.Context) {
	userDetail := helpers.GetAuthUser(c)
	var payload NewsRequest
	validateError := helpers.Validate(c, &payload)
	if validateError != nil {
		return
	}

	insert := models.News{
		NewsCreatedAt: primitive.NewDateTimeFromTime(time.Now()),
		NewsCreatedBy: userDetail.UsersId,
	}
	r.ProcessData(c, &insert, payload)

	result, err := config.DB.Collection("News").InsertOne(context.TODO(), insert)
	if err != nil {
		fmt.Println("ERROR", err.Error())
		return
	}

	var News models.News
	config.DB.Collection("News").FindOne(context.TODO(), bson.D{{Key: "_id", Value: result.InsertedID}}).Decode(&News)
	c.JSON(http.StatusOK, News)
}

func (r NewsRepository) Read(c *gin.Context) {
	var News models.News
	err := r.ReadOne(c, &News, "")
	if err == nil {
		c.JSON(http.StatusOK, News)
	}
}

func (r NewsRepository) Update(c *gin.Context) {
	//userDetail := user.(*models.Users)
	var payload NewsRequest
	validateError := helpers.Validate(c, &payload)
	if validateError != nil {
		return
	}
	var News models.News
	errRead := r.ReadOne(c, &News, "")

	if errRead != nil {
		return
	}

	r.ProcessData(c, &News, payload)
	_, errUpd := config.DB.Collection("News").ReplaceOne(context.TODO(), bson.D{{Key: "_id", Value: News.NewsId}}, News)
	if errUpd == nil {
		c.JSON(http.StatusOK, News)
	}
}

func (r NewsRepository) ProcessData(c *gin.Context, News *models.News, payload NewsRequest) {
	newsDate := helpers.StringDateToPrimitiveDateTime(payload.NewsDate)
	News.NewsDate = newsDate
	News.NewsTitle = payload.NewsTitle
	News.NewsContent = payload.NewsContent
}

func (r NewsRepository) ReadOne(c *gin.Context, PocketList *models.News, newsId string) error {
	idVal := c.Param("id")
	if newsId != "" {
		idVal = newsId
	}

	id, _ := primitive.ObjectIDFromHex(idVal)
	return r.ReadById(c, PocketList, id)
}

func (r NewsRepository) ReadById(c *gin.Context, PocketList *models.News, id primitive.ObjectID) error {
	filter := bson.D{{Key: "_id", Value: id}}
	err := config.DB.Collection("News").FindOne(context.TODO(), filter).Decode(&PocketList)
	if err != nil {
		helpers.ResultNotFound(c, err, "News not found")
	}
	return err
}

func (r NewsRepository) Delete(c *gin.Context) {
	var News models.News
	err := r.ReadOne(c, &News, "")
	if err == nil {
		filters := bson.D{{Key: "_id", Value: News.NewsId}}
		config.DB.Collection("News").DeleteOne(context.TODO(), filters)
		helpers.ResultMessageSuccess(c, "News deleted")
	}
}
