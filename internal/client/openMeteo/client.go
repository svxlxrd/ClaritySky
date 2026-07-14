package openmeteo

import (
	"github.com/svxlxrd/weather-rest-api/internal/domain"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type HTTPClient interface {
	Do(ctx context.Context, req *http.Request) (*http.Response, error)
}

type Provider struct {
	httpClient HTTPClient
}

func NewProvider(hc HTTPClient) *Provider {
	return &Provider{
		httpClient: hc,
	}
}

func (p *Provider) getCoordinates(
	ctx context.Context,
	cityName string,
) (GeoResponse, error) {
	url := fmt.Sprintf(
		"https://geocoding-api.open-meteo.com/v1/search?name=%s&count=5",
		cityName,
	)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return GeoResponse{}, err
	}

	resp, err := p.httpClient.Do(ctx, req)
	if err != nil {
		if resp != nil {
			resp.Body.Close()
		}
		return GeoResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return GeoResponse{}, fmt.Errorf("geocoding API return status: %d", resp.StatusCode)
	}

	var cities GeoResponse
	if err := json.NewDecoder(resp.Body).Decode(&cities); err != nil {
		return GeoResponse{}, err
	}

	if len(cities.Results) == 0 {
		return GeoResponse{}, fmt.Errorf("city not found")
	}

	return cities, nil
}

func (p *Provider) GetWeather(
	ctx context.Context,
	cityName string,
) ([]domain.Weather, error) {
	if cityName == "" {
		return nil, fmt.Errorf("city name is empty")
	}

	cities, err := p.getCoordinates(ctx, cityName)
	if err != nil {
		return nil, err
	}

	var results []domain.Weather

	for _, city := range cities.Results {
		url := fmt.Sprintf(
			"https://api.open-meteo.com/v1/forecast?latitude=%f&longitude=%f&current=temperature_2m",
			city.Latitude,
			city.Longitude,
		)

		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			continue
		}

		resp, err := p.httpClient.Do(ctx, req)
		if err != nil {
			if resp != nil {
				resp.Body.Close()
			}
			continue
		}

		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			continue
		}

		var weather WeatherResponse
		if err := json.NewDecoder(resp.Body).Decode(&weather); err != nil {
			resp.Body.Close()
			continue
		}
		resp.Body.Close()

		updatedAt, err := time.Parse("2006-01-02T15:04", weather.Current.Time)
		if err != nil {
			continue
		}

		results = append(results, domain.Weather{
			City:        city.Name,
			Region:      city.Region,
			Country:     city.Country,
			Temperature: weather.Current.Temperature,
			UpdatedAt:   updatedAt,
		})
	}

	return results, nil
}
