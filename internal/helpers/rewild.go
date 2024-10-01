package helpers

import (
	"context"
	"io"
	"net/http"
	"oosa_rewild/internal/config"
	"oosa_rewild/internal/models"

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

func RewildGooglePhotos(c *gin.Context, photos []*places.GoogleMapsPlacesV1Photo) []models.RewildingPhotos {
	var RewildingPhotos []models.RewildingPhotos
	for _, item := range photos {
		photo := GooglePlacePhoto(c, item.Name)
		RewildingPhotos = append(RewildingPhotos, models.RewildingPhotos{
			RewildingPhotosID:   primitive.NewObjectID(),
			RewildingPhotosPath: photo.PhotoUri,
		})
	}
	return RewildingPhotos
}
