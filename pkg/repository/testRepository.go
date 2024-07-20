package repository

import (
	"oosa_rewild/internal/helpers"

	"github.com/gin-gonic/gin"
)

type TestRepository struct{}

func (r TestRepository) CreateBadge(c *gin.Context) {
	helpers.BadgeAllocate(c, "R2")
}

func (r TestRepository) CreateNotification(c *gin.Context) {
	userDetail := helpers.GetAuthUser(c)
	notifCode := "Invitation"
	notifType := "New"
	helpers.NotificationsCreate(c, notifCode, notifType, userDetail.UsersId, "TEST", helpers.StringToPrimitiveObjId("66005b3ef4ca36269a55b468"))
}
