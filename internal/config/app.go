package config

import (
	"os"
	"path/filepath"
	"strconv"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	BaseUrl                    string
	GoogleOauthClientId        string
	GoogleApiKey               string
	AppPort                    string
	DbConnection               string
	DbApiDatabase              string
	FlickrApiKey               string
	FlickrSecret               string
	FlickrUserName             string
	CloudflareImageAuthToken   string
	ClourdlareImageAccountId   string
	ClourdlareImageAccountHash string
	ClourdlareImageDeliveryUrl string
	OpenWeather                string
	OpenWeatherApiKey          string
	NotificationContextKey     string
	NotificationHeaderName     string
}

type AppLimit struct {
	PocketList                     int64
	PocketListItems                int64
	EventPolaroidLimit             int64
	EventAccountingLimit           int64
	EventAnnouncementLimit         int64
	EventMessageBoardLimit         int64
	LengthPocketListName           int64
	LengthRewildingName            int64
	LengthRewildingImage           int64
	LengthRewildingReferenceLink   int64
	LengthEventName                int64
	LengthEventMessageBoardMessage int64
	LengthEventAccountingMessage   int64
	LengthEventInvitationMessage   int64
	LengthEventPolaroidMessage     int64
	LengthEventParticipantMessage  int64
	MinimumTopRanking              int64
	PolaroidAchievementRadius      float64
}

var APP AppConfig
var APP_LIMIT AppLimit

