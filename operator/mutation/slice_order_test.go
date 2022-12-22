package mutation

import (
	"math/rand"
	"testing"

	"github.com/arl/evolve/generator"
	"github.com/stretchr/testify/assert"
)

func TestSliceOrderMutation(t *testing.T) {
	rng := rand.New(rand.NewSource(99))
	pop := [][]int{{1, 2, 3, 4, 5}}

	op := SliceOrder[int]{
		Count:       generator.Const(1),
		Amount:      generator.Const(1),
		Probability: generator.Const(1.0),
	}
	mutpop := op.Apply(pop, rng)

	assert.Len(t, mutpop, len(pop), "population size should be unchanged after mutation")

	mutant := mutpop[0]
	t.Logf("original: %+v", pop[0])
	t.Logf("mutant : %+v", mutant)

	assert.Len(t, mutant, len(pop[0]), "mutant should be same length as original")

	// The mutant should have the same elements but in a different order
	matches := 0
	for i := range pop[0] {
		cand := pop[0]
		if cand[i] == mutant[i] {
			matches++
		} else {
			// If positions don't match, an adjacent item should be a match.
			next := (i + 1) % len(cand)
			prev := ((i - 1) + len(cand)) % len(cand)
			matchAdjacent := (cand[i] == mutant[next]) != (cand[i] == mutant[prev])
			assert.True(t, matchAdjacent, "Mutated item is not in one of the expected positions")
		}
	}

	assert.Equal(t, matches, len(pop[0])-2, "All but 2 positions should be unchanged.")
}
