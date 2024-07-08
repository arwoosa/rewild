package helpers

import (
	"fmt"
	"math"
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

func Haversine(lat1 float64, lon1 float64, lat2 float64, lon2 float64) float64 {

	fmt.Println(lat1, lon1, lat2, lon2)
	var R float64 = 6371
	var x1 = lat2 - lat1
	dLat := toRad(x1)
	var x2 = lon2 - lon1
	dLon := toRad(x2)

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(toRad(lat1))*math.Cos(toRad(lat2))*
			math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	d := R * c
	return d
}

func toRad(deg float64) float64 {
	return deg * (math.Pi / 180)
}
