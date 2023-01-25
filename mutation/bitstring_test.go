package mutation

import (
	"math/rand"
	"testing"

	"github.com/arl/bitstring"
	"github.com/arl/evolve"
	"github.com/arl/evolve/generator"
)

func TestBitstringMutation(t *testing.T) {
	rng := rand.New(rand.NewSource(99))

	bs, err := bitstring.NewFromString("111100101")
	if err != nil {
		t.Fatal(err)
	}

	// Uses a probability of 1 to make the outcome predictable (all bits are
	// flipped).
	mut := &Bitstring{
		FlipCount:   generator.Const(1),
		Probability: generator.Const(0.5),
	}

	org := bitstring.Clone(bs)
	pop := evolve.NewPopulationOf([]*bitstring.Bitstring{bs}, nil)
	mut.Apply(pop, rng)

	mutated := pop.Candidates[0]
	if mutated.Equals(org) {
		t.Errorf("mutated and original are equals, should be different")
	}
	if mutated.Len() != org.Len() {
		t.Errorf("mutated.Len() = %v, want same length as original (%d)", mutated.Len(), org.Len())
	}
	if ones := mutated.OnesCount(); ones != 7 {
		t.Errorf("mutated string has %d ones, want %d", ones, 7)
	}
	if zeroes := mutated.ZeroesCount(); zeroes != 2 {
		t.Errorf("mutated string has %d zeroes, want %d", zeroes, 2)
	}
}
