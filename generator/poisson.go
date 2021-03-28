package generator

import (
	"math"
	"math/rand"
)

// Poisson generates Poisson-distributed values.
type Poisson struct {
	rng  *rand.Rand
	mean Float
}

func NewPoisson(mean Float, rng *rand.Rand) *Poisson {
	return &Poisson{mean: mean, rng: rng}
}

// Next returns the next generated Poisson-distributed value.
func (p *Poisson) Next() int64 {
	var (
		x int64
		t float64
	)
	for {
		t -= math.Log(p.rng.Float64()) / p.mean.Next()
		if t > 1.0 {
			break
		}
		x++
	}

	return x
}
