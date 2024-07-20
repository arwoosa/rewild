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

func BadgeAllocate(c *gin.Context, badgeCode string) {
	badgeDetail := BadgeDetail(badgeCode)
	userDetail := GetAuthUser(c)

	var UserBadges models.UserBadges
	if badgeDetail.BadgesIsOnce {
		filter := bson.D{
			{Key: "user_badges_user", Value: userDetail.UsersId},
			{Key: "user_badges_badge", Value: badgeDetail.BadgesId},
		}
		config.DB.Collection("UserBadges").FindOne(context.TODO(), filter).Decode(&UserBadges)

		if !MongoZeroID(userDetail.UsersId) {
			fmt.Println("This badge is only received once")
			return
		}
	}

	insert := models.UserBadges{
		UserBadgesUser:      userDetail.UsersId,
		UserBadgesBadge:     badgeDetail.BadgesId,
		UserBadgesCreatedAt: primitive.NewDateTimeFromTime(time.Now()),
	}

	_, err := config.DB.Collection("UserBadges").InsertOne(context.TODO(), insert)
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
