package routes

import (
	"oosa_rewild/pkg/repository"

	"github.com/gin-gonic/gin"
)

func FlickrRoutes(r *gin.Engine) *gin.Engine {
	repo := repository.FlickrRepository{}

	main := r.Group("/flickr")
	{
		main.GET("", repo.Retrieve)
		main.GET(":id", repo.Read)
		main.POST("", repo.Upload)
		main.GET("oauth", repo.Oauth)
		main.GET("oauth/callback", repo.OauthCb)
	}

	return r
}
