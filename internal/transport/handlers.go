package transport

import (
	"github.com/svxlxrd/weather-rest-api/internal/domain"
	"context"
	"encoding/json"
	"net/http"
	"time"
)

type WeatherService interface {
	GetWeather(ctx context.Context, cityName string) ([]domain.Weather, error)
}

type WeatherHandler struct {
	weatherProvider WeatherService
}

func NewWeatherHandler(ws WeatherService) *WeatherHandler {
	return &WeatherHandler{
		weatherProvider: ws,
	}
}

func (h *WeatherHandler) CityWeatherHandler(w http.ResponseWriter, r *http.Request) {
	city := r.PathValue("city")

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	weather, err := h.weatherProvider.GetWeather(ctx, city)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(weather); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
