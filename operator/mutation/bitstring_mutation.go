package mutation

import (
	"errors"
	"fmt"
	"math"
	"math/rand"

	"github.com/arl/evolve/pkg/bitstring"
)

// ErrInvalidMutationCount is the error returned when trying to set an invalid
// mutation count
var ErrInvalidMutationCount = errors.New("mutation count must be in the [0,MaxInt32] range")

// Bitstring is a an evolutionary operator that mutates individual bits in
// a bitstring.Bitstring according to some probability.
type Bitstring struct {
	*Mutation
	nmut             int
	varnmut          bool
	nmutmin, nmutmax int
}

// NewBitstring creates a new BitString mutation operator, pre-configured with a
// probability of mutation of 1.0 and mutation count of 1.
//
// The mutation probability is the (possibly variable) probability of a
// candidate bit string being mutated at all. It can be modified with SetProb
// and SetProbRange.
// The mutation count is the (possibly variable) number of bits that will be
// flipped on any candidate bit string that is selected for mutation. It can be
// modified with SetMutations and SetMutationsRange.
func NewBitstring() *Bitstring {
	bsmut := &Bitstring{
		nmut: 1, varnmut: false, nmutmin: 1, nmutmax: 1,
	}
	bsmut.Mutation = New(bsmut)
	err := bsmut.SetProb(1.0)
	if err != nil {
		panic(fmt.Errorf("cannot set mutation probability: %v", err))
	}
	return bsmut
}

// SetMutations sets the number of mutations (i.e the number of bits that will
// be flipped if the bitstring candidate is selected for mutation).
//
// If nmut is not in the [0,MaxInt32] range SetMutations will return
// ErrInvalidMutationCount.
func (op *Bitstring) SetMutations(nmut int) error {
	if nmut < 0 || nmut > math.MaxInt32 {
		return ErrInvalidMutationCount
	}
	op.nmut = nmut
	op.varnmut = false
	return nil
}

// SetMutationsRange sets the range of possible number of mutations (i.e the
// numnber of bits that will be flipped if the bitstring candidate is selected
// for mutation).
//
// The specific number of mutations will be randomly chosen with the pseudo
// random number generator argument of Apply, by linearly converting from
// [0,MaxInt32) to [min,max).
//
// If min and max are not bounded by [0,MaxInt32] SetMutationsRange will return
func (op *Bitstring) SetMutationsRange(min, max int) error {
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
func (op *Bitstring) Mutate(c interface{}, rng *rand.Rand) interface{} {
	// get/decide a probability for this run
	prob := op.prob
	if op.varprob {
		prob = op.probmin + (op.probmax-op.probmin)*rng.Float64()
	}

	if rng.Float64() < prob {
		bs := c.(*bitstring.Bitstring)
		mutated := bitstring.Copy(bs)
		// there is a mutation to perform, get/decide how many bits to flip
		var nmut int
		if op.varnmut {
			nmut = op.nmutmin + rng.Intn(op.nmutmax-op.nmutmin)
		} else {
			nmut = op.nmut
		}

		for i := 0; i < nmut; i++ {
			mutated.FlipBit(uint(rng.Intn(mutated.Len())))
		}
		return mutated
	}
	return c
}
