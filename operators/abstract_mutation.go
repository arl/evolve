package operators

import (
	"fmt"
	"math/rand"

	"github.com/aurelien-rainone/evolve/framework"
	"github.com/aurelien-rainone/evolve/number"
)

// Mutater is the interface implemented by objects defining the Mutate function.
type Mutater interface {

	// Mutate performs mutation on a candidate.
	//
	// The original candidate is let untouched while the mutant is returned.
	Mutate(framework.Candidate, *rand.Rand) framework.Candidate
}

// AbstractMutation is a generic struct for mutation implementations.
//
// It supports all mutation processes that operate on an unique candidate.
// The mutation probability is configurable, its effect depends on the specific
// mutation implementation, where it will be documented.
type AbstractMutation struct {
	mutationProbability number.ProbabilityGenerator
	Mutater
}

// NewAbstractMutation creates an AbstractMutation configured with the
// provided options.
//
// TODO: example of use of how setting options
func NewAbstractMutation(mutater Mutater, options ...Option) (*AbstractMutation, error) {
	// create with default options, a mutation probability of zero (the default
	// may be changed by the specific mutation implementation)
	op := &AbstractMutation{
		mutationProbability: number.NewConstantProbabilityGenerator(number.ProbabilityZero),
		Mutater:             mutater,
	}

	// set client options
	for _, option := range options {
		if err := option.Apply(op); err != nil {
			return nil, fmt.Errorf("can't apply abstract mutation option: %v", err)
		}
	}
	return op, nil
}

// Apply applies the mutation operation to each entry in the list of selected
// candidates.
func (op *AbstractMutation) Apply(selectedCandidates []framework.Candidate, rng *rand.Rand) []framework.Candidate {
	mutatedPopulation := make([]framework.Candidate, len(selectedCandidates))
	for i, candidate := range selectedCandidates {
		mutatedPopulation[i] = op.Mutate(candidate, rng)
	}
	return mutatedPopulation
}
