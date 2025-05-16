package routes

import (
	"oosa_rewild/internal/middleware"
	"oosa_rewild/pkg/repository"

	"github.com/gin-gonic/gin"
)

func UserRoutes(r gin.IRouter) gin.IRouter {
	//repoUser := repository.OosaUserRepository{}
	repoUserEvent := repository.OosaUserEventRepository{}
	repoUserAchievement := repository.OosaUserAchievementRepository{}
	repoCollaborativeLog := repository.CollaborativeLogRepository{}

	me := r.Group("/user/:id")
	{
		//me.GET("", repoUser.Read)
		me.GET("/events", middleware.AuthMiddleware(), repoUserEvent.Retrieve)
		me.GET("/achievement", repoUserAchievement.Retrieve)
		me.GET("/collaborative-log", middleware.AuthMiddleware(), repoCollaborativeLog.Retrieve)
	}

	return r
}
