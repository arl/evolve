package mutation

import (
	"math/rand"

	"github.com/arl/evolve/generator"
)

// String mutates individual characters (single bytes) in a string according to
// some mutation probabilty.
//
// The mutation probablity is generated once for each call to Mutate, then there
// might be zero or more modified characters. A modified character gets replaced
// by any other one, randomly chosen from the alphabet string.
type String struct {
	Alphabet string
	MutProb  generator.Float
}

// Mutate modifies a string with respect to a mutation probabilty.
func (op *String) Mutate(c interface{}, rng *rand.Rand) interface{} {
	s := c.(string)

	buffer := make([]byte, len(s))
	copy(buffer, []byte(s))

	// Find out the probability for this run.
	prob := op.MutProb.Next()

	for i := range buffer {
		if rng.Float64() < prob {
			buffer[i] = op.Alphabet[rng.Intn(len(op.Alphabet))]
		}
	}

	return string(buffer)
}
