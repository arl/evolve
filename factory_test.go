package evolve

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

type intFactory struct{}

func (intFactory) New(rng *rand.Rand) int { return rng.Int() }

func generateInt(rng *rand.Rand) int { return rng.Int() }

func testGeneratePopulation(t *testing.T, g Factory[int]) {
	rng := rand.New(rand.NewSource(99))

	pop := GeneratePopulation[int](intFactory{}, 10, rng)
	assert.Len(t, pop, 10)
}

func testSeedPopulation(t *testing.T, g Factory[int]) {
	rng := rand.New(rand.NewSource(99))

	// seed 5 candidates over 10
	seeds := make([]int, 5)
	for i := 0; i < 5; i++ {
		seeds[i] = i
	}

	pop, err := SeedPopulation[int](intFactory{}, 10, seeds, rng)
	assert.NoError(t, err)
	assert.Len(t, pop, 10)
}

func testSeedPopulationError(t *testing.T, g Factory[int]) {
	rng := rand.New(rand.NewSource(99))

	seeds := make([]int, 10)
	for i := 0; i < 10; i++ {
		seeds[i] = i
	}

	pop, err := SeedPopulation[int](intFactory{}, 5, seeds, rng)
	assert.Nil(t, pop)
	assert.ErrorIs(t, err, ErrTooManySeedCandidates)
}

func TestGeneratePopulation(t *testing.T) {
	testGeneratePopulation(t, intFactory{})
}

func TestSeedPopulation(t *testing.T) {
	testSeedPopulation(t, intFactory{})
}

func TestSeedPopulationError(t *testing.T) {
	testSeedPopulationError(t, intFactory{})
}

func TestGeneratePopulationFunc(t *testing.T) {
	testGeneratePopulation(t, FactoryFunc[int](generateInt))
}

func TestSeedPopulationFunc(t *testing.T) {
	testSeedPopulation(t, FactoryFunc[int](generateInt))
}

func TestSeedPopulationErrorFunc(t *testing.T) {
	testSeedPopulationError(t, FactoryFunc[int](generateInt))
}
