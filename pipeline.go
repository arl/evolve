package evolve

import (
	"math/rand"
)

// A Pipeline is a compound evolutionary operator that applies multiple
// operators in sequence to a population.
type Pipeline[T any] []Operator[T]

// Apply applies each operator in the pipeline in sequence to the selection.
func (ops Pipeline[T]) Apply(pop *Population[T], rng *rand.Rand) {
	for _, op := range ops {
		op.Apply(pop, rng)
	}
}
