package mutation

import (
	"math/rand"
	"testing"

	"github.com/arl/evolve"
	"github.com/arl/evolve/generator"
	"github.com/stretchr/testify/assert"
)

func TestPermutation(t *testing.T) {
	rng := rand.New(rand.NewSource(99))

	mut := Permutation[int]{
		Count:       generator.Const(1),
		Amount:      generator.Const(1),
		Probability: generator.Const(1.0),
	}

	org := []int{1, 2, 3, 4, 5}

	// Place another -equal- candidate in the population, so we can later
	// compare with the original
	cpy := []int{1, 2, 3, 4, 5}
	pop := evolve.NewPopulationOf([][]int{cpy}, nil)
	mut.Apply(pop, rng)

	mutant := pop.Candidates[0]
	t.Logf("original: %+v", org[0])
	t.Logf("mutant : %+v", mutant)

	if len(mutant) != len(org) {
		t.Fatalf("len(mutant) = %v, want %v", len(mutant), len(org))
	}

	// The mutant should have the same elements but in a different order.
	matches := 0
	for i := range mutant {
		if org[i] == mutant[i] {
			matches++
		} else {
			// If positions don't match, an adjacent item should be a match.
			next := (i + 1) % len(org)
			prev := ((i - 1) + len(org)) % len(org)
			matchAdjacent := (org[i] == mutant[next]) != (org[i] == mutant[prev])
			assert.True(t, matchAdjacent, "Mutated item is not in one of the expected positions")
		}
	}

	assert.Equal(t, matches, len(org)-2, "All but 2 positions should be unchanged.")
}
