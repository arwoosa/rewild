package helpers

import (
	"oosa_rewild/internal/models"

	"github.com/gin-gonic/gin"
)

func GetAuthUser(c *gin.Context) models.Users {
	user, _ := c.Get("user")
	userDetail := user.(*models.Users)

	return *userDetail
}
