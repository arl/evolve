package mutation

import (
	"math/rand"
	"testing"

	"github.com/aurelien-rainone/evolve/pkg/api"
	"github.com/aurelien-rainone/evolve/pkg/bitstring"
	"github.com/stretchr/testify/assert"
)

func errcheck(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("want error = nil, got %v", err)
	}
}

// Ensures that mutation occurs correctly. Because of the random aspect we
// can't actually make many assertions. This just ensures that there are no
// unexpected exceptions and that the length of the bit strings remains as
// expected.
func TestBitstringMutationRandom(t *testing.T) {
	rng := rand.New(rand.NewSource(99))

	mut := NewBitstringMutation()
	errcheck(t, mut.SetProb(0.5))

	org, err := bitstring.MakeFromString("111100101")
	errcheck(t, err)

	pop := []api.Candidate{org}
	for i := 0; i < 20; i++ {
		// Perform several iterations to get different mutations.
		pop = mut.Apply(pop, rng)
		mutated := pop[0]
		assert.IsType(t, &bitstring.Bitstring{}, mutated)
		ms := mutated.(*bitstring.Bitstring)
		assert.Equalf(t, 9, ms.Len(), "want ms.Len() = 9, got %v", ms.Len())
	}
}

// Ensures that mutation occurs correctly.  Uses a probability of 1 to
// make the outcome predictable (all bits will be flipped).
func TestBitstringMutationSingleBit(t *testing.T) {
	rng := rand.New(rand.NewSource(99))

	mut := NewBitstringMutation()
	errcheck(t, mut.SetProb(1.0))

	org, err := bitstring.MakeFromString("111100101")
	errcheck(t, err)

	pop := []api.Candidate{org}
	pop = mut.Apply(pop, rng)

	mutated := pop[0]
	assert.IsType(t, &bitstring.Bitstring{}, mutated)
	ms := mutated.(*bitstring.Bitstring)

	assert.False(t, ms.Equals(org), "want mutant to be different from original, got equals")
	assert.Equalf(t, 9, ms.Len(), "want mutated bit string to not change length, 9, got %v", ms.Len())
	set := ms.OnesCount()
	unset := ms.ZeroesCount()
	assert.Truef(t, set == 5 || set == 7, "want 5 or 7 set bits in mutated bit string, got %v", set)
	assert.Truef(t, unset == 2 || unset == 4, "want 2 or 4 unset bits in mutated bit string, got %v", unset)
}

// TODO:  test BitstringMutation constructed with options other than default ones
//func
