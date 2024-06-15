package helpers

import (
	"net/url"
	"strings"

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
