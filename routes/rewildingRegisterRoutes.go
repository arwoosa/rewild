package routes

import (
	"oosa_rewild/internal/middleware"
	"oosa_rewild/pkg/repository"

	"github.com/gin-gonic/gin"
)

func RewildingRegisterRoutes(r *gin.Engine) *gin.Engine {
	repoRegisterRewilding := repository.RewildingRegisterRepository{}
	repoRegisterRewildingPhoto := repository.RewildingRegisterPhotoRepository{}

	rewildingRegister := r.Group("/rewilding-register", middleware.AuthMiddleware())
	{
		rewildingRegister.GET("", repoRegisterRewilding.Retrieve)
		rewildingRegister.POST("", repoRegisterRewilding.Create)
		rewildingRegister.GET("/references", repoRegisterRewilding.Options)
	}

	rewildingRegisterDetail := rewildingRegister.Group("/:id")
	{
		rewildingRegisterDetail.GET("", repoRegisterRewilding.Read)
		rewildingRegisterDetail.PUT("", repoRegisterRewilding.Update)
	}

	rewildingRegisterPhoto := rewildingRegisterDetail.Group("/photos", middleware.AuthMiddleware())
	{
		rewildingRegisterPhoto.GET("", repoRegisterRewildingPhoto.Retrieve)
		rewildingRegisterPhoto.POST("", repoRegisterRewildingPhoto.Create)
		rewildingRegisterPhoto.GET(":photosId", repoRegisterRewildingPhoto.Read)
		rewildingRegisterPhoto.DELETE(":photosId", repoRegisterRewildingPhoto.Delete)
	}

	return r
}
