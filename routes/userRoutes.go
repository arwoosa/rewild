package routes

import (
	"github.com/gin-gonic/gin"
)

func UserRoutes(r *gin.Engine) *gin.Engine {
	/*oosa := repository.OosaUserEventRepository{}

	me := r.Group("/user/:id", middleware.AuthMiddleware())
	{
		me.GET("/events", repoUserEvents.Retrieve)
		me.GET("/badges", repoUserEvents.Retrieve)
	}*/

	return r
}
