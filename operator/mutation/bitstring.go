package mutation

import (
	"errors"
	"math/rand"

	"github.com/arl/bitstring"
	"github.com/arl/evolve/generator"
)

// ErrInvalidMutationCount is the error returned when trying to set an invalid
// mutation count
var ErrInvalidMutationCount = errors.New("mutation count must be in [0, MaxInt32]")

// A Bitstring mutates individual bits in a bitstring.Bitstring according to
// some probability.
//
// Probability is the probability of a bitstring being mutated at all.
// FlipCount is the the number of bits to flip on the bitstring in case it's
// selected for mutation.
type Bitstring struct {
	Probability generator.Float
	FlipCount   generator.Generator[int]
}

// Mutate modifies a bitstring.Bitstring with respect to a mutation probabilty.
func (op *Bitstring) Mutate(bs *bitstring.Bitstring, rng *rand.Rand) *bitstring.Bitstring {
	// Find out the mutation probabilty
	prob := op.Probability.Next()

	if rng.Float64() < prob {
		mutated := bitstring.Clone(bs)

		// Since there's a mutation to perform, find out how many bits to flip.
		nmuts := op.FlipCount.Next()
		for i := 0; i < nmuts; i++ {
			mutated.FlipBit(rng.Intn(mutated.Len()))
		}

		return mutated
	}

	return bs
}
