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

func (r TestRepository) CreatePairs(c *gin.Context) {
	var pairs [][]string
	arr := []string{"1", "2", "3", "4", "5"}
	arrLeng := len(arr)

	for i := 0; i < arrLeng; i++ {
		for j := i + 1; j < arrLeng; j++ {
			pairString := []string{arr[i], arr[j]}
			pairs = append(pairs, pairString)
		}
	}

	c.JSON(200, pairs)
}

func (r TestRepository) EventFriend(c *gin.Context) {
	eventId, err := primitive.ObjectIDFromHex(c.Param("eventId"))
	if err != nil {
		helpers.ResponseBadRequestError(c, err.Error())
		return
	}

	EventRepository{}.HandleParticipantFriend(c, eventId)
}
