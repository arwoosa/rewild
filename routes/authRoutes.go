package routes

import (
	"oosa_rewild/internal/middleware"
	"oosa_rewild/pkg/repository"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(r *gin.Engine) *gin.Engine {
	repo := repository.AuthRepository{}

	main := r.Group("/auth")
	{
		main.GET("", middleware.AuthMiddleware(), repo.Auth)
		main.GET("test-badge", middleware.AuthMiddleware(), repo.TestBadge)
	}

	return r
}
