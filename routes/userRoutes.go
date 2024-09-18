package routes

import (
	"oosa_rewild/internal/middleware"
	"oosa_rewild/pkg/repository"

	"github.com/gin-gonic/gin"
)

func UserRoutes(r *gin.Engine) *gin.Engine {
	//repoUser := repository.OosaUserRepository{}
	repoUserEvent := repository.OosaUserEventRepository{}
	repoUserAchievement := repository.OosaUserAchievementRepository{}

	me := r.Group("/user/:id", middleware.AuthMiddleware())
	{
		//me.GET("", repoUser.Read)
		me.GET("/event", repoUserEvent.Retrieve)
		me.GET("/achievement", repoUserAchievement.Retrieve)
	}

	return r
}
