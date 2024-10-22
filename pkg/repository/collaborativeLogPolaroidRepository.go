package repository

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"oosa_rewild/internal/config"
	"oosa_rewild/internal/helpers"
	"oosa_rewild/internal/models"
	"strconv"
	"time"

	exif "github.com/dsoprea/go-exif/v3"
	exifcommon "github.com/dsoprea/go-exif/v3/common"
	hemp "github.com/dsoprea/go-heic-exif-extractor/v2"
	jis "github.com/dsoprea/go-jpeg-image-structure/v2"
	pis "github.com/dsoprea/go-png-image-structure/v2"
	riimage "github.com/dsoprea/go-utility/v2/image"
	"github.com/gabriel-vasile/mimetype"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CollaborativeLogPolaroidRepository struct{}
type CollaborativeLogPolaroidRequest struct {
	EventPolaroidsMessage string `form:"event_polaroids_message"`
	EventPolaroidsTag     string `form:"event_polaroids_tag"`
}

func (r CollaborativeLogPolaroidRepository) Retrieve(c *gin.Context) {
	var Events models.Events
	err := CollaborativeLogRepository{}.ReadOne(c, &Events)
	if err != nil {
		return
	}

	var EventPolaroids []models.EventPolaroids
	agg := mongo.Pipeline{
		bson.D{{
			Key: "$match", Value: bson.M{
				"event_polaroids_event": Events.EventsId,
			},
		}},
		bson.D{{
			Key: "$lookup", Value: bson.M{
				"from":         "Users",
				"localField":   "event_polaroids_created_by",
				"foreignField": "_id",
				"as":           "event_polaroids_created_by_user",
			},
		}},
		bson.D{{
			Key: "$unwind", Value: "$event_polaroids_created_by_user",
		}},
	}
	cursor, err := config.DB.Collection("EventPolaroids").Aggregate(context.TODO(), agg)
	cursor.All(context.TODO(), &EventPolaroids)

	if err != nil {
		return
	}

	if len(EventPolaroids) == 0 {
		helpers.ResponseNoData(c, "No Data")
		return
	}
	c.JSON(http.StatusOK, EventPolaroids)
}

