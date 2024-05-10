package routes

import (
	"oosa_rewild/internal/middleware"
	"oosa_rewild/pkg/repository"

	"github.com/gin-gonic/gin"
)

func CollaborativeLogRoutes(r *gin.Engine) *gin.Engine {
	repo := repository.CollaborativeLogRepository{}
	repoAlbumLink := repository.CollaborativeLogAlbumLinkRepository{}
	repoPolaroid := repository.CollaborativeLogPolaroidRepository{}

	main := r.Group("/collaborative-log", middleware.AuthMiddleware())
	{
		main.GET("", repo.Retrieve)
	}

	detail := main.Group("/:id", middleware.AuthMiddleware())
	{
	}

	albumLink := detail.Group("/album-link", middleware.AuthMiddleware())
	{
		albumLink.GET("", repoAlbumLink.Retrieve)
		albumLink.POST("", repoAlbumLink.Create)
		// albumLink.GET("/:messageBoardId", repoAlbumLink.Read)
		// albumLink.PUT("/:messageBoardId", repoAlbumLink.Update)
		// albumLink.DELETE("/:messageBoardId", repoAlbumLink.Delete)
	}

	polaroid := detail.Group("/polaroid", middleware.AuthMiddleware())
	{
		polaroid.GET("", repoPolaroid.Retrieve)
		polaroid.POST("", repoPolaroid.Create)
		// albumLink.GET("/:messageBoardId", repoAlbumLink.Read)
		// albumLink.PUT("/:messageBoardId", repoAlbumLink.Update)
		// albumLink.DELETE("/:messageBoardId", repoAlbumLink.Delete)
	}

	/*detail := main.Group("/:id", middleware.AuthMiddleware())
	{
		detail.GET("", repo.Read)
		detail.PUT("", repo.Update)
	}*/

	return r
}
