package mutation

import (
	"errors"
	"math/rand"
)

// ErrInvalidMutationProb is the error returned when trying to set an invalid
// mutation probability
var ErrInvalidMutationProb = errors.New("mutation probability must be in the [0,1] range")

// Mutater is the interface that wraps the Mutate method.
type Mutater interface {

	// Mutate performs mutation on a candidate.
	//
	// The original candidate is let untouched while the mutant is returned.
	Mutate(interface{}, *rand.Rand) interface{}
}

// Mutation implements a mutation operator.
//
// It supports all mutation processes that operate on an unique candidate.
// The mutation probability is configurable, its effect depends on the specific
// mutation implementation, where it will be documented.
//
// Note: unless you are implementing your own mutation operator, you generally
// don't need to directly instantiate a Mutation as specific mutation operators
// like BitStringMutation, StringMutation, etc. already create and embed a
// Mutation.
type Mutation struct {
	Mutater
	prob             float64
	varprob          bool
	probmin, probmax float64
}

// New creates a Mutation operator with rhe provided Mutater.
//
// The returned Mutation is preconfigured with a 0.01 mutation probability.
func New(mut Mutater) *Mutation {
	return &Mutation{
		Mutater: mut,
		prob:    0.01, varprob: false, probmin: 0.01, probmax: 0.01,
	}
}

// SetProb sets the mutation probability,
//
// If prob is not in the [0,1] range SetProb will return
// ErrInvalidMutationCount.
func (op *Mutation) SetProb(prob float64) error {
	if prob < 0.0 || prob > 1.0 {
		return ErrInvalidMutationProb
	}
	op.prob = prob
	op.varprob = false
	return nil
}

// SetProbRange sets the range of possible mutation probabilities.
//
// The specific mutation probability will be randomly chosen with the pseudo
// random number generator argument of Apply, by linearly converting from
// [0,1) to [min,max).
//
// If min and max are not bounded by [0,1] SetProbRange will return
// ErrInvalidMutationProb.
func (op *Mutation) SetProbRange(min, max float64) error {
	if min > max || min < 0.0 || max > 1.0 {
		return ErrInvalidMutationProb
	}
	op.probmin = min
	op.probmax = max
	op.varprob = true
	return nil
}

// Apply applies the mutation operator to each entry in the list of selected
// candidates.
func (op *Mutation) Apply(sel []interface{}, rng *rand.Rand) []interface{} {
	mutpop := make([]interface{}, len(sel))
	for i, cand := range sel {
		mutpop[i] = op.Mutate(cand, rng)
	}
	return mutpop
}
