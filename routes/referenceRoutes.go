package routes

import (
	"oosa_rewild/pkg/repository"

	"github.com/gin-gonic/gin"
)

func ReferenceRoutes(r gin.IRouter) gin.IRouter {
	repo := repository.ReferenceRepository{}

	main := r.Group("/references")
	{
		main.GET("", repo.Options)
	}

	return r
}
