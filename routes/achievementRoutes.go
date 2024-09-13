package routes

import (
	"oosa_rewild/internal/middleware"
	"oosa_rewild/pkg/repository"

	"github.com/gin-gonic/gin"
)

func AchievementRoutes(r *gin.Engine) *gin.Engine {
	repo := repository.AchievementRepository{}

	main := r.Group("/achievement", middleware.AuthMiddleware())
	{
		main.GET("", repo.Retrieve)
	}

	return r
}
