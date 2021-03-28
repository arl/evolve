package mutation

import (
	"math/rand"

	"github.com/arl/evolve/generator"
)

// String mutates individual characters (single bytes) in a string according to
// some mutation probabilty.
//
// Probability governs the probabilty for each character to be modified by
// Mutate. If this mutation happens, the mutated character gets replaced by any
// character in Alphabet.
type String struct {
	Alphabet    string
	Probability generator.Float
}

// Mutate modifies a string with respect to a mutation probabilty.
func (op *String) Mutate(c interface{}, rng *rand.Rand) interface{} {
	s := c.(string)

	buffer := make([]byte, len(s))
	copy(buffer, []byte(s))

	// Find out the probability for this run.
	prob := op.Probability.Next()

	for i := range buffer {
		if rng.Float64() < prob {
			buffer[i] = op.Alphabet[rng.Intn(len(op.Alphabet))]
		}
	}

	return string(buffer)
}
