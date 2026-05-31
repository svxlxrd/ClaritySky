package transport

import (
	"claritysky/internal/client/openmeteo"
	"context"
	"encoding/json"
	"net/http"
	"time"
)

func CityWeatherHandler(w http.ResponseWriter, r *http.Request) {
	city := r.PathValue("city")

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	weather, err := openmeteo.GetWeather(ctx, city)
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
