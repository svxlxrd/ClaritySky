package openmeteo

type GeoResponse struct {
	Results []struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	} `json:"results"`
}

type WeatherResponse struct {
	Current struct {
		Temperature float64 `json:"temperature_2m"`
		WindSpeed   float64 `json:"wind_speed_10m"`
		Time        string  `json:"time"`
	} `json:"current"`
}
