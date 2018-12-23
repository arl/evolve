package operator

import (
	"math/rand"

	"github.com/arl/evolve"
)

// A Pipeline is a compound evolutionary operator that applies multiple
// operators in sequence to a population.
type Pipeline []evolve.Operator

// Apply applies each operator in the pipeline in sequence to the selection.
func (ops Pipeline) Apply(sel []interface{}, rng *rand.Rand) []interface{} {
	for _, op := range ops {
		sel = op.Apply(sel, rng)
	}
	return sel
}
