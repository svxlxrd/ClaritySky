package retry

import (
	"math"
	"math/rand/v2"
	"time"
)

func CalculateDelay(attempt int, baseDelay time.Duration, maxDelay time.Duration) time.Duration {
	exp := math.Pow(2, float64(attempt))
	delay := time.Duration(float64(baseDelay) * exp)

	jitterRange := float64(delay) * 0.2
	jitter := (rand.Float64() * 2 - 1) * jitterRange

	finalDelay := delay + time.Duration(jitter)

	if finalDelay > maxDelay {
		finalDelay = maxDelay
	}

	return finalDelay
}
