package middleware

import (
	"net/http"
	"oosa_rewild/internal/auth"
	"oosa_rewild/internal/helpers"
	"strings"

	"github.com/gin-gonic/gin"
)

type UserBindByHeader struct {
	Id       string `header:"X-User-Id"`
	User     string `header:"X-User-Account"`
	Email    string `header:"X-User-Email"`
	Name     string `header:"X-User-Name"`
	Language string `header:"X-User-Language"`
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if ssoAuth(c) {
			return
		}
		reqToken := c.Request.Header.Get("Authorization")
		if reqToken == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "AUTH01-REWILDING: You are not authorized to access this resource"})
			c.Abort()
			return
		}

		splitToken := strings.Split(reqToken, "Bearer ")
		reqToken = splitToken[1]

		user, err := auth.VerifyToken(reqToken)

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "AUTH02-REWILDING: You are not authorized to access this resource"})
			c.Abort()
			return
		}

		if helpers.MongoZeroID(user.UsersId) {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "AUTH03-REWILDING: You are not authorized to access this resource"})
			c.Abort()
			return
		}

		c.Set("user", &user)
		c.Next()
	}
}

func CheckIfAuth(c *gin.Context) {
	if ssoCheckIfAuth(c) {
		return
	}
	reqToken := c.Request.Header.Get("Authorization")
	if reqToken == "" {
		return
	}

	splitToken := strings.Split(reqToken, "Bearer ")
	reqToken = splitToken[1]

	user, err := auth.VerifyToken(reqToken)

	if err != nil {
		return
	}

	if helpers.MongoZeroID(user.UsersId) {
		return
	}

	c.Set("user", &user)
}

func ssoAuth(c *gin.Context) bool {
	if _, ok := c.Get("user"); ok {
		c.Next()
		return true
	}
	return false
}

func ssoCheckIfAuth(c *gin.Context) bool {
	if _, ok := c.Get("user"); ok {
		return true
	}
	return false
}
