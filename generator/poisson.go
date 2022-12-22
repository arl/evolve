package generator

import (
	"math"
	"math/rand"

	"golang.org/x/exp/constraints"
)

// Poisson generates Poisson-distributed values, that are always positive.
type Poisson[I constraints.Integer] struct {
	rng  *rand.Rand
	mean Float
}

func NewPoisson[I constraints.Integer](mean Float, rng *rand.Rand) *Poisson[I] {
	return &Poisson[I]{mean: mean, rng: rng}
}

// Next returns the next generated Poisson-distributed value.
func (p *Poisson[I]) Next() I {
	var (
		i I
		t float64
	)
	for {
		t -= math.Log(p.rng.Float64()) / p.mean.Next()
		if t > 1.0 {
			break
		}
		i++
	}

	return i
}
