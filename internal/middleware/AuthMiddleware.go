package middleware

import (
	"context"
	"net/http"
	"oosa_rewild/internal/auth"
	"oosa_rewild/internal/config"
	"oosa_rewild/internal/helpers"
	"oosa_rewild/internal/models"
	"strings"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
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
		headerUserId := c.GetHeader("X-User-Id")
		if headerUserId != "" && headerUserId != "guest" {
			var headerUser UserBindByHeader
			err := c.BindHeader(&headerUser)
			if err != nil || headerUser.Id == "" {
				c.JSON(http.StatusUnauthorized, gin.H{"message": "AUTH01-USER: You are not authorized to access this resource"})
				c.Abort()
				return
			}
			var user models.Users
			err = config.DB.Collection("Users").FindOne(context.TODO(), bson.D{{Key: "users_source_id", Value: headerUser.Id}}).Decode(&user)

			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"message": "AUTH01-USER: You are not authorized to access this resource"})
				c.Abort()
				return
			}
			c.Set("user", &user)
			c.Next()
		} else {

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
}

func CheckIfAuth(c *gin.Context) {
	headerUserId := c.GetHeader("X-User-Id")
	if headerUserId != "" && headerUserId == "guest" {
		return
	} else if headerUserId != "" && headerUserId != "guest" {
		var headerUser UserBindByHeader
		err := c.BindHeader(&headerUser)
		if err != nil || headerUser.Id == "" {
			return
		}
		var user models.Users
		err = config.DB.Collection("Users").FindOne(context.TODO(), bson.D{{Key: "users_source_id", Value: headerUser.Id}}).Decode(&user)

		if err != nil {
			return
		}
		c.Set("user", &user)
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
