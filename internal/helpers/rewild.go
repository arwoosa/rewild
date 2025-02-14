package helpers

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"oosa_rewild/internal/config"
	"oosa_rewild/internal/models"
	"sort"
	"strings"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/api/places/v1"
)

func GetRewildByPlaceId(placeId string) (models.Rewilding, bool) {
	var Rewilding models.Rewilding
	err := config.DB.Collection("Rewilding").FindOne(context.TODO(), bson.D{{Key: "rewilding_place_id", Value: placeId}}).Decode(&Rewilding)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return Rewilding, false
		}
	}
	return Rewilding, true
}

func RewildSaveGooglePhotos(c *gin.Context, photos []*places.GoogleMapsPlacesV1Photo) ([]models.RewildingPhotos, error) {
	var RewildingPhotos []models.RewildingPhotos
	for _, item := range photos {
		photo := GooglePlacePhoto(c, item.Name)

		resp, err := http.Get(photo.PhotoUri)
		if err != nil {
			//return "", fmt.Errorf("GET error: %v", err)
			return RewildingPhotos, err
		}
		defer resp.Body.Close()

		data, _ := io.ReadAll(resp.Body)

		RwPhoto := models.RewildingPhotos{
			RewildingPhotosID:   primitive.NewObjectID(),
			RewildingPhotosData: data,
			// RewildingPhotosPath:   item.Name,
		}
		RewildingPhotos = append(RewildingPhotos, RwPhoto)
	}
	return RewildingPhotos, nil
}

const googlePhotoUrlTpl = "https://maps.googleapis.com/maps/api/place/photo?photo_reference=%s&maxwidth=400&key=%s"

func RewildGooglePhotos(c *gin.Context, photos []*places.GoogleMapsPlacesV1Photo) []models.RewildingPhotos {
	var RewildingPhotos []models.RewildingPhotos
	for _, item := range photos {
		slice := strings.Split(item.Name, "/")
		if len(slice) != 4 {
			continue
		}
		// fmt.Println(slice[3])
		// fmt.Printf(googlePhotoUrlTpl+"\n", slice[3], config.APP.GoogleApiTestKey)
		// photo := GooglePlacePhoto(c, item.Name)
		RewildingPhotos = append(RewildingPhotos, models.RewildingPhotos{
			RewildingPhotosID:   primitive.NewObjectID(),
			RewildingPhotosPath: fmt.Sprintf(googlePhotoUrlTpl, slice[3], config.APP.GoogleApiKey),
		})
	}
	return RewildingPhotos
}

func RewildingAchievementByLatLng(c *gin.Context, lat float64, lng float64) (models.RefAchievementPlaces, error) {
	var RefAchievementPlaces []models.RefAchievementPlaces
	cursor, _ := config.DB.Collection("RefAchievementPlaces").Find(context.TODO(), bson.D{})
	cursor.All(context.TODO(), &RefAchievementPlaces)

	for k, v := range RefAchievementPlaces {
		dist := Haversine(lat, lng, v.RefAchievementPlacesLat, v.RefAchievementPlacesLng)
		RefAchievementPlaces[k].RefAchievementPlacesDistance = dist
	}

	sort.Slice(RefAchievementPlaces, func(i, j int) bool {
		return RefAchievementPlaces[i].RefAchievementPlacesDistance < RefAchievementPlaces[j].RefAchievementPlacesDistance
	})

	if RefAchievementPlaces[0].RefAchievementPlacesDistance < 200 {
		return RefAchievementPlaces[0], nil
	}

	return models.RefAchievementPlaces{}, errors.New("no achievement for this location")
}
