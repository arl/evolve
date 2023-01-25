package mutation

import (
	"math/rand"

	"github.com/arl/bitstring"
	"github.com/arl/evolve"
	"github.com/arl/evolve/generator"
)

// Bitstring mutation mutates random individual bits in a population of
// *bitstring.Bitstring.
//
// Probability governs the probability of a bitstring to be mutated. FlipCount
// governs the number of bits to flip on a bitstring selected for mutation.
type Bitstring struct {
	Probability generator.Float
	FlipCount   generator.Generator[int]
}

// Apply mutates the provided population.
func (op *Bitstring) Apply(pop *evolve.Population[*bitstring.Bitstring], rng *rand.Rand) {
	for i := 0; i < pop.Len(); i++ {
		if rng.Float64() >= op.Probability.Next() {
			continue
		}

		// This candidate will be mutated. Find out the number of bits to flip.
		cand := pop.Candidates[i]
		nbits := op.FlipCount.Next()
		for i := 0; i < nbits; i++ {
			cand.FlipBit(rng.Intn(cand.Len()))
		}
	}
}
