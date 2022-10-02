package operator

import (
	"math/rand"

	"github.com/arl/evolve"
)

// A Switch is a compound evolutionary operator that switches, from one
// generation to the next, the evolutionary operator to apply.
type Switch[T any] struct {
	ops []evolve.Operator[T]
	cur int
}

// NewSwitch returns a new Switch operator.
func NewSwitch[T any](ops ...evolve.Operator[T]) *Switch[T] {
	return &Switch[T]{ops: ops}
}

// Switch applies a single one of the operators to the selection. At the next
// generation the next operator will be applied, etc.
func (s *Switch[T]) Apply(sel []T, rng *rand.Rand) []T {
	sel = s.ops[s.cur].Apply(sel, rng)
	s.cur++
	if s.cur == len(s.ops) {
		s.cur = 0
	}
	return sel
}
