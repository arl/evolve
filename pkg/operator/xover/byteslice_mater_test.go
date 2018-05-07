package xover

import (
	"math/rand"
	"testing"

	"github.com/aurelien-rainone/evolve/pkg/api"
	"github.com/stretchr/testify/assert"
)

func TestByteSliceMater(t *testing.T) {
	rng := rand.New(rand.NewSource(99))
	crossover := New(ByteSliceMater{})
	population := make([]api.Candidate, 4)
	population[0] = []byte{1, 2, 3, 4, 5}
	population[1] = []byte{6, 7, 8, 9, 10}
	population[2] = []byte{11, 12, 13, 14, 15}
	population[3] = []byte{16, 17, 18, 19, 20}

	for i := 0; i < 20; i++ {
		values := make(map[byte]struct{}, 20) // used as a set of runes
		population = crossover.Apply(population, rng)
		assert.Len(t, population, 4, "Population size changed after crossover.")
		for _, individual := range population {
			s := individual.([]byte)
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

// The ByteArrayCrossover operator is only defined to work on populations
// containing arrays of equal lengths. Any attempt to apply the operation to
// populations that contain different length arrays should panic. Not panicking
// should be considered a bug since it could lead to hard to trace bugs
// elsewhere.
func TestByteSliceMaterWithDifferentLengthParents(t *testing.T) {
	rng := rand.New(rand.NewSource(99))

	crossover := New(ByteSliceMater{})
	population := make([]api.Candidate, 2)
	population[0] = []byte{1, 2, 3, 4, 5}
	population[1] = []byte{2, 4, 8, 10, 12, 14, 16}

	// This should panic since the parents are different lengths.
	// TODO: why panicking and not returning an error?
	assert.Panics(t, func() {
		crossover.Apply(population, rng)
	})
}
