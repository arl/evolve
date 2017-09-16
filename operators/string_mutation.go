package operators

import (
	"fmt"
	"math/rand"

	"github.com/aurelien-rainone/evolve/framework"
	"github.com/aurelien-rainone/evolve/number"
)

// WithConstantStringMutationProbability sets up a constant probability that,
// once selected, a candidate will be mutated.
func WithConstantStringMutationProbability(mutationProbability number.Probability) StringMutationOption {
	return func(op *StringMutation) error {
		op.mutationProbability = number.NewConstantProbabilityGenerator(mutationProbability)
		return nil
	}
}

// WithVariableStringMutationProbability sets up a variable probability that,
// once selected, a candidate will be mutated.
func WithVariableStringMutationProbability(variable number.ProbabilityGenerator) StringMutationOption {
	return func(op *StringMutation) error {
		op.mutationProbability = variable
		return nil
	}
}

// StringMutation is an evolutionary operator that mutates individual
// characters in a string according to some probability.
type StringMutation struct {
	alphabet            string
	mutationProbability number.ProbabilityGenerator
}

// StringMutationOption is the type of the functions used to set string mutation
// options.
type StringMutationOption func(*StringMutation) error

// NewStringMutation creates a StringMutation configured with the provided
// options.
func NewStringMutation(alphabet string, options ...StringMutationOption) (*StringMutation, error) {
	// create with default options, mutation probability of zero
	op := &StringMutation{
		alphabet:            alphabet,
		mutationProbability: number.NewConstantProbabilityGenerator(number.ProbabilityZero),
	}

	// set client options
	for _, option := range options {
		if err := option(op); err != nil {
			return op, fmt.Errorf("can't apply string mutation options: %v", err)
		}
	}
	return op, nil
}

// Apply applies the string mutation to each entry in the list of selected
// candidates.
func (op *StringMutation) Apply(selectedCandidates []framework.Candidate, rng *rand.Rand) []framework.Candidate {
	mutatedPopulation := make([]framework.Candidate, len(selectedCandidates))
	for i, candidate := range selectedCandidates {
		mutatedPopulation[i] = op.mutateString(candidate.(string), rng)
	}
	return mutatedPopulation
}

// mutateString mutates a single string.
//
// Zero or more characters may be modified. The probability of any given
// character being modified is governed by the probability generator configured
// for this mutation operator.
// Returns the mutated string
func (op *StringMutation) mutateString(s string, rng *rand.Rand) string {
	buffer := make([]byte, len(s))
	copy(buffer, []byte(s))
	for i := range buffer {
		if op.mutationProbability.NextValue().NextEvent(rng) {
			buffer[i] = op.alphabet[rng.Intn(len(op.alphabet))]
		}
	}
	return string(buffer)
}
