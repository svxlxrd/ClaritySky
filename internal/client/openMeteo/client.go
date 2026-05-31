package openMeteo

import (
	"claritysky/internal/domain"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type geoResponse struct {
	Results []struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	} `json:"results"`
}

type weatherResponse struct {
	Current struct {
		Temperature float64 `json:"temperature_2m"`
		WindSpeed   float64 `json:"wind_speed_10m"`
		Time        string  `json:"time"`
	} `json:"current"`
}

func getCoordinates(cityName string) (float64, float64, error) {
	url := fmt.Sprintf(
		"https://geocoding-api.open-meteo.com/v1/search?name=%s&country_code=RU&limit=1",
		cityName,
	)

	resp, err := http.Get(url)
	if err != nil {
		return 0, 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, 0, err
	}

	var city geoResponse
	if err := json.NewDecoder(resp.Body).Decode(&city); err != nil {
		return 0, 0, err
	}

	if len(city.Results) == 0 {
		return 0, 0, fmt.Errorf("city not found")
	}

	return city.Results[0].Latitude, city.Results[0].Longitude, nil
}

func GetWeather(cityName string) (domain.Weather, error) {
	latitude, longitude, err := getCoordinates(cityName)
	if err != nil {
		return domain.Weather{}, err
	}

	url := fmt.Sprintf(
		"https://api.open-meteo.com/v1/forecast?latitude=%f&longitude=%f&current=temperature_2m,wind_speed_10m",
		latitude,
		longitude,
	)

	resp, err := http.Get(url)
	if err != nil {
		return domain.Weather{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return domain.Weather{}, err
	}

	var weather weatherResponse
	if err := json.NewDecoder(resp.Body).Decode(&weather); err != nil {
		return domain.Weather{}, err
	}

	updatedAt, err := time.Parse("2006-01-02T15:04", weather.Current.Time)
	if err != nil {
		return domain.Weather{}, err
	}

	return domain.Weather{
		Temperature: weather.Current.Temperature,
		WindSpeed:   weather.Current.WindSpeed,
		UpdatedAt:   updatedAt,
	}, nil
}
