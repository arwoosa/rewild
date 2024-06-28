package middleware

import (
	"net/http"
	"oosa_rewild/internal/models"

	"github.com/gin-gonic/gin"
)

func AuthBusinessUserMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := c.Get("user")

		if !exists {
			c.JSON(http.StatusBadRequest, gin.H{"message": "AUTH-BUSINESS-REWILDING01: Invalid user"})
			c.Abort()
			return
		}

		userDetail := user.(*models.Users)

		if !userDetail.UsersIsBusiness {
			c.JSON(http.StatusBadRequest, gin.H{"message": "AUTH-BUSINESS-REWILDING02: Not a business user"})
			c.Abort()
			return
		}

		c.Next()
	}
}
