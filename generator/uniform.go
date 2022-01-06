package generator

import (
	"constraints"
	"math/rand"
)

// An UniformInt generates a random sequence of discrete and uniformly
// distributed integers.
type UniformInt[T constraints.Integer] struct {
	rng       *rand.Rand
	min, rang T
}

// NewUniformInt returns an UniformInt using rng to generator of integers
// number generator of uniformly distributed integers.
func NewUniformtInt[T constraints.Integer](min, max T, rng *rand.Rand) *UniformInt[T] {
	if min > max {
		panic("min > max")
	}

	return &UniformInt[T]{min: min, rang: max - min, rng: rng}
}

// Next returns the next element in the sequence.
func (g *UniformInt[T]) Next() T {
	return T(g.min) + T(g.rng.Int63n(int64(g.rang)))
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
