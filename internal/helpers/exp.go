package helpers

import (
	"context"
	"oosa_rewild/internal/config"
	"oosa_rewild/internal/models"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	EXP_REWILDING  = "REWILDING"
	EXP_OOSA_DAILY = "OOSA_DAILY"
)

func ExpAward(c *gin.Context, expSource string, expAmount int, referenceId primitive.ObjectID) {
	userDetail := GetAuthUser(c)

	insert := models.Exp{
		ExpUser:      userDetail.UsersId,
		ExpPoints:    expAmount,
		ExpSource:    expSource,
		ExpCreatedAt: primitive.NewDateTimeFromTime(time.Now()),
	}

	if expSource == EXP_REWILDING {
		insert.ExpRewilding = referenceId
	}
	config.DB.Collection("Exp").InsertOne(context.TODO(), insert)
}
