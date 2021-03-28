package generator

import (
	"math/rand"
)

// An UniformInt generates a random sequence of discrete and uniformly
// distributed integers.
type UniformInt struct {
	rng       *rand.Rand
	min, rang int64
}

// NewUniformInt returns an UniformInt using rng to generator of integers
// number generator of uniformly distributed integers.
func NewUniformtInt(min, max int64, rng *rand.Rand) *UniformInt {
	if min > max {
		panic("min > max")
	}

	return &UniformInt{min: min, rang: max - min, rng: rng}
}

// Next returns the next element in the sequence.
func (g *UniformInt) Next() int64 {
	return g.min + g.rng.Int63n(g.rang)
}

// An UniformFloat generates a random sequence of continuous and uniformly
// distributed floating point numbers.
type UniformFloat struct {
	rng       *rand.Rand
	min, rang float64
}

// NewUniformFloat returns a generator of uniformly distributed floating point
// numbers.
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
