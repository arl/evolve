package mutation

import (
	"math/rand"

	"github.com/arl/evolve"
	"github.com/arl/evolve/generator"
)

// String mutates individual characters (single bytes) in a string according to
// some mutation probabilty.
//
// Probability governs the probabilty for a character to be mutated. When
// mutation happens, a new character gets randomly picked in the provided
// Alphabet.
type String struct {
	Alphabet    string
	Probability generator.Float
}

// Apply mutates the provided population.
func (op *String) Apply(pop *evolve.Population[string], rng *rand.Rand) {
	for i := 0; i < pop.Len(); i++ {
		// Find out the probability for this candidate.
		prob := op.Probability.Next()

		buf := []byte(pop.Candidates[i])
		for j := range buf {
			if rng.Float64() < prob {
				buf[j] = op.Alphabet[rng.Intn(len(op.Alphabet))]
			}
		}
		pop.Candidates[i] = string(buf)
	}
}
