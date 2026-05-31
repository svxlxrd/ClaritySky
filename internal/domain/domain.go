package domain

import "time"

type Weather struct {
	Temperature float64
	WindSpeed   float64
	UpdatedAt   time.Time
}
