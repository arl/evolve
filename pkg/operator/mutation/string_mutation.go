package operators

import (
	"math/rand"

	"github.com/aurelien-rainone/evolve/framework"
)

// NewStringMutation creates an operator that mutates individual characters in a
// string according to some probability.
//
// Zero or more characters may be modified. The probability of any given
// character being modified is governed by the mutation probability.
func NewStringMutation(alphabet string) *stringMutater {
	smut := &stringMutater{alphabet: alphabet}
	smut.Mutation = NewMutation(smut)
	return smut
}

type stringMutater struct {
	*Mutation
	alphabet string
}

func (op *stringMutater) Mutate(c framework.Candidate, rng *rand.Rand) framework.Candidate {
	s := c.(string)
	buffer := make([]byte, len(s))
	copy(buffer, []byte(s))

	// get/decide a probability for this run
	prob := op.prob
	if op.varprob {
		prob = op.probmin + (op.probmax-op.probmin)*rng.Float64()
	}

	for i := range buffer {
		if rng.Float64() < prob {
			buffer[i] = op.alphabet[rng.Intn(len(op.alphabet))]
		}
	}
	return string(buffer)
}
