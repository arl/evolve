package mutation

import (
	"errors"
	"math/rand"

	"github.com/arl/evolve/generator"
	"github.com/arl/evolve/pkg/bitstring"
)

// ErrInvalidMutationCount is the error returned when trying to set an invalid
// mutation count
var ErrInvalidMutationCount = errors.New("mutation count must be in the [0,MaxInt32] range")

// A Bitstring mutates individual bits in a bitstring.Bitstring according to
// some probability.
//
// Probability is the probability of a bitstring being mutated at all.
// FlipCount is the the number of bits to flip on the bitstring in case it's
// selected for mutation.
type Bitstring struct {
	Probability generator.Float
	FlipCount   generator.Int
}

// Mutate modifies a bitstring.Bitstring with respect to a mutation probabilty.
func (op *Bitstring) Mutate(c interface{}, rng *rand.Rand) interface{} {
	// Find out the mutation probabilty
	prob := op.Probability.Next()

	if rng.Float64() < prob {
		bs := c.(*bitstring.Bitstring)
		mutated := bitstring.Copy(bs)

		// Since there's a mutation to perform, find out how many bits to flip.
		nmuts := op.FlipCount.Next()
		for i := int64(0); i < nmuts; i++ {
			mutated.FlipBit(uint(rng.Intn(mutated.Len())))
		}

		return mutated
	}

	return c
}
