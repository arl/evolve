package generator

import (
	"constraints"
	"math"
	"math/rand"
)

// Poisson generates Poisson-distributed values.
type Poisson[U constraints.Unsigned] struct {
	rng  *rand.Rand
	mean Float
}

func NewPoisson[U constraints.Unsigned](mean Float, rng *rand.Rand) *Poisson[U] {
	return &Poisson[U]{mean: mean, rng: rng}
}

// Next returns the next generated Poisson-distributed value.
func (p *Poisson[U]) Next() U {
	var (
		x U
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
