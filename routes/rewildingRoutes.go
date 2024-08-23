package routes

import (
	"oosa_rewild/internal/middleware"
	"oosa_rewild/pkg/repository"

	"github.com/gin-gonic/gin"
)

func RewildingRoutes(r *gin.Engine) *gin.Engine {
	repoRewilding := repository.RewildingRepository{}
	repoRewildingPhoto := repository.RewildingPhotoRepository{}
	repoRewildingSearch := repository.RewildingSearchRepository{}

	rewilding := r.Group("/rewilding")
	{
		rewilding.GET("", repoRewilding.Retrieve)
		rewilding.POST("", middleware.AuthMiddleware(), repoRewilding.Create)
		// rewilding.GET("/references", middleware.AuthMiddleware(), repoRewilding.Options)
	}

	detail := rewilding.Group("/:id")
	{
		detail.GET("", repoRewilding.Read)
		detail.DELETE("", middleware.AuthMiddleware(), repoRewilding.Delete)
		detail.GET("/photos", repoRewildingPhoto.Retrieve)
		detail.GET("/photos/:photosId", repoRewildingPhoto.Read)
	}

	main := r.Group("/rewilding-search")
	{
		main.GET("", repoRewildingSearch.Retrieve)
		main.POST("", middleware.AuthMiddleware(), repoRewildingSearch.Create)
		main.GET("/:id", middleware.AuthMiddleware(), repoRewildingSearch.Read)
		main.PUT("/:id", middleware.AuthMiddleware(), repoRewildingSearch.Update)
	}

	return r
}
