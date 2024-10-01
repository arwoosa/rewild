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
	PocketList                     int64
	PocketListItems                int64
	EventPolaroidLimit             int64
	LengthPocketListName           int64
	LengthEventName                int64
	LengthEventMessageBoardMessage int64
	LengthEventAccountingMessage   int64
	MinimumTopRanking              int64
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

	APP_LIMIT.EventPolaroidLimit = 0
	APP_LIMIT.PocketList = 0
	APP_LIMIT.PocketListItems = 0
	APP_LIMIT.LengthPocketListName = 0
	APP_LIMIT.LengthEventName = 0
	APP_LIMIT.LengthEventMessageBoardMessage = 0
	APP_LIMIT.LengthEventAccountingMessage = 0
	APP_LIMIT.MinimumTopRanking = 0

	polaroidLimit, err := strconv.ParseInt(os.Getenv("EVENT_POLAROID_LIMIT"), 10, 64)
	pocketListLimit, pocketlistLimitErr := strconv.ParseInt(os.Getenv("POCKET_LIST_LIMIT"), 10, 64)
	pocketListitemsLimit, pocketlistitemsLimitErr := strconv.ParseInt(os.Getenv("POCKET_LIST_ITEMS_LIMIT"), 10, 64)
	lengthPocketListName, lengthPocketListNameErr := strconv.ParseInt(os.Getenv("LENGTH_POCKET_LIST_NAME"), 10, 64)
	lengthEventName, lengthEventNameErr := strconv.ParseInt(os.Getenv("LENGTH_EVENT_NAME"), 10, 64)
	lengthEventMessageBoardMessage, lengthEventMessageBoardMessageErr := strconv.ParseInt(os.Getenv("LENGTH_EVENT_MESSAGE_BOARD_MESSAGE"), 10, 64)
	lengthEventAccountingMessage, lengthEventAccountingMessageErr := strconv.ParseInt(os.Getenv("LENGTH_EVENT_ACCOUNTING_MESSAGE"), 10, 64)
	minimumTopRanking, minimumTopRankingErr := strconv.ParseInt(os.Getenv("MINIMUM_TOP_RANKING"), 10, 64)
	if err == nil {
		APP_LIMIT.EventPolaroidLimit = polaroidLimit
	}
	if pocketlistLimitErr == nil {
		APP_LIMIT.PocketList = pocketListLimit
	}
	if pocketlistitemsLimitErr == nil {
		APP_LIMIT.PocketListItems = pocketListitemsLimit
	}
	if lengthPocketListNameErr == nil {
		APP_LIMIT.LengthPocketListName = lengthPocketListName
	}
	if lengthEventNameErr == nil {
		APP_LIMIT.LengthEventName = lengthEventName
	}
	if minimumTopRankingErr == nil {
		APP_LIMIT.MinimumTopRanking = minimumTopRanking
	}
	if lengthEventMessageBoardMessageErr == nil {
		APP_LIMIT.LengthEventMessageBoardMessage = lengthEventMessageBoardMessage
	}
	if lengthEventAccountingMessageErr == nil {
		APP_LIMIT.LengthEventAccountingMessage = lengthEventAccountingMessage
	}
}
