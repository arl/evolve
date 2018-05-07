package mutation

import (
	"errors"
	"math"
	"math/rand"

	"github.com/aurelien-rainone/evolve/pkg/api"
	"github.com/aurelien-rainone/evolve/pkg/bitstring"
)

// ErrInvalidMutationCount is the error returned when trying to set an invalid
// mutation count
var ErrInvalidMutationCount = errors.New("mutation count must be in the [0,MaxInt32] range")

// TODO: document + rename (with pkg name it could be named mutation.Bitstring)
type BitStringMutation struct {
	*Mutation
	nmut             int
	varnmut          bool
	nmutmin, nmutmax int
}

// NewBitstringMutation creates an evolutionary operator that mutates individual
// bits in a bitstring.Bitstring according to some probability.
//
// Possible options:
// the mutation probability is the (possibly variable) probability of a
// candidate bit string being mutated at all; set it with ConstantProbability or
// VariableProbability. The default is a constant probability of 1.
// The mutation count is the (possibly variable) number of bits that will be
// flipped on any candidate bit string that is selected for mutation; set it
// with ConstantMutationCount or VariableMutationCount. The default is a
// constant mutation count of exactly 1 bit flipped.
func NewBitstringMutation() *BitStringMutation {
	bsmut := &BitStringMutation{
		nmut: 1, varnmut: false, nmutmin: 1, nmutmax: 1,
	}
	bsmut.Mutation = NewMutation(bsmut)
	return bsmut
}

// SetMutations sets the number of mutations (i.e the number of bits that will
// be flipped if the bitstring candidate is selected for mutation).
//
// If nmut is not in the [0,MaxInt32] range SetMutations will return
// ErrInvalidMutationCount.
func (op *BitStringMutation) SetMutations(nmut int) error {
	if nmut < 0 || nmut > math.MaxInt32 {
		return ErrInvalidMutationCount
	}
	op.nmut = nmut
	op.varnmut = false
	return nil
}

// SetMutationsRange sets the range of possible number of mutations (i.e the
// numnber of bits that will be filpped if the bitstring candidate is selected
// for mutation).
//
// The specific number of mutations will be randomly chosen with the pseudo
// random number generator argument of Apply, by linearly converting from
// [0,MaxInt32) to [min,max).
//
// If min and max are not bounded by [0,MaxInt32] SetMutationsRange will return
func (op *BitStringMutation) SetMutationsRange(min, max int) error {
	if min > max || min < 0 || max > math.MaxInt32 {
		return ErrInvalidMutationCount
	}
	op.nmutmin = min
	op.nmutmax = max
	op.varnmut = true
	return nil
}

// Mutate mutates a single bit string. Zero or more bits may be flipped.
//
// The probability of any given bit being flipped is governed by the probability
// generator configured for this mutation operator.
func (op *BitStringMutation) Mutate(c api.Candidate, rng *rand.Rand) api.Candidate {
	// get/decide a probability for this run
	prob := op.prob
	if op.varprob {
		prob = op.probmin + (op.probmax-op.probmin)*rng.Float64()
	}

	if rng.Float64() < prob {
		bs := c.(*bitstring.Bitstring)
		mutated := bs.Copy()
		// there is a mutation to perform, get/decide how many bits to flip
		var nmut int
		if op.varnmut {
			nmut = op.nmutmin + rng.Intn(op.nmutmax-op.nmutmin)
		} else {
			nmut = op.nmut
		}

		for i := 0; i < nmut; i++ {
			mutated.FlipBit(rng.Intn(mutated.Len()))
		}
		return mutated
	}
	return c
}
