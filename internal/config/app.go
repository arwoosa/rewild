package config

import (
	"os"
	"path/filepath"
	"strconv"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	GoogleOauthClientId        string
	GoogleApiKey               string
	AppPort                    string
	DbApiHost                  string
	DbApiPort                  string
	DbApiDatabase              string
	DbApiUsername              string
	DbApiPassword              string
	FlickrApiKey               string
	FlickrSecret               string
	FlickrUserName             string
	CloudflareImageAuthToken   string
	ClourdlareImageAccountId   string
	ClourdlareImageAccountHash string
	ClourdlareImageDeliveryUrl string
	OpenWeather                string
	OpenWeatherApiKey          string
}

type AppLimit struct {
	PocketList         int
	EventPolaroidLimit int64
}

var APP AppConfig
var APP_LIMIT AppLimit

func InitialiseConfig() {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	errEnv := godotenv.Load(filepath.Join(dir, ".env"))
	if errEnv != nil {
		godotenv.Load()
	}

	APP.GoogleOauthClientId = os.Getenv("GOOGLE_OAUTH_CLIENT_ID")
	APP.GoogleApiKey = os.Getenv("GOOGLE_API_KEY")
	APP.AppPort = os.Getenv("APP_PORT")
	APP.DbApiHost = os.Getenv("DB_API_HOST")
	APP.DbApiPort = os.Getenv("DB_API_PORT")
	APP.DbApiDatabase = os.Getenv("DB_API_DATABASE")
	APP.DbApiUsername = os.Getenv("DB_API_USERNAME")
	APP.DbApiPassword = os.Getenv("DB_API_PASSWORD")
	APP.FlickrApiKey = os.Getenv("FLICKR_API_KEY")
	APP.FlickrSecret = os.Getenv("FLICKR_SECRET_KEY")
	APP.FlickrUserName = os.Getenv("FLICKR_UPLOAD_USERNAME")
	APP.CloudflareImageAuthToken = os.Getenv("CLOUDFLARE_IMAGE_AUTH_TOKEN")
	APP.ClourdlareImageAccountId = os.Getenv("CLOURDLARE_IMAGE_ACCOUNT_ID")
	APP.ClourdlareImageAccountHash = os.Getenv("CLOURDLARE_IMAGE_ACCOUNT_HASH")
	APP.ClourdlareImageDeliveryUrl = os.Getenv("CLOURDLARE_IMAGE_DELIVERY_URL")
	APP.OpenWeather = os.Getenv("OPENWEATHER_API_BASE_URL")
	APP.OpenWeatherApiKey = os.Getenv("OPENWEATHER_API_KEY")

	polaroidLimit, err := strconv.ParseInt(os.Getenv("EVENT_POLAROID_LIMIT"), 10, 64)
	if err == nil {
		APP_LIMIT.EventPolaroidLimit = polaroidLimit
	} else {
		APP_LIMIT.EventPolaroidLimit = 0
	}

	APP_LIMIT.PocketList = 100
}
