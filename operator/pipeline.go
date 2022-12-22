package operator

import (
	"math/rand"

	"github.com/arl/evolve"
)

// A Pipeline is a compound evolutionary operator that applies multiple
// operators in sequence to a population.
type Pipeline[T any] []evolve.Operator[T]

// Apply applies each operator in the pipeline in sequence to the selection.
func (ops Pipeline[T]) Apply(sel []T, rng *rand.Rand) []T {
	for _, op := range ops {
		sel = op.Apply(sel, rng)
	}
	return sel
}
