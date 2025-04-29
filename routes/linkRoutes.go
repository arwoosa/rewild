package routes

import (
	"oosa_rewild/pkg/repository"

	"github.com/gin-gonic/gin"
)

func LinkRoutes(r gin.IRouter) gin.IRouter {
	repo := repository.LinkRepository{}

	main := r.Group("/link-query")
	{
		main.POST("", repo.Query)
	}

	return r
}
