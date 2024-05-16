package routes

import (
	"oosa_rewild/pkg/repository"

	"github.com/gin-gonic/gin"
)

func CloudfareRoutes(r *gin.Engine) *gin.Engine {
	repo := repository.CloudflareRepository{}

	main := r.Group("/cloudflare")
	{
		main.GET("", repo.Retrieve)
		main.GET(":imageId", repo.Read)
		main.POST("", repo.Upload)
	}

	return r
}
