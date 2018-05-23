package api

import (
	"math/rand"
	"testing"
)

type intGenerator struct{}

func (intGenerator) Generate(rng *rand.Rand) interface{} { return rng.Int() }

func generateInt(rng *rand.Rand) interface{} { return rng.Int() }

func testGeneratePopulation(t *testing.T, g Generator) {
	rng := rand.New(rand.NewSource(99))

	pop := GeneratePopulation(intGenerator{}, 10, rng)
	if len(pop) != 10 {
		t.Errorf("GeneratePopulation: want len(pop) = %v, got %v", 10, len(pop))
	}
}

func testSeedPopulation(t *testing.T, g Generator) {
	rng := rand.New(rand.NewSource(99))

	// seed 5 candidates over 10
	seeds := make([]interface{}, 5)
	for i := 0; i < 5; i++ {
		seeds[i] = i
	}

	pop, err := SeedPopulation(intGenerator{}, 10, seeds, rng)
	if len(pop) != 10 {
		t.Errorf("SeedPopulation: want len(pop) = %v, got %v", 10, len(pop))
	}
	if err != nil {
		t.Errorf("SeedPopulation: want err = nil, got %v", err)
	}
}

func testSeedPopulationError(t *testing.T, g Generator) {
	rng := rand.New(rand.NewSource(99))

	seeds := make([]interface{}, 10)
	for i := 0; i < 10; i++ {
		seeds[i] = i
	}

	pop, err := SeedPopulation(intGenerator{}, 5, seeds, rng)
	if pop != nil {
		t.Errorf("SeedPopulation: want pop = nil, got %v", pop)
	}
	if err != ErrTooManySeedCandidates {
		t.Errorf("SeedPopulation: want err = %v, got %v", ErrTooManySeedCandidates, err)
	}
}

func TestGeneratePopulation(t *testing.T) {
	testGeneratePopulation(t, intGenerator{})
}
func TestSeedPopulation(t *testing.T) {
	testSeedPopulation(t, intGenerator{})
}
func TestSeedPopulationError(t *testing.T) {
	testSeedPopulationError(t, intGenerator{})
}

func TestGeneratePopulationFunc(t *testing.T) {
	testGeneratePopulation(t, GeneratorFunc(generateInt))
}
func TestSeedPopulationFunc(t *testing.T) {
	testSeedPopulation(t, GeneratorFunc(generateInt))
}
func TestSeedPopulationErrorFunc(t *testing.T) {
	testSeedPopulationError(t, GeneratorFunc(generateInt))
}
