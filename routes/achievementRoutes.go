package routes

import (
	"oosa_rewild/internal/middleware"
	"oosa_rewild/pkg/repository"

	"github.com/gin-gonic/gin"
)

func AchievementRoutes(r gin.IRouter) gin.IRouter {
	repo := repository.AchievementRepository{}

	main := r.Group("achievement", middleware.AuthMiddleware())
	{
		main.GET("", repo.Retrieve) // TODO: add deprecated header
		main.GET("/places", repo.Places)
	}

	return r
}