func (r CollaborativeLogPolaroidRepository) Create(c *gin.Context) {
	var payload CollaborativeLogPolaroidRequest
	validateError := helpers.ValidateForm(c, &payload)
	if validateError != nil {
		return
	}

	var Events models.Events
	err := CollaborativeLogRepository{}.ReadOne(c, &Events)
	if err != nil {
		return
	}

	countPolaroid := r.CountTotalPolaroids(Events.EventsId)
	fmt.Println("countPolaroid", countPolaroid)
	if countPolaroid >= config.APP_LIMIT.EventPolaroidLimit {
		helpers.ResponseBadRequestError(c, "Unable to add more polaroids. Maximum allowed: "+strconv.Itoa(int(config.APP_LIMIT.EventPolaroidLimit)))
		return
	}

	file, fileErr := c.FormFile("event_polaroids_file")
	if fileErr != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "No file is received",
		})
		return
	}

	uploadedFile, err := file.Open()

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Unable to open file",
		})
		return
	}

	b, _ := io.ReadAll(uploadedFile)
	mimeType := mimetype.Detect(b)

	var ecc riimage.MediaContext
	switch mimeType.String() {
	case "image/heic_":
		ecc, _ = hemp.NewHeicExifMediaParser().ParseBytes(b)
	case "image/jpeg":
		ecc, _ = jis.NewJpegMediaParser().ParseBytes(b)
	case "image/png":
		ecc, _ = pis.NewPngMediaParser().ParseBytes(b)
	default:
		c.JSON(http.StatusBadRequest, "Mime: "+mimeType.String()+" not supported")
		return
	}

	exif1, _, _ := ecc.Exif()

	//fmt.Println("DumpTree", exif1.DumpTree())
	//fmt.Println("DumpTags", exif1.DumpTags())
	//fmt.Println("Children", exif1.Children())
	//exif1.PrintTagTree(true)

	exifChildren := exif1.Children()

	latRef := ""
	lngRef := ""
	var latInterface []exifcommon.Rational
	var lngInterface []exifcommon.Rational
	lat := float64(0)
	lng := float64(0)

	for _, v := range exifChildren {
		entries := v.Entries()
		for _, v1 := range entries {
			v1Val, _ := v1.Value()

			switch v1.TagName() {
			case "GPSLatitudeRef":
				latRef = v1Val.(string)
			case "GPSLatitude":
				latInterface = v1Val.([]exifcommon.Rational)
			case "GPSLongitudeRef":
				lngRef = v1Val.(string)
			case "GPSLongitude":
				lngInterface = v1Val.([]exifcommon.Rational)
			}
		}
	}

	if len(latInterface) > 0 && latRef != "" {
		deg, _ := exif.NewGpsDegreesFromRationals(latRef, latInterface)
		lat = deg.Decimal()
	}
	if len(lngInterface) > 0 && lngRef != "" {
		deg, _ := exif.NewGpsDegreesFromRationals(lngRef, lngInterface)
		lng = deg.Decimal()
	}

	cloudflare := CloudflareRepository{}
	cloudflareResponse, postErr := cloudflare.Post(c, file)
	if postErr != nil {
		helpers.ResponseBadRequestError(c, postErr.Error())
		return
	}
	fileName := cloudflare.ImageDelivery(cloudflareResponse.Result.Id, "public")

	userDetail := helpers.GetAuthUser(c)
	insert := models.EventPolaroids{
		EventPolaroidsEvent:     Events.EventsId,
		EventPolaroidsUrl:       fileName,
		EventPolaroidsLat:       lat,
		EventPolaroidsLng:       lng,
		EventPolaroidsMessage:   payload.EventPolaroidsMessage,
		EventPolaroidsTag:       payload.EventPolaroidsTag,
		EventPolaroidsCreatedBy: userDetail.UsersId,
		EventPolaroidsCreatedAt: primitive.NewDateTimeFromTime(time.Now()),
	}

	radius := helpers.Haversine(lat, lng, Events.EventsLat, Events.EventsLng)

	eligibleAchievement := false
	if Events.EventsRewildingAchievementType != "" {
		if Events.EventsRewildingAchievementType == "big" || Events.EventsRewildingAchievementType == "small" {
			if radius <= 200 {
				eligibleAchievement = true
			}
		} else if Events.EventsRewildingAchievementType == "protect" {
			if radius <= 2000 {
				eligibleAchievement = true
			}
		}
	}

	insert.EventPolaroidsRadiusFromEvent = radius
	insert.EventPolaroidsAchievementEligible = &eligibleAchievement

	result, err := config.DB.Collection("EventPolaroids").InsertOne(context.TODO(), insert)
	if err != nil {
		fmt.Println("ERROR", err.Error())
		return
	}
	r.EventAchievementEligibility(c, Events)
	var EventPolaroids models.EventPolaroids
	config.DB.Collection("EventPolaroids").FindOne(context.TODO(), bson.D{{Key: "_id", Value: result.InsertedID}}).Decode(&EventPolaroids)
	c.JSON(http.StatusOK, EventPolaroids)
}

func (r CollaborativeLogPolaroidRepository) CountTotalPolaroids(eventId primitive.ObjectID) int64 {
	filter := bson.D{{Key: "event_polaroids_event", Value: eventId}}
	count, err := config.DB.Collection("EventPolaroids").CountDocuments(context.TODO(), filter)
	if err != nil {
		return 0
	}
	return count
}

func (r CollaborativeLogPolaroidRepository) EventAchievementEligibility(c *gin.Context, Event models.Events) {
	filter := bson.D{{Key: "event_polaroids_event", Value: Event.EventsId}, {Key: "event_polaroids_achievement_eligible", Value: true}}
	count, _ := config.DB.Collection("EventPolaroids").CountDocuments(context.TODO(), filter)
	if count > 0 {
		filterUpd := bson.D{{Key: "_id", Value: Event.EventsId}}
		eligible := true
		eventUpd := bson.D{{Key: "$set", Value: map[string]interface{}{
			"events_rewilding_achievement_eligible": &eligible,
		}}}
		config.DB.Collection("Events").UpdateOne(context.TODO(), filterUpd, eventUpd)
	}
}
