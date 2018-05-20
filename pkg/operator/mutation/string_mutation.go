package mutation

import (
	"math/rand"
)

// NewString creates an operator that mutates individual characters in a
// string according to some probability.
//
// Zero or more characters may be modified. The probability of any given
// character being modified is governed by the mutation probability.
func NewString(alphabet string) *stringMutater { // nolint: golint
	smut := &stringMutater{alphabet: alphabet}
	smut.Mutation = New(smut)
	return smut
}

type stringMutater struct {
	*Mutation
	alphabet string
}

func (op *stringMutater) Mutate(c interface{}, rng *rand.Rand) interface{} {
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
