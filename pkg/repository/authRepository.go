package repository

import (
	"github.com/gin-gonic/gin"
)

type AuthRepository struct{}

func (t AuthRepository) Auth(c *gin.Context) {
	user, _ := c.Get("user")
	c.JSON(200, user)
}

func (t AuthRepository) TestBadge(c *gin.Context) {
	user, _ := c.Get("user")
	c.JSON(200, user)
}
