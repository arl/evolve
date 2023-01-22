package operator

import (
	"math/rand"

	"github.com/arl/evolve"
)

// A Pipeline is a compound evolutionary operator that applies multiple
// operators in sequence to a population.
//
// TODO(arl): move in another package? evolve maybe?
type Pipeline[T any] []evolve.Operator[T]

// Apply applies each operator in the pipeline in sequence to the selection.
func (ops Pipeline[T]) Apply(sel *evolve.Population[T], rng *rand.Rand) {
	for _, op := range ops {
		op.Apply(sel, rng)
	}
}
