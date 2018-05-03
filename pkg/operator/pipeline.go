package operator

import (
	"math/rand"

	"github.com/aurelien-rainone/evolve/framework"
)

// A Pipeline is a compound evolutionary operator that applies multiple
// operators, in sequence, to a starting population.
type Pipeline []framework.EvolutionaryOperator

// Apply applies each operation in the pipeline in turn to the selection.
func (ops Pipeline) Apply(
	sel []framework.Candidate,
	rng *rand.Rand) []framework.Candidate {

	for _, op := range ops {
		sel = op.Apply(sel, rng)
	}
	return sel
}
