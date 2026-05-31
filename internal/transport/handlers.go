package transport

import (
	"claritysky/internal/client/openMeteo"
	"encoding/json"
	"net/http"
)

func CityWeatherHandler(w http.ResponseWriter, r *http.Request) {
	city := r.PathValue("city")

	weather, err := openMeteo.GetWeather(city)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(weather); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
