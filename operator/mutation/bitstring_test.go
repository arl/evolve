package mutation

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/arl/bitstring"
	"github.com/arl/evolve/generator"
)

// Ensures that mutation occurs correctly. Because of the random aspect we can't
// actually make many assertions. This just ensures that there are no unexpected
// exceptions and that the length of the bit strings remains as expected.
func TestBitstringMutationRandom(t *testing.T) {
	rng := rand.New(rand.NewSource(99))

	bs := &Bitstring{
		FlipCount:   generator.Const(1),
		Probability: generator.Const(0.5),
	}

	mut := New[*bitstring.Bitstring](bs)

	org, err := bitstring.NewFromString("111100101")
	assert.NoError(t, err)

	pop := []*bitstring.Bitstring{org}
	for i := 0; i < 20; i++ {
		// Perform several iterations to get different mutations.
		pop = mut.Apply(pop, rng)
		mutated := pop[0]
		assert.IsType(t, &bitstring.Bitstring{}, mutated)
		assert.Equalf(t, 9, mutated.Len(), "want mutated.Len() = 9, got %v", mutated.Len())
	}
}

// Ensures that mutation occurs correctly. Uses a probability of 1 to make the
// outcome predictable (all bits will be flipped).
func TestBitstringMutationSingleBit(t *testing.T) {
	rng := rand.New(rand.NewSource(99))

	bs := &Bitstring{
		FlipCount:   generator.Const(1),
		Probability: generator.Const(0.5),
	}

	mut := New[*bitstring.Bitstring](bs)

	org, err := bitstring.NewFromString("111100101")
	assert.NoError(t, err)

	pop := []*bitstring.Bitstring{org}
	pop = mut.Apply(pop, rng)

	mutated := pop[0]
	assert.IsType(t, &bitstring.Bitstring{}, mutated)

	assert.False(t, mutated.Equals(org), "want mutant to be different from original, got equals")
	assert.Equalf(t, 9, mutated.Len(), "want mutated bit string to not change length, 9, got %v", mutated.Len())
	set := mutated.OnesCount()
	unset := mutated.ZeroesCount()
	assert.Truef(t, set == 5 || set == 7, "want 5 or 7 set bits in mutated bit string, got %v", set)
	assert.Truef(t, unset == 2 || unset == 4, "want 2 or 4 unset bits in mutated bit string, got %v", unset)
}
