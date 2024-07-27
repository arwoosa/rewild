package repository

import (
	"net/http"
	"oosa_rewild/internal/helpers"

	"github.com/gin-gonic/gin"
)

type ReferenceRepository struct{}

func (r ReferenceRepository) Options(c *gin.Context) {
	RefRewildingTypes := helpers.RefRewildingTypes()
	c.JSON(http.StatusOK, gin.H{"rewilding_types": RefRewildingTypes})
}
