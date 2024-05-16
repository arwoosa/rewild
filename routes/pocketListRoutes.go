package routes

import (
	"oosa_rewild/internal/middleware"
	"oosa_rewild/pkg/repository"

	"github.com/gin-gonic/gin"
)

func PocketListRoutes(r *gin.Engine) *gin.Engine {
	repo := repository.PocketListRepository{}
	repoDetail := repository.PocketListItemsRepository{}

	main := r.Group("/pocket-list", middleware.AuthMiddleware())
	{
		main.GET("", repo.Retrieve)
		main.POST("", repo.Create)
	}

	detail := main.Group("/:id", middleware.AuthMiddleware())
	{
		detail.GET("", repo.Read)
		detail.PUT("", repo.Update)
		detail.DELETE("", repo.Delete)
	}

	detailItems := detail.Group("/items", middleware.AuthMiddleware())
	{
		detailItems.GET("", repoDetail.Retrieve)
		detailItems.POST("", repoDetail.Create)
		detailItems.GET("/:itemsId", repoDetail.Read)
		detailItems.PUT("/:itemsId", repoDetail.Update)
		detailItems.DELETE("/:itemsId", repoDetail.Delete)
	}

	return r
}
