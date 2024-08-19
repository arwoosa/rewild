package repository

import (
	"oosa_rewild/internal/helpers"
	"oosa_rewild/internal/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TestRepository struct{}

func (r TestRepository) CreateBadge(c *gin.Context) {
	helpers.BadgeAllocate(c, "M5", 0, primitive.NilObjectID)
}

func (r TestRepository) CreateNotification(c *gin.Context) {
	userDetail := helpers.GetAuthUser(c)
	notifCode := helpers.NOTIFICATION_INVITATION
	NotificationMessage := models.NotificationMessage{
		Message: "TEST",
		Data:    []map[string]interface{}{},
	}
	helpers.NotificationsCreate(c, notifCode, userDetail.UsersId, NotificationMessage, helpers.StringToPrimitiveObjId("66005b3ef4ca36269a55b468"))
}

func (r TestRepository) CreateExp(c *gin.Context) {
	helpers.ExpAward(c, helpers.EXP_REWILDING, 1, primitive.NilObjectID)
}
