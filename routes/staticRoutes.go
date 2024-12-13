package routes

import (
	"github.com/gin-gonic/gin"
)

func StaticRoutes(r *gin.Engine) *gin.Engine {
	r.Static("/event/cover", "./public/event/cover")
	r.Static("/badges", "./public/badges")
	return r
}
