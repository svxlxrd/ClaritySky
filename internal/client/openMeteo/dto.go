package openmeteo

type GeoResponse struct {
	Results []struct {
		Name      string  `json:"name"`
		Region    string  `json:"admin1"`
		Country   string  `json:"country"`
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	} `json:"results"`
}

type WeatherResponse struct {
	Current struct {
		Temperature float64 `json:"temperature_2m"`
		Time        string  `json:"time"`
	} `json:"current"`
}
