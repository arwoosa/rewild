package repository

import (
	"net/http"
	"oosa_rewild/internal/helpers"
	"oosa_rewild/internal/models"

	"github.com/gin-gonic/gin"
)

type OosaUserAchievementRepository struct{}

func (r OosaUserAchievementRepository) Retrieve(c *gin.Context) {
	var results []models.AchievementRewildingV2
	userId := helpers.StringToPrimitiveObjId(c.Param("id")) // 被查詢者id
	authUserId := helpers.GetAuthUser(c).UsersId            // 登入者id

	var userStatus string
	if authUserId.IsZero() {
		userStatus = "stranger" // 未登入
	} else if authUserId == userId {
		userStatus = "owner" // 擁有者
	} else {
		userStatus = "others" // 查看者(有登入)
	}

	err := AchievementRepository{}.GetAchievementsByUserIdV2(c, userId, &results)
	if err != nil {
		helpers.ResponseError(c, err.Error())
		return
	}

	resp := models.AchievementEvent{
		UserStatus:   userStatus,
		Achievements: results,
	}

	c.JSON(http.StatusOK, resp)
}