func InitialiseConfig() {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	errEnv := godotenv.Load(filepath.Join(dir, ".env"))
	if errEnv != nil {
		godotenv.Load()
	}

	APP.BaseUrl = os.Getenv("APP_BASE_URL")
	APP.GoogleOauthClientId = os.Getenv("GOOGLE_OAUTH_CLIENT_ID")
	APP.GoogleApiKey = os.Getenv("GOOGLE_API_KEY")
	APP.AppPort = os.Getenv("APP_PORT")
	APP.DbConnection = os.Getenv("DB_CONNECTION")
	APP.DbApiDatabase = os.Getenv("DB_API_DATABASE")
	APP.FlickrApiKey = os.Getenv("FLICKR_API_KEY")
	APP.FlickrSecret = os.Getenv("FLICKR_SECRET_KEY")
	APP.FlickrUserName = os.Getenv("FLICKR_UPLOAD_USERNAME")
	APP.CloudflareImageAuthToken = os.Getenv("CLOUDFLARE_IMAGE_AUTH_TOKEN")
	APP.ClourdlareImageAccountId = os.Getenv("CLOURDLARE_IMAGE_ACCOUNT_ID")
	APP.ClourdlareImageAccountHash = os.Getenv("CLOURDLARE_IMAGE_ACCOUNT_HASH")
	APP.ClourdlareImageDeliveryUrl = os.Getenv("CLOURDLARE_IMAGE_DELIVERY_URL")
	APP.OpenWeather = os.Getenv("OPENWEATHER_API_BASE_URL")
	APP.OpenWeatherApiKey = os.Getenv("OPENWEATHER_API_KEY")
	APP.NotificationContextKey = os.Getenv("NOTIFICATION_CONTEXT_KEY")
	APP.NotificationHeaderName = os.Getenv("NOTIFICATION_HEADER_NAME")

	APP_LIMIT.EventPolaroidLimit = 0
	APP_LIMIT.EventAccountingLimit = 0
	APP_LIMIT.EventAnnouncementLimit = 0
	APP_LIMIT.EventMessageBoardLimit = 0
	APP_LIMIT.PocketList = 0
	APP_LIMIT.PocketListItems = 0
	APP_LIMIT.LengthPocketListName = 0
	APP_LIMIT.LengthRewildingName = 0
	APP_LIMIT.LengthRewildingImage = 0
	APP_LIMIT.LengthRewildingReferenceLink = 0
	APP_LIMIT.LengthEventName = 0
	APP_LIMIT.LengthEventMessageBoardMessage = 0
	APP_LIMIT.LengthEventAccountingMessage = 0
	APP_LIMIT.LengthEventInvitationMessage = 0
	APP_LIMIT.LengthEventPolaroidMessage = 0
	APP_LIMIT.LengthEventParticipantMessage = 0
	APP_LIMIT.MinimumTopRanking = 0
	APP_LIMIT.PolaroidAchievementRadius = 0

	polaroidLimit, err := strconv.ParseInt(os.Getenv("EVENT_POLAROID_LIMIT"), 10, 64)
	eventAccountingLimit, eventAccountingLimitErr := strconv.ParseInt(os.Getenv("EVENT_ACCOUNTING_LIMIT"), 10, 64)
	eventAnnouncementLimit, eventAnnouncementLimitErr := strconv.ParseInt(os.Getenv("EVENT_ANNOUNCEMENT_LIMIT"), 10, 64)
	eventMessageBoardLimit, eventMessageBoardLimitErr := strconv.ParseInt(os.Getenv("EVENT_MESSAGE_BOARD_LIMIT"), 10, 64)
	pocketListLimit, pocketlistLimitErr := strconv.ParseInt(os.Getenv("POCKET_LIST_LIMIT"), 10, 64)
	pocketListitemsLimit, pocketlistitemsLimitErr := strconv.ParseInt(os.Getenv("POCKET_LIST_ITEMS_LIMIT"), 10, 64)
	lengthPocketListName, lengthPocketListNameErr := strconv.ParseInt(os.Getenv("LENGTH_POCKET_LIST_NAME"), 10, 64)
	lengthRewildingName, lengthRewildingNameErr := strconv.ParseInt(os.Getenv("LENGTH_REWILDING_NAME"), 10, 64)
	lengthRewildingImage, lengthRewildingImageErr := strconv.ParseInt(os.Getenv("LENGTH_REWILDING_IMAGE"), 10, 64)
	lengthRewildingReferenceLink, lengthRewildingReferenceLinkErr := strconv.ParseInt(os.Getenv("LENGTH_REWILDING_REFERENCE_LINK"), 10, 64)
	lengthEventName, lengthEventNameErr := strconv.ParseInt(os.Getenv("LENGTH_EVENT_NAME"), 10, 64)
	lengthEventMessageBoardMessage, lengthEventMessageBoardMessageErr := strconv.ParseInt(os.Getenv("LENGTH_EVENT_MESSAGE_BOARD_MESSAGE"), 10, 64)
	lengthEventAccountingMessage, lengthEventAccountingMessageErr := strconv.ParseInt(os.Getenv("LENGTH_EVENT_ACCOUNTING_MESSAGE"), 10, 64)
	lengthEventInvitationMessage, lengthEventInvitationMessageErr := strconv.ParseInt(os.Getenv("LENGTH_EVENT_INVITATION_MESSAGE"), 10, 64)
	lengthEventPolaroidMessage, lengthEventPolaroidMessageErr := strconv.ParseInt(os.Getenv("LENGTH_EVENT_POLAROID_MESSAGE"), 10, 64)
	lengthEventParticipantMessage, lengthEventParticipantMessageErr := strconv.ParseInt(os.Getenv("LENGTH_EVENT_PARTICIPANT_MESSAGE"), 10, 64)
	minimumTopRanking, minimumTopRankingErr := strconv.ParseInt(os.Getenv("MINIMUM_TOP_RANKING"), 10, 64)
	polaroidAchievementRadius, polaroidAchievementRadiusErr := strconv.ParseFloat(os.Getenv("POLAROID_ACHIEVEMENT_RADIUS"), 64)

	if err == nil {
		APP_LIMIT.EventPolaroidLimit = polaroidLimit
	}
	if eventAccountingLimitErr == nil {
		APP_LIMIT.EventAccountingLimit = eventAccountingLimit
	}
	if eventAnnouncementLimitErr == nil {
		APP_LIMIT.EventAnnouncementLimit = eventAnnouncementLimit
	}
	if eventMessageBoardLimitErr == nil {
		APP_LIMIT.EventMessageBoardLimit = eventMessageBoardLimit
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
	if lengthRewildingNameErr == nil {
		APP_LIMIT.LengthRewildingName = lengthRewildingName
	}
	if lengthRewildingImageErr == nil {
		APP_LIMIT.LengthRewildingImage = lengthRewildingImage
	}
	if lengthRewildingReferenceLinkErr == nil {
		APP_LIMIT.LengthRewildingReferenceLink = lengthRewildingReferenceLink
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
	if lengthEventInvitationMessageErr == nil {
		APP_LIMIT.LengthEventInvitationMessage = lengthEventInvitationMessage
	}
	if lengthEventPolaroidMessageErr == nil {
		APP_LIMIT.LengthEventPolaroidMessage = lengthEventPolaroidMessage
	}
	if lengthEventParticipantMessageErr == nil {
		APP_LIMIT.LengthEventParticipantMessage = lengthEventParticipantMessage
	}
	if polaroidAchievementRadiusErr == nil {
		APP_LIMIT.PolaroidAchievementRadius = polaroidAchievementRadius
	}
}
