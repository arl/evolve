package selection

import (
	"math/rand"
	"testing"

	"github.com/aurelien-rainone/evolve/framework"
	"github.com/stretchr/testify/assert"
)

// Unit test for roulette selection strategy. We cannot easily test
// that the correct candidates are returned because of the random aspect
// of the selection, but we can at least make sure the right number of
// candidates are selected.
func TestRouletteWheelSelectionNaturalFitnessSolution(t *testing.T) {
	rng := rand.New(rand.NewSource(99))
	selector := &RouletteWheelSelection{}

	steve, _ := framework.NewEvaluatedCandidate("Steve", 10.0)
	mary, _ := framework.NewEvaluatedCandidate("Mary", 9.1)
	john, _ := framework.NewEvaluatedCandidate("John", 8.4)
	gary, _ := framework.NewEvaluatedCandidate("Gary", 6.2)
	population := framework.EvaluatedPopulation{steve, mary, john, gary}

	// Run several iterations to get different outcomes from the "roulette wheel".
	for i := 0; i < 20; i++ {
		selection := selector.Select(population, true, 2, rng)
		assert.Len(t, selection, 2, "want selection size = 2, got %v", len(selection))
	}
}

func TestRouletteWheelSelectionNonNaturalFitnessSolution(t *testing.T) {
	rng := rand.New(rand.NewSource(99))
	selector := &RouletteWheelSelection{}

	gary, _ := framework.NewEvaluatedCandidate("Gary", 6.2)
	john, _ := framework.NewEvaluatedCandidate("John", 8.4)
	mary, _ := framework.NewEvaluatedCandidate("Mary", 9.1)
	steve, _ := framework.NewEvaluatedCandidate("Steve", 10.0)
	population := framework.EvaluatedPopulation{gary, john, mary, steve}

	// Run several iterations to get different outcomes from the "roulette wheel".
	for i := 0; i < 20; i++ {
		selection := selector.Select(population, false, 2, rng)
		assert.Len(t, selection, 2, "want selection size = 2, got %v", len(selection))
	}
}

// Make sure that the code still functions for non-natural fitness scores even
// when one of them is a zero (a perfect score).
func TestRouletteWheelSelectionNonNaturalFitnessPerfectSolution(t *testing.T) {
	rng := rand.New(rand.NewSource(99))
	selector := &RouletteWheelSelection{}

	gary, _ := framework.NewEvaluatedCandidate("Gary", 0)
	john, _ := framework.NewEvaluatedCandidate("John", 8.4)
	mary, _ := framework.NewEvaluatedCandidate("Mary", 9.1)
	steve, _ := framework.NewEvaluatedCandidate("Steve", 10.0)
	population := framework.EvaluatedPopulation{gary, john, mary, steve}

	// Run several iterations to get different outcomes from the "roulette wheel".
	for i := 0; i < 20; i++ {
		selection := selector.Select(population, false, 2, rng)
		assert.Len(t, selection, 2, "want selection size = 2, got %v", len(selection))
	}
}
