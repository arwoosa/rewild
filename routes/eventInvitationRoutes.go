package routes

import (
	"oosa_rewild/internal/middleware"
	"oosa_rewild/pkg/repository"

	"github.com/gin-gonic/gin"
)

func EventInvitationRoutes(r gin.IRouter) gin.IRouter {
	repo := repository.EventInvitationRepository{}

	main := r.Group("/event-invitations", middleware.AuthMiddleware())
	{
		main.GET("", repo.Retrieve)
	}

	detail := main.Group("/:id", middleware.AuthMiddleware())
	{
		detail.GET("", repo.Read)
		detail.PUT("", repo.Update)
	}

	return r
}
