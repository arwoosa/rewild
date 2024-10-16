package routes

import (
	"github.com/gin-gonic/gin"
)

func StaticRoutes(r *gin.Engine) *gin.Engine {
	r.Static("/event/cover", "./public/event/cover")
	return r
}
