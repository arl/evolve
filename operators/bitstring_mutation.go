package operators

import (
	"math/rand"

	"github.com/aurelien-rainone/evolve/framework"
	"github.com/aurelien-rainone/evolve/number"
	"github.com/aurelien-rainone/evolve/pkg/bitstring"
)

// NewBitStringMutation creates an evolutionary operator that mutates individual
// bits in a bitstring.BitString according to some probability.
//
// Possible options:
// the mutation probability is the (possibly variable) probability of a
// candidate bit string being mutated at all; set it with ConstantProbability or
// VariableProbability. The default is a constant probability of 1.
// the mutation count is the (possibly variable) number of bits that will be
// flipped on any candidate bit string that is selected for mutation; set it
// with ConstantMutationCount or VariableMutationCount. The default is a
// constant mutation count of exactly 1 bit flipped.
func NewBitStringMutation(options ...Option) (*AbstractMutation, error) {
	// set default mutation count to 1
	mutater := &bitStringMutater{
		mutationCount: number.NewConstantIntegerGenerator(1),
	}
	// set default mutation probability to 1, this option is prepended to the
	// client slice of options, so it will be applied before and remains as
	// default if ever it's not overwritten by subsequent (in the slice)
	// options.
	impl, err := NewAbstractMutation(
		mutater,
		append([]Option{ConstantProbability(number.ProbabilityOne)}, options...)...,
	)
	// in the current case the actual mutater needs the abstract mutation
	// implementation back in order to access the mutation probability.
	mutater.impl = impl
	return impl, err
}

type bitStringMutater struct {
	mutationCount number.IntegerGenerator
	impl          *AbstractMutation
}

// Mutate mutates a single bit string. Zero or more bits may be flipped.
//
// The probability of any given bit being flipped is governed by the probability
// generator configured for this mutation operator.
func (op *bitStringMutater) Mutate(c framework.Candidate, rng *rand.Rand) framework.Candidate {
	if op.impl.mutationProbability.NextValue().NextEvent(rng) {
		bitString := c.(*bitstring.BitString)
		mutatedBitString := bitString.Copy()
		mutations := op.mutationCount.NextValue()
		for i := int64(0); i < mutations; i++ {
			mutatedBitString.FlipBit(rng.Intn(mutatedBitString.Len()))
		}
		return mutatedBitString
	}
	return c
}
