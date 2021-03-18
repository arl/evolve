package generator

import (
	"math/rand"
)

// An Int generates sequences of integers.
type Int interface {
	// Next returns the next number in the sequence.
	Next() int64
}

type ConstInt int64

// Next always returns i.
func (i ConstInt) Next() int64 {
	return int64(i)
}

// An UniformInt generates a random sequence of discrete and uniformly
// distributed integers.
type UniformInt struct {
	rng       *rand.Rand
	min, rang int64
}

// NewUniformInt returns an UniformInt using rng to generator of integers  number generator of uniformly distributed integers.
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
