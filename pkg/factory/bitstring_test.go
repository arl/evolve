package factory

import (
	"math/rand"
	"testing"

	"github.com/aurelien-rainone/evolve/pkg/bitstring"
	"github.com/stretchr/testify/assert"
)

func TestBitstringFactory(t *testing.T) {
	const (
		candidateLength = 10
		populationSize  = 5
	)

	// TODO: rename variables

	// local test function
	validatePopulation := func(population []interface{}) {
		// Make sure the correct number of candidates were generated.
		assert.Lenf(t, population, populationSize, "want population size = %v, got %v", populationSize, len(population))
		// Make sure that each individual is the right length.
		for _, cand := range population {
			bitString := cand.(*bitstring.Bitstring)
			assert.Equalf(t, candidateLength, bitString.Len(), "want bitstring length = %v, got %v", candidateLength, bitString.Len())
		}
	}

	rng := rand.New(rand.NewSource(99))

	t.Run("unseed population", func(t *testing.T) {

		f := NewBitstring(candidateLength)
		population := f.GenPopulation(populationSize, rng)
		validatePopulation(population)
	})

	t.Run("seeded population", func(t *testing.T) {

		f := NewBitstring(candidateLength)
		seed1, _ := bitstring.MakeFromString("1111100000")
		seed2, _ := bitstring.MakeFromString("1010101010")
		seeds := []interface{}{seed1, seed2}
		population := f.SeedPopulation(populationSize, seeds, rng)

		// Check that the seed candidates appear in the generated population.
		assert.Containsf(t, population, seed1, "Population does not contain seed candidate 1.")
		assert.Containsf(t, population, seed2, "Population does not contain seed candidate 2.")
		validatePopulation(population)
	})

	t.Run("too many seed candidates", func(t *testing.T) {

		f := NewBitstring(candidateLength)
		candidate, _ := bitstring.New(candidateLength)
		// The following call should panic since the 3 seed candidates won't fit
		// into a population of size 2.
		assert.Panics(t, func() {
			f.SeedPopulation(2,
				[]interface{}{candidate, candidate, candidate},
				rng)
		})
	})
}
