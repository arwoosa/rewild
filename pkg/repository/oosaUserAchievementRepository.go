package repository

import (
	"net/http"
	"oosa_rewild/internal/helpers"
	"oosa_rewild/internal/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OosaUserAchievementRepository struct{}

func (r OosaUserAchievementRepository) Retrieve(c *gin.Context) {
	var results []models.EventsCountryCount
	userIdVal := c.Param("id")
	userId, _ := primitive.ObjectIDFromHex(userIdVal)

	err := AchievementRepository{}.GetAchievementsByUserId(c, userId, &results)
	if err != nil {
		helpers.ResponseError(c, err.Error())
		return
	}

	if len(results) == 0 {
		helpers.ResponseNoData(c, "No Data")
		return
	}

	c.JSON(http.StatusOK, results)
}
