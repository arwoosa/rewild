package repository

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"oosa_rewild/internal/config"
	"oosa_rewild/internal/helpers"
	"oosa_rewild/internal/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/mknote"
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
	userDetail := helpers.GetAuthUser(c)
	isCheck := false
	if c.Query("is_check") != "" {
		isCheck = true
	}

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

	var EventParticipantsCheck models.EventParticipants
	checkParticipant := config.DB.Collection("EventParticipants").FindOne(context.TODO(), bson.D{{Key: "event_participants_event", Value: Events.EventsId}, {Key: "event_participants_user", Value: userDetail.UsersId}}).Decode(&EventParticipantsCheck)

	if checkParticipant == mongo.ErrNoDocuments {
		helpers.ResponseBadRequestError(c, "You are not a participant of this event")
		return
	}

	if !isCheck {
		// countPolaroid := r.CountTotalPolaroids(Events.EventsId)
		countPolaroid := int64(EventParticipantsCheck.EventParticipantsPolaroidCount)
		if countPolaroid >= config.APP_LIMIT.EventPolaroidLimit {
			helpers.ResponseBadRequestError(c, "Unable to add more polaroids. Maximum allowed: "+strconv.Itoa(int(config.APP_LIMIT.EventPolaroidLimit)))
			return
		}

		match, errMessage := helpers.ValidateStringLength(payload.EventPolaroidsMessage, int(config.APP_LIMIT.LengthEventPolaroidMessage))
		if !match {
			helpers.ResponseBadRequestError(c, errMessage)
			return
		}
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

	lat := float64(0)
	lng := float64(0)

	b, _ := io.ReadAll(uploadedFile)
	reader := bytes.NewReader(b)

	exif.RegisterParsers(mknote.All...)
	x, err := exif.Decode(reader)
	if err != nil {
		//helpers.ResponseBadRequestError(c, "EXIF: "+err.Error())
	} else {
		lat, lng, _ = x.LatLong()
	}

	tm, _ := x.DateTime()

	fileName := ""
	if !isCheck {
		cloudflare := CloudflareRepository{}
		cloudflareResponse, postErr := cloudflare.Post(c, file)
		if postErr != nil {
			helpers.ResponseBadRequestError(c, postErr.Error())
			return
		}
		fileName = cloudflare.ImageDelivery(cloudflareResponse.Result.Id, "public")
	}

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

	radius := helpers.Haversine(lat, lng, Events.EventsLat, Events.EventsLng) * 1000

	eligibleAchievement := false
	isEventPeriod := false
	starType := 2

	if tm.Before(Events.EventsDateEnd.Time()) && tm.After(Events.EventsDate.Time()) {
		isEventPeriod = true
	}

	if Events.EventsRewildingAchievementType != "" {
		if radius <= config.APP_LIMIT.PolaroidAchievementRadius && isEventPeriod {
			starType = 1
		}
		eligibleAchievement = true
	} else {
		var Rewilding models.Rewilding
		filter := bson.D{{Key: "_id", Value: Events.EventsRewilding}}
		config.DB.Collection("Rewilding").FindOne(context.TODO(), filter).Decode(&Rewilding)
		//radius := helpers.Haversine(lat, lng, Rewilding.RewildingLat, Rewilding.RewildingLng) * 1000

		if radius <= config.APP_LIMIT.PolaroidAchievementRadius && isEventPeriod {
			starType = 1
		}
		eligibleAchievement = true
	}

	insert.EventPolaroidsIsEventPeriod = &isEventPeriod
	insert.EventPolaroidsRadiusFromEvent = &radius
	insert.EventPolaroidsAchievementEligible = &eligibleAchievement
	insert.EventPolaroidsStarType = starType

	if !isCheck {
		result, err := config.DB.Collection("EventPolaroids").InsertOne(context.TODO(), insert)
		if err != nil {
			fmt.Println("ERROR", err.Error())
			return
		}
		r.EventAchievementEligibility(c, Events)
		r.CountUploadPolaroidByParticipant(c, Events.EventsId, userDetail.UsersId)
		var EventPolaroids models.EventPolaroids
		config.DB.Collection("EventPolaroids").FindOne(context.TODO(), bson.D{{Key: "_id", Value: result.InsertedID}}).Decode(&EventPolaroids)

		helpers.BadgeAllocate(c, "P1", helpers.BADGE_EVENTS, Events.EventsId, primitive.NilObjectID)
		r.EventAchievementEligibility(c, Events)
		c.JSON(http.StatusOK, EventPolaroids)
	} else {
		c.JSON(http.StatusOK, insert)
	}
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
	var EventPolaroids []models.EventPolaroids
	var oneStarUsers []primitive.ObjectID
	filter := bson.D{{Key: "event_polaroids_event", Value: Event.EventsId}, {Key: "event_polaroids_achievement_eligible", Value: true}}
	count, _ := config.DB.Collection("EventPolaroids").CountDocuments(context.TODO(), filter)

	if count > 0 {
		filterUpd := bson.D{{Key: "_id", Value: Event.EventsId}}
		eligible := true
		eventUpd := bson.D{{Key: "$set", Value: map[string]interface{}{
			"events_rewilding_achievement_eligible": &eligible,
		}}}
		config.DB.Collection("Events").UpdateOne(context.TODO(), filterUpd, eventUpd)

		filterType2 := bson.D{
			{Key: "event_polaroids_event", Value: Event.EventsId},
			{Key: "event_polaroids_achievement_eligible", Value: true},
		}
		cursor, _ := config.DB.Collection("EventPolaroids").Find(context.TODO(), filterType2)
		cursor.All(context.TODO(), &EventPolaroids)

		if len(EventPolaroids) > 0 {
			unlockedDate := EventPolaroids[0].EventPolaroidsCreatedAt

			for _, v := range EventPolaroids {
				if v.EventPolaroidsStarType == 1 {
					oneStarUsers = append(oneStarUsers, v.EventPolaroidsCreatedBy)
				}
			}

			filterTwoStars := bson.D{
				{Key: "event_participants_event", Value: Event.EventsId},
			}
			if len(oneStarUsers) > 0 {
				filterOneStars := bson.D{
					{Key: "event_participants_event", Value: Event.EventsId},
					{Key: "event_participants_user", Value: bson.M{"$in": oneStarUsers}},
				}
				updOneStarType := bson.D{{Key: "$set", Value: map[string]interface{}{
					"event_participants_star_type":               1,
					"event_participants_achievement_eligible":    &eligible,
					"event_participants_achievement_unlocked_at": unlockedDate,
				}}}
				config.DB.Collection("EventParticipants").UpdateMany(context.TODO(), filterOneStars, updOneStarType)

				filterTwoStars = append(filterTwoStars, primitive.E{Key: "event_participants_user", Value: bson.M{"$nin": oneStarUsers}})
			}

			updTwoStarType := bson.D{{Key: "$set", Value: map[string]interface{}{
				"event_participants_star_type":               2,
				"event_participants_achievement_eligible":    &eligible,
				"event_participants_achievement_unlocked_at": unlockedDate,
			}}}
			config.DB.Collection("EventParticipants").UpdateMany(context.TODO(), filterTwoStars, updTwoStarType)
		}
	}
	r.HandleBadges(c, Event)
}

