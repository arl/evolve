package factory

import (
	"math/rand"
	"testing"

	"github.com/arl/evolve"
	"github.com/arl/evolve/pkg/bitstring"
	"github.com/stretchr/testify/assert"
)

func TestBitstring(t *testing.T) {
	const (
		length  = 10
		popsize = 5
	)

	// local test function
	checkPop := func(pop []interface{}) {
		// Make sure the correct number of candidates were generated.
		assert.Lenf(t, pop, popsize, "want population size = %v, got %v", popsize, len(pop))
		// Make sure that each individual is the right length.
		for _, cand := range pop {
			bs := cand.(*bitstring.Bitstring)
			assert.Equalf(t, length, bs.Len(), "want bitstring length = %v, got %v", length, bs.Len())
		}
	}

	rng := rand.New(rand.NewSource(99))

	t.Run("unseed population", func(t *testing.T) {
		f := Bitstring(length)
		pop := evolve.GeneratePopulation(f, popsize, rng)
		checkPop(pop)
	})

	t.Run("seeded population", func(t *testing.T) {
		f := Bitstring(length)
		seed1, _ := bitstring.MakeFromString("1111100000")
		seed2, _ := bitstring.MakeFromString("1010101010")
		seeds := []interface{}{seed1, seed2}
		pop, err := evolve.SeedPopulation(f, popsize, seeds, rng)
		if err != nil {
			t.Error(err)
		}

		// Check that the seed candidates appear in the generated population.
		assert.Containsf(t, pop, seed1, "Population does not contain seed candidate 1.")
		assert.Containsf(t, pop, seed2, "Population does not contain seed candidate 2.")
		checkPop(pop)
	})

	t.Run("too many seed candidates", func(t *testing.T) {
		f := Bitstring(length)
		cand := bitstring.New(length)
		// The following call should panic since the 3 seed candidates won't fit
		// into a population of size 2.
		seeds := []interface{}{cand, cand, cand}
		_, err := evolve.SeedPopulation(f, 2, seeds, rng)
		if err != evolve.ErrTooManySeedCandidates {
			t.Errorf("wantErr = %v, got %v", evolve.ErrTooManySeedCandidates, err)
		}
	})
}
