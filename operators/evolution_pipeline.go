package operators

import (
	"errors"
	"math/rand"

	"github.com/aurelien-rainone/evolve/framework"
)

// An EvolutionPipeline is a compound evolutionary operator that applies
// multiple operators (of the same Candidate type) in series.
//
// By combining EvolutionPipeline operators with SplitEvolution operators,
// elaborate evolutionary schemes can be constructed.
type EvolutionPipeline struct {
	pipeline []framework.EvolutionaryOperator
}

// NewEvolutionPipeline creates a pipeline consisting of the specified operators
// in the order that they are supplied.
func NewEvolutionPipeline(pipeline ...framework.EvolutionaryOperator) (*EvolutionPipeline, error) {
	ep := &EvolutionPipeline{pipeline: pipeline}
	if len(ep.pipeline) == 0 {
		return nil, errors.New("pipeline must contain at least one operator")
	}
	return ep, nil
}

// Apply applies each operation in the pipeline in turn to the selection.
func (ep *EvolutionPipeline) Apply(
	selectedCandidates []framework.Candidate,
	rng *rand.Rand) []framework.Candidate {

	population := selectedCandidates
	for _, op := range ep.pipeline {
		population = op.Apply(population, rng)
	}
	return population
}
