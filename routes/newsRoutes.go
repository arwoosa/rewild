package routes

import (
	"oosa_rewild/internal/middleware"
	"oosa_rewild/pkg/repository"

	"github.com/gin-gonic/gin"
)

func NewsRoutes(r *gin.Engine) *gin.Engine {
	repo := repository.NewsRepository{}

	main := r.Group("/news")
	{
		main.GET("", repo.Retrieve)
		main.POST("", repo.Create, middleware.AuthMiddleware())
	}

	detail := main.Group("/:id", middleware.AuthMiddleware())
	{
		detail.GET("", repo.Read)
		detail.PUT("", repo.Update)
		detail.DELETE("", repo.Delete)
	}

	return r
}
