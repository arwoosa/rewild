package helpers

import (
	"context"
	"fmt"
	"oosa_rewild/internal/config"
	"oosa_rewild/internal/models"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	BADGE_REWILDING = 1
	BADGE_FRIENDS   = 2
	BADGE_SOCIAL    = 3
)

func BadgeAllocate(c *gin.Context, badgeCode string, badgeSource int, badgeReference primitive.ObjectID) {
	badgeDetail := BadgeDetail(badgeCode)
	userDetail := GetAuthUser(c)

	var UserBadges models.UserBadges

	if badgeDetail.BadgesIsOnce {
		filter := bson.D{
			{Key: "user_badges_user", Value: userDetail.UsersId},
			{Key: "user_badges_badge", Value: badgeDetail.BadgesId},
		}
		config.DB.Collection("UserBadges").FindOne(context.TODO(), filter).Decode(&UserBadges)

		if !MongoZeroID(UserBadges.UserBadgesId) {
			return
		}
	}

	insert := models.UserBadges{
		UserBadgesUser:      userDetail.UsersId,
		UserBadgesBadge:     badgeDetail.BadgesId,
		UserBadgesCreatedAt: primitive.NewDateTimeFromTime(time.Now()),
	}

	if badgeSource == BADGE_REWILDING {
		insert.UserBadgesRewilding = badgeReference
	}

	result, err := config.DB.Collection("UserBadges").InsertOne(context.TODO(), insert)

	NotificationMessage := models.NotificationMessage{
		Message: "太棒了！恭喜你獲得了一枚新的徽章!",
		Data:    []map[string]interface{}{NotificationFormatBadges(badgeDetail)},
	}
	NotificationsCreate(c, NOTIFICATION_BADGE_NEW, userDetail.UsersId, NotificationMessage, result.InsertedID.(primitive.ObjectID))
	if err != nil {
		fmt.Println("ERROR", err.Error())
		return
	}
}

func BadgeDetail(badgeCode string) models.Badges {
	var results models.Badges
	filter := bson.D{{Key: "badges_code", Value: badgeCode}}
	config.DB.Collection("Badges").FindOne(context.TODO(), filter).Decode(&results)
	return results
}

func BadgeEvents(c *gin.Context, eventId primitive.ObjectID) {
	BadgeAllocate(c, "N1", BADGE_REWILDING, eventId)
}
