package generator

import (
	"math"
	"math/rand"
)

// Exponential generates a random sequence following an exponential distribution.
type Exponential struct {
	rate Generator[float64]
	rng  *rand.Rand
}

// NewExponential creates a generator of exponentially-distributed values from a
// distribution with a rate controlled by the rate generator parameter.
//
// The mean of this distribution is 1/rate and its variance is 1/rate².
// Note: the rate generator must only return strictly positive values.
func NewExponential(rate Generator[float64], rng *rand.Rand) *Exponential {
	return &Exponential{rate: rate, rng: rng}
}

// Next returns the next exponentially-distributed value.
func (g *Exponential) Next() float64 {
	var u float64
	for {
		// Get a uniformly-distributed random double in [0 1)
		u = g.rng.Float64()
		if u == 0 {
			// Reject zero
			continue
		}
		break
	}
	return -math.Log(u) / g.rate.Next()
}
