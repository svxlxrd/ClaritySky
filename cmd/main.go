package main

import (
	"claritysky/internal/transport"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/weather/{city}", transport.CityWeatherHandler)

	log.Print("server started on port :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}

}
