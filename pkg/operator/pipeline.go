package operator

import (
	"math/rand"

	"github.com/aurelien-rainone/evolve/pkg/api"
)

// A Pipeline is a compound evolutionary operator that applies multiple
// operators, in sequence, to a starting population.
type Pipeline []api.Operator

// Apply applies each operation in the pipeline in turn to the selection.
func (ops Pipeline) Apply(
	sel []interface{},
	rng *rand.Rand) []interface{} {

	for _, op := range ops {
		sel = op.Apply(sel, rng)
	}
	return sel
}