func (r CollaborativeLogPolaroidRepository) HandleBadges(c *gin.Context, Event models.Events) {
	agg := mongo.Pipeline{
		bson.D{{
			Key: "$match", Value: bson.M{
				"event_participants_event":                Event.EventsId,
				"event_participants_status":               GetEventParticipantStatus("ACCEPTED"),
				"event_participants_achievement_eligible": true,
			},
		}},
		/*bson.D{{
			Key: "$lookup", Value: bson.M{
				"from":         "UserBadges",
				"localField":   "event_participants_user",
				"foreignField": "user_badges_events_participant_user",
				"as":           "event_participant_badges",
			},
		}},
		*/
		bson.D{{Key: "$lookup", Value: bson.M{
			"as":   "event_participant_badges",
			"from": "UserBadges",
			"let": bson.M{
				"event_participants_event": "$event_participants_event",
				"event_participants_user":  "$event_participants_user",
			},
			"pipeline": []bson.M{
				/*{
					"$match": bson.M{
						"$expr": bson.M{
							"$and": []bson.M{
								{"eq": bson.A{"$user_badges_user", "$$event_participants_user"}},
								//{"eq": bson.A{"$user_badges_events", "$$event_participants_event"}},
							},
						},
					},
				},*/
				{"$match": bson.M{"user_badges_badge": helpers.StringToPrimitiveObjId("6695534c1da4d6fee70fc3ee")}},
				{
					"$match": bson.M{
						"$expr": bson.M{
							"$and": []bson.M{
								{"$eq": bson.A{"$user_badges_user", "$$event_participants_user"}},
								{"$eq": bson.A{"$user_badges_events", "$$event_participants_event"}},
							}},
					},
				},
			},
		}}},
		bson.D{{
			Key: "$unwind", Value: bson.M{
				"path":                       "$event_participant_badges",
				"preserveNullAndEmptyArrays": true,
			},
		}},
	}

	var results []models.EventParticipants
	cursor, err := config.DB.Collection("EventParticipants").Aggregate(context.TODO(), agg)
	if err != nil {
		panic(err)
	}
	cursor.All(context.TODO(), &results)

	for _, v := range results {
		if v.EventParticipantBadges == nil {
			helpers.BadgeAllocate(c, "R1", helpers.BADGE_EVENT_STARS, v.EventParticipantsEvent, v.EventParticipantsUser)
		}
	}

	c.JSON(200, results)
}

func (r CollaborativeLogPolaroidRepository) CountUploadPolaroidByParticipant(c *gin.Context, eventId primitive.ObjectID, userId primitive.ObjectID) {
	var EventParticipants models.EventParticipants
	filter := bson.D{
		{Key: "event_participants_event", Value: eventId},
		{Key: "event_participants_user", Value: userId},
	}
	config.DB.Collection("EventParticipants").FindOne(context.TODO(), filter).Decode(&EventParticipants)

	countFilter := bson.D{{Key: "event_polaroids_event", Value: eventId}, {Key: "event_polaroids_created_by", Value: userId}}
	count, _ := config.DB.Collection("EventPolaroids").CountDocuments(context.TODO(), countFilter)

	filterUpd := bson.D{{Key: "_id", Value: EventParticipants.EventParticipantsId}}
	eventParticipantUpd := bson.D{{Key: "$set", Value: map[string]interface{}{
		"event_participants_polaroid_count": count,
	}}}
	config.DB.Collection("EventParticipants").UpdateOne(context.TODO(), filterUpd, eventParticipantUpd)
	helpers.BadgeAllocate(c, "N4", helpers.BADGE_REWILDING, eventId, userId)
}
