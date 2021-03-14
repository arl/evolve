package generator

import (
	"math/rand"
)

// An FloatGenerator generates sequences of floating point numbers.
type FloatGenerator interface {
	// Next returns the next number in the sequence.
	Next() float64
}

type ConstFloat64 float64

// Next always returns i.
func (f ConstFloat64) Next() float64 {
	return float64(f)
}

// An UniformFloat generates a random sequence of continuous and uniformly
// distributed floating point numbers..
type UniformFloat struct {
	rng       *rand.Rand
	min, rang float64
}

// NewUniformFloat returns a generator of uniformly distributed floating point numbers.
func NewUniformFloat(min, max float64, rng *rand.Rand) *UniformFloat {
	if min > max {
		panic("min > max")
	}

	return &UniformFloat{min: min, rang: max - min, rng: rng}
}

// Next returns the next element in the sequence.
func (g *UniformFloat) Next() float64 {
	return g.min + g.rng.Float64()*g.rang
}
