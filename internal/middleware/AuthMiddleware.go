package middleware

import (
	"net/http"
	"oosa_rewild/internal/auth"
	"oosa_rewild/internal/helpers"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		reqToken := c.Request.Header.Get("Authorization")
		if reqToken == "" {
			c.JSON(http.StatusBadRequest, gin.H{"message": "AUTH01-REWILDING: You are not authorized to access this resource"})
			c.Abort()
			return
		}

		splitToken := strings.Split(reqToken, "Bearer ")
		reqToken = splitToken[1]

		user, err := auth.VerifyToken(reqToken)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "AUTH02-REWILDING: You are not authorized to access this resource"})
			c.Abort()
			return
		}

		if helpers.MongoZeroID(user.UsersId) {
			c.JSON(http.StatusBadRequest, gin.H{"message": "AUTH03-REWILDING: You are not authorized to access this resource"})
			c.Abort()
			return
		}

		c.Set("user", &user)
		c.Next()
	}
}
