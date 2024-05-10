package helpers

import (
	"fmt"
	"oosa_rewild/internal/config"

	"github.com/gin-gonic/gin"
	"googlemaps.github.io/maps"
)

func GoogleMapsInitialise() *maps.Client {
	apiKey := config.APP.GoogleApiKey
	mapService, err := maps.NewClient(maps.WithAPIKey(apiKey))
	if err != nil {
		fmt.Println("ERROR...", err.Error())
	}
	return mapService
}

func GoogleMapsGeocode(c *gin.Context, lat float64, lng float64) maps.GeocodingResult {
	mapService := GoogleMapsInitialise()
	req := &maps.GeocodingRequest{
		LatLng: &maps.LatLng{
			Lat: lat,
			Lng: lng,
		},
	}
	result, err := mapService.Geocode(c, req)
	if err != nil {
		fmt.Println("ERROR...", err.Error())
	}
	return result[0]
}

func GoogleMapsElevation(c *gin.Context, lat float64, lng float64) maps.ElevationResult {
	mapService := GoogleMapsInitialise()
	var latLngData []maps.LatLng

	latLngData = append(latLngData, maps.LatLng{Lat: lat, Lng: lng})

	req := &maps.ElevationRequest{
		Locations: latLngData,
	}
	result, err := mapService.Elevation(c, req)
	if err != nil {
		fmt.Println("ERROR...", err.Error())
	}
	return result[0]
}
