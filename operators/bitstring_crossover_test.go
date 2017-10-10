package operators

import (
	"math/rand"
	"testing"

	"github.com/aurelien-rainone/evolve/bitstring"
	"github.com/aurelien-rainone/evolve/factory"
	"github.com/aurelien-rainone/evolve/framework"
	"github.com/stretchr/testify/assert"
)

func TestBitStringCrossover(t *testing.T) {
	rng := rand.New(rand.NewSource(99))
	crossover, _ := NewBitStringCrossover()
	factory := factory.NewBitStringFactory(50)

	population := factory.GenerateInitialPopulation(2, rng)
	// Test to make sure that cross-over correctly preserves all genetic material
	// originally present in the population and does not introduce anything new.
	totalSetBits := population[0].(*bitstring.BitString).CountSetBits() +
		population[1].(*bitstring.BitString).CountSetBits()
	for i := 0; i < 50; i++ {
		// Test several generations.
		population = crossover.Apply(population, rng)

		setBits := population[0].(*bitstring.BitString).CountSetBits() +
			population[1].(*bitstring.BitString).CountSetBits()
		assert.Equal(t, setBits, totalSetBits, "bitstring crossover should not change the total number of set bits in population")
	}
}

func TestBitStringCrossoveWithDifferentLengthParents(t *testing.T) {
	// The BitStringCrossover operator is only defined to work on populations
	// containing Strings of equal lengths. Any attempt to apply the operation
	// to populations that contain different length strings should panic. Not
	// panicking should be considered a bug since it could lead to hard to trace
	// bugs elsewhere.
	rng := rand.New(rand.NewSource(99))
	crossover, _ := NewBitStringCrossover(WithConstantCrossoverPoints(1))

	bs1, _ := bitstring.NewRandom(32, rng)
	bs2, _ := bitstring.NewRandom(33, rng)
	population := []framework.Candidate{bs1, bs2}

	assert.Panics(t, func() {
		// This should panic since the parents are different lengths.
		crossover.Apply(population, rng)
	})
}
