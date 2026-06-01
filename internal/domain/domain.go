package domain

import "time"

type Weather struct {
	City 		string 
	Region 		string
	Country 	string
	Temperature float64
	UpdatedAt   time.Time
}
