package helpers

import (
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func IndexStringInSlice(a string, list []string) int {
	for key, b := range list {
		if b == a {
			return key
		}
	}
	return -1
}

func MongoZeroID(a primitive.ObjectID) bool {
	zeroValue, _ := primitive.ObjectIDFromHex("000000000000000000000000")
	return a == zeroValue
}

func GetDomain(link string) (string, error) {
	var hostname string
	var temp []string
	url, err := url.Parse(link)
	if err != nil {
		return "", err
	}

	var urlstr string = url.String()

	if strings.HasPrefix(urlstr, "https") {
		hostname = strings.TrimPrefix(urlstr, "https://")
	} else if strings.HasPrefix(urlstr, "http") {
		hostname = strings.TrimPrefix(urlstr, "http://")
	} else {
		hostname = urlstr
	}

	if strings.HasPrefix(hostname, "www") {
		hostname = strings.TrimPrefix(hostname, "www.")
	}

	if strings.Contains(hostname, "/") {
		temp = strings.Split(hostname, "/")
		hostname = temp[0]
	}

	return hostname, err
}

func TimeIsBetween(t, min, max time.Time) bool {
	if min.After(max) {
		min, max = max, min
	}
	return (t.Equal(min) || t.After(min)) && (t.Equal(max) || t.Before(max))
}

func StringToDatetime(value string) time.Time {
	dt, _ := time.Parse("2006-01-02 15:04:05", value)
	return dt
}

func DataPaginate(c *gin.Context, length int64) []bson.D {
	var params []bson.D
	pageQuery := c.Query("page")

	if pageQuery == "" {
		pageQuery = "1"
	}

	page, err := strconv.ParseInt(pageQuery, 10, 64)

	if err != nil {
		page = 1
	}

	skip := (page - 1) * length

	return append(params, bson.D{{Key: "$skip", Value: skip}}, bson.D{{Key: "$limit", Value: length}})
}
