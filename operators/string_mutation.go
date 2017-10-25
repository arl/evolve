package operators

import (
	"math/rand"

	"github.com/aurelien-rainone/evolve/framework"
)

// NewStringMutation creates an evolutionary operator that mutates individual
// characters in a string according to some probability.
//
// Zero or more characters may be modified. The probability of any given
// character being modified is governed by the probability generator configured
// for this mutation operator, using ConstantProbability or VariableProbability.
func NewStringMutation(alphabet string, options ...Option) (*AbstractMutation, error) {
	mutater := &stringMutater{alphabet: alphabet}
	impl, err := NewAbstractMutation(mutater, options...)
	// in the current case the actual mutater needs the abstract mutation
	// implementation back in order to access the mutation probability.
	mutater.impl = impl
	return impl, err
}

type stringMutater struct {
	alphabet string
	impl     *AbstractMutation
}

func (op *stringMutater) Mutate(c framework.Candidate, rng *rand.Rand) framework.Candidate {
	s := c.(string)
	buffer := make([]byte, len(s))
	copy(buffer, []byte(s))
	for i := range buffer {
		if op.impl.mutationProbability.NextValue().NextEvent(rng) {
			buffer[i] = op.alphabet[rng.Intn(len(op.alphabet))]
		}
	}
	return string(buffer)
}
