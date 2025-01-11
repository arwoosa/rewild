package routes

import (
	_ "oosa_rewild/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
)

func RegisterRoutes() *gin.Engine {
	r := gin.Default()
	RewildingRoutes(r)
	RewildingRegisterRoutes(r)
	PocketListRoutes(r)
	EventUserRoutes(r)
	EventRoutes(r)
	AchievementRoutes(r)
	EventInvitationRoutes(r)
	CollaborativeLogRoutes(r)
	FlickrRoutes(r)
	CloudfareRoutes(r)
	LinkRoutes(r)
	ReferenceRoutes(r)
	TestRoutes(r)
	UserRoutes(r)
	NewsRoutes(r)
	StaticRoutes(r)
	healthRoutes(r)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return r
}
