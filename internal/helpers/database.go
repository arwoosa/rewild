package helpers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func FloatToDecimal128(float float64) primitive.Decimal128 {
	formatted, _ := primitive.ParseDecimal128(fmt.Sprint(float))
	return formatted
}

func FloatToString(float float64) string {
	return strconv.FormatFloat(float, 'f', -1, 64)
}

func StringToFloat(stringFloat string) float64 {
	val, _ := strconv.ParseFloat(stringFloat, 64)
	return val
}

func StringToInt(stringInt string) int {
	val, _ := strconv.Atoi(stringInt)
	return val
}

func Decimal128ToFloat(decimal128 primitive.Decimal128) float64 {
	decimalString := decimal128.String()
	float, _ := strconv.ParseFloat(decimalString, 64)
	return float
}

func StringToPrimitiveObjId(value string) primitive.ObjectID {
	id, _ := primitive.ObjectIDFromHex(value)
	return id
}

func StringToPrimitiveDateTime(value string) primitive.DateTime {
	time := StringToDateTime(value)
	return primitive.NewDateTimeFromTime(time)
}

func StringDateToPrimitiveDateTime(value string) primitive.DateTime {
	time := StringDateToDateTime(value)
	return primitive.NewDateTimeFromTime(time)
}

func StringDateToDateTime(value string) time.Time {
	date, _ := time.Parse("2006-01-02", value)
	return date
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

func ResultNotFound(c *gin.Context, err error, message string) {
	if err == mongo.ErrNoDocuments {
		if message == "" {
			message = err.Error()
		}
		ResponseNotFound(c, message)
		return
	}
}

func ResultMessageSuccess(c *gin.Context, message string) {
	c.JSON(http.StatusOK, gin.H{"message": message})
}

func ResultMessageError(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, gin.H{"message": message})
}
