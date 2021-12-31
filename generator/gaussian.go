package generator

import (
	"math/rand"
)

// Gaussian generates normally-distributed values where mean and standard
// deviation are determined by a Float64Generator.
type Gaussian struct {
	rng          *rand.Rand
	mean, stddev Float
}

// TODO: doc
func NewGaussian(mean, stddev Float, rng *rand.Rand) *Gaussian {
	return &Gaussian{mean: mean, stddev: stddev, rng: rng}
}

// Next returns the next normally-distributed value.
func (g *Gaussian) Next() float64 {
	return g.rng.NormFloat64()*float64(g.stddev.Next()) + float64(g.mean.Next())
}
