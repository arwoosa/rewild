package repository

import (
	"oosa_rewild/internal/helpers"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TestRepository struct{}

func (r TestRepository) CreateBadge(c *gin.Context) {
	// helpers.BadgeAllocate(c, "R2")
}

func (r TestRepository) CreateNotification(c *gin.Context) {
	userDetail := helpers.GetAuthUser(c)
	notifCode := "Invitation"
	notifType := "New"
	helpers.NotificationsCreate(c, notifCode, notifType, userDetail.UsersId, "TEST", helpers.StringToPrimitiveObjId("66005b3ef4ca36269a55b468"))
}

func (r TestRepository) CreateExp(c *gin.Context) {
	helpers.ExpAward(c, helpers.EXP_REWILDING, 1, primitive.NilObjectID)
}
