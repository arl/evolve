package xover

import (
	"math/rand"
	"testing"

	"github.com/aurelien-rainone/evolve/pkg/api"
	"github.com/stretchr/testify/assert"
)

func TestStringMater(t *testing.T) {
	rng := rand.New(rand.NewSource(99))

	xover := New(StringMater{})
	pop := make([]api.Candidate, 4)
	pop[0] = "abcde"
	pop[1] = "fghij"
	pop[2] = "klmno"
	pop[3] = "pqrst"

	for i := 0; i < 20; i++ {
		values := make(map[rune]struct{}, 20) // used as a set of runes
		pop = xover.Apply(pop, rng)
		assert.Len(t, pop, 4, "Population size changed after crossover.")
		for _, individual := range pop {
			s := individual.(string)
			assert.Lenf(t, s, 5, "Invalid candidate length: %v", len(s))
			for _, value := range s {
				values[value] = struct{}{}
			}
		}
		// All of the individual elements should still be present, just jumbled up
		// between individuals.
		assert.Len(t, values, 20, "Information lost during crossover.")
	}
}

// StringMater is only defined to work on populations
// containing strings of equal lengths. Any attempt to apply the operation to
// populations that contain different length strings should panic.
func TestStringMaterWithDifferentLengthParents(t *testing.T) {
	rng := rand.New(rand.NewSource(99))

	xover := New(StringMater{})
	pop := []api.Candidate{"abcde", "fghijklm"}

	// This should panic since the parents are different lengths.
	// TODO: why panicking and not returning an error?
	assert.Panics(t, func() {
		xover.Apply(pop, rng)
	})
}
