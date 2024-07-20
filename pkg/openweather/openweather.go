package openweather

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"oosa_rewild/internal/config"
	"oosa_rewild/internal/helpers"

	"github.com/gin-gonic/gin"
)

type OpenWeatherRepository struct {
	ApiKey      string
	TokenSecret string
	SecretKey   string
	Method      string
	Url         string
	RequestUrl  string
	Args        map[string]string
}

type OWForecast struct {
	City OWForecastCity `json:"city"`
}

type OWForecastCity struct {
	Coord      OWForecastCityCoord `json:"coord"`
	Country    string              `json:"country"`
	Id         int                 `json:"id"`
	Name       string              `json:"name"`
	Population int                 `json:"population"`
	Sunrise    int                 `json:"sunrise"`
	Sunset     int                 `json:"sunset"`
	Timezone   int                 `json:"timezone"`
}

type OWForecastCityCoord struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

func (request OpenWeatherRepository) Forecast(c *gin.Context, lat float64, lng float64) OWForecast {
	endpoint := config.APP.OpenWeather + "forecast?lat=" + helpers.FloatToString(lat) + "&lon=" + helpers.FloatToString(lng) + "&appid=" + config.APP.OpenWeatherApiKey
	client := &http.Client{}
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		log.Fatal(err)
	}

	var responseRaw OWForecast
	// Create Response
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
		return responseRaw
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	json.Unmarshal(body, &responseRaw)

	return responseRaw
}
