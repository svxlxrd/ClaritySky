package openmeteo

import (
	"claritysky/internal/domain"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

var client = &http.Client{
	Timeout: 5 * time.Second,
}

func getCoordinates(ctx context.Context, cityName string) (float64, float64, error) {
	url := fmt.Sprintf(
		"https://geocoding-api.open-meteo.com/v1/search?name=%s&limit=1",
		cityName,
	)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return 0, 0, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return 0, 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, 0, fmt.Errorf("geocoding API return status: %d", resp.StatusCode)
	}

	var city GeoResponse
	if err := json.NewDecoder(resp.Body).Decode(&city); err != nil {
		return 0, 0, err
	}

	if len(city.Results) == 0 {
		return 0, 0, fmt.Errorf("city not found")
	}

	return city.Results[0].Latitude, city.Results[0].Longitude, nil
}

func GetWeather(ctx context.Context, cityName string) (domain.Weather, error) {
	if cityName == "" {
		return domain.Weather{}, fmt.Errorf("city name is empty")
	}

	latitude, longitude, err := getCoordinates(ctx, cityName)
	if err != nil {
		return domain.Weather{}, err
	}

	url := fmt.Sprintf(
		"https://api.open-meteo.com/v1/forecast?latitude=%f&longitude=%f&current=temperature_2m,wind_speed_10m",
		latitude,
		longitude,
	)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return domain.Weather{}, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return domain.Weather{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return domain.Weather{}, fmt.Errorf("weather API return status: %d", resp.StatusCode)
	}

	var weather WeatherResponse
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
