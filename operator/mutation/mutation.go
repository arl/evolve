package mutation

import (
	"math/rand"
)

// Mutation implements the mutation evolutionnary operator. It modifies the
// genetic content of individuals in order to maintain diversity from one
// population to the next.
//
// At individual level, mutation is applied through a Mutater, which performs
// modification on a single element at once.
type Mutation struct {
	Mutater
}

// New Mutation returns an Operator based on mutater.
func NewMutation(mutater Mutater) *Mutation {
	return &Mutation{Mutater: mutater}
}

// Apply applies the mutation operator to all individuals in the provided population.
func (op *Mutation) Apply(population []interface{}, rng *rand.Rand) []interface{} {
	muted := make([]interface{}, len(population))
	for i, cand := range population {
		muted[i] = op.Mutate(cand, rng)
	}
	return muted
}

// A Mutater mutates individuals.
type Mutater interface {

	// Mutate performs mutation on an individual.
	//
	// The original individual must be let untouched while the mutant is
	// returned, unless no mutation is performed, in which case the original
	// individual can be returned.
	Mutate(interface{}, *rand.Rand) interface{}
}
