package routes

import (
	_ "oosa_rewild/docs"
	"oosa_rewild/internal/middleware"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
)

func RegisterRoutes() *gin.Engine {
	r := gin.Default()

	checkSsoUserGroup := r.Group("", middleware.CheckRegisterMiddleware())
	RewildingRoutes(checkSsoUserGroup)
	RewildingRegisterRoutes(checkSsoUserGroup)
	PocketListRoutes(checkSsoUserGroup)
	EventUserRoutes(checkSsoUserGroup)
	EventRoutes(checkSsoUserGroup)
	AchievementRoutes(checkSsoUserGroup)
	EventInvitationRoutes(checkSsoUserGroup)
	CollaborativeLogRoutes(checkSsoUserGroup)
	FlickrRoutes(checkSsoUserGroup)
	CloudfareRoutes(checkSsoUserGroup)
	LinkRoutes(checkSsoUserGroup)
	ReferenceRoutes(checkSsoUserGroup)
	TestRoutes(checkSsoUserGroup)
	UserRoutes(checkSsoUserGroup)
	NewsRoutes(checkSsoUserGroup)
	StaticRoutes(checkSsoUserGroup)

	healthRoutes(r)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return r
}
