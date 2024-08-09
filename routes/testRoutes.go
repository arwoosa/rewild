package routes

import (
	"oosa_rewild/internal/middleware"
	"oosa_rewild/pkg/repository"

	"github.com/gin-gonic/gin"
)

func TestRoutes(r *gin.Engine) *gin.Engine {
	repo := repository.TestRepository{}

	main := r.Group("/test")
	{
		main.GET("/badge", middleware.AuthMiddleware(), repo.CreateBadge)
		main.GET("/notifications", middleware.AuthMiddleware(), repo.CreateNotification)
		main.GET("/exp", middleware.AuthMiddleware(), repo.CreateExp)
	}

	return r
}
