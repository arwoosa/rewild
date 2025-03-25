package routes

import (
	"oosa_rewild/internal/middleware"
	"oosa_rewild/pkg/repository"

	"github.com/gin-gonic/gin"
)

func EventUserRoutes(r gin.IRouter) gin.IRouter {
	repo := repository.UserEventRepository{}

	main := r.Group("/my/event", middleware.AuthMiddleware())
	{
		main.GET("", repo.Retrieve)
	}

	return r
}
