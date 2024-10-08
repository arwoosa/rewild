package repository

import (
	"net/http"
	"oosa_rewild/internal/helpers"
	"oosa_rewild/internal/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OosaUserEventRepository struct{}

func (r OosaUserEventRepository) Retrieve(c *gin.Context) {
	userIdVal := c.Param("id")
	userId, _ := primitive.ObjectIDFromHex(userIdVal)
	var results []models.Events

	err := UserEventRepository{}.GetEventByUserId(c, userId, &results)

	if err != nil {
		return
	}

	if len(results) == 0 {
		helpers.ResponseNoData(c, "No Data")
		return
	}

	c.JSON(http.StatusOK, results)
}
