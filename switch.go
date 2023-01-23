package evolve

import (
	"math/rand"
)

// A Switch holds multiple evolutionary operators and switches from one to the
// next each time it is called.
type Switch[T any] struct {
	ops []Operator[T]
	cur int
}

// NewSwitch returns a new Switch operator.
func NewSwitch[T any](ops ...Operator[T]) *Switch[T] {
	return &Switch[T]{ops: ops}
}

// Switch applies a single one of the operators to the selection. At the next
// generation the next operator will be applied, etc.
func (s *Switch[T]) Apply(pop *Population[T], rng *rand.Rand) {
	s.ops[s.cur].Apply(pop, rng)
	s.cur++
	if s.cur == len(s.ops) {
		s.cur = 0
	}
}
