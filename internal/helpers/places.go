package helpers

import (
	"log"
	"net/http"
	"oosa_rewild/internal/config"

	"github.com/gin-gonic/gin"
	"google.golang.org/api/option"
	"google.golang.org/api/places/v1"
	"googlemaps.github.io/maps"
)

func GooglePlacesInitialise(c *gin.Context) *places.Service {
	apiKey := config.APP.GoogleApiKey
	placesService, err := places.NewService(c, option.WithAPIKey(apiKey))
	if err != nil {
		log.Println("GooglePlacesInitialise-Error: ", err.Error())
		return nil
	}
	return placesService
}

func GooglePlaceById(c *gin.Context, id string) *places.GoogleMapsPlacesV1Place {
	placesService := GooglePlacesInitialise(c)
	if placesService != nil {
		placeReq := placesService.Places.Get("places/" + id).LanguageCode("zh-TW")
		placeReq.Header().Add("X-Goog-FieldMask", "id,types,displayName,formattedAddress,addressComponents,location,rating,userRatingCount,photos")
		places, errPlace := placeReq.Do()

		if errPlace != nil {
			log.Println("GooglePlaceById-Error: "+id, errPlace)
			c.JSON(http.StatusBadRequest, gin.H{"message": errPlace.Error()})
			return nil
		}
		return places
	}
	return nil
}

func GooglePlaceV1Search(c *gin.Context) {

}

func GoogleGeocoding(c *gin.Context) {
	//placesService := GooglePlacesInitialise(c)
}

func GooglePlacesGetArea(addresses []maps.AddressComponent, addressType string) (string, string) {
	longName := ""
	shortName := ""
	for _, val := range addresses {
		if StringInSlice(addressType, val.Types) {
			shortName = val.ShortName
			longName = val.LongName
		}
	}
	return longName, shortName
}

func GooglePlacesToLocationArray(addresses []maps.AddressComponent) []string {
	loc := []string{}
	key := []string{"administrative_area_level_2", "administrative_area_level_1", "country"}
	for _, v := range key {
		longName, _ := GooglePlacesGetArea(addresses, v)
		loc = append(loc, longName)
	}
	return loc
}

func GooglePlacesV1GetArea(addresses []*places.GoogleMapsPlacesV1PlaceAddressComponent, addressType string) (string, string) {
	longName := ""
	shortName := ""
	for _, val := range addresses {
		if StringInSlice(addressType, val.Types) {
			shortName = val.ShortText
			longName = val.LongText
		}
	}
	return longName, shortName
}

func GooglePlacesV1ToLocationArray(addresses []*places.GoogleMapsPlacesV1PlaceAddressComponent) []string {
	loc := []string{}
	key := []string{"administrative_area_level_1", "administrative_area_level_2", "country"}
	for _, v := range key {
		longName, _ := GooglePlacesV1GetArea(addresses, v)
		loc = append(loc, longName)
	}
	return loc
}

func GooglePlacePhoto(c *gin.Context, photoName string) *places.GoogleMapsPlacesV1PhotoMedia {
	placesService := GooglePlacesInitialise(c)
	if placesService != nil {
		url := photoName + "/media"
		placeReq := placesService.Places.Photos.GetMedia(url).SkipHttpRedirect(true).MaxHeightPx(400).MaxWidthPx(400)
		places, errPlace := placeReq.Do()

		if errPlace != nil {
			log.Println("GooglePlacePhoto-Error: "+photoName, errPlace)
			c.JSON(http.StatusBadRequest, gin.H{"message": errPlace.Error()})
			return nil
		}
		return places
	}
	return nil
}
