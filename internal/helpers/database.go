package helpers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func FloatToDecimal128(float float64) primitive.Decimal128 {
	formatted, _ := primitive.ParseDecimal128(fmt.Sprint(float))
	return formatted
}

func StringToPrimitiveObjId(value string) primitive.ObjectID {
	id, _ := primitive.ObjectIDFromHex(value)
	return id
}

func StringToPrimitiveDateTime(value string) primitive.DateTime {
	time := StringToDateTime(value)
	return primitive.NewDateTimeFromTime(time)
}

func StringToDateTime(value string) time.Time {
	date, _ := time.Parse("2006-01-02T15:04:05Z07:00", value)
	return date
}

func ResultEmpty(c *gin.Context, err error) {
	if err == mongo.ErrNoDocuments {
		ResponseNoData(c, err.Error())
		return
	}
}

func ResultMessageSuccess(c *gin.Context, message string) {
	c.JSON(http.StatusOK, gin.H{"message": message})
}

func ResultMessageError(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, gin.H{"message": message})
}
