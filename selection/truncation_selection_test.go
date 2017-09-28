package selection

import (
	"testing"

	"github.com/aurelien-rainone/evolve/framework"
	"github.com/stretchr/testify/assert"
)

// Unit test for truncation selection strategy. Ensures the correct candidates
// are selected.

func TestTruncationSelectionNaturalFitness(t *testing.T) {
	selector, err := NewTruncationSelection(WithConstantSelectionRatio(0.5))
	assert.NoError(t, err)
	population := make(framework.EvaluatedPopulation, 4)

	// Higher score is better.
	steve, _ := framework.NewEvaluatedCandidate("Steve", 10.0)
	mary, _ := framework.NewEvaluatedCandidate("Mary", 9.1)
	john, _ := framework.NewEvaluatedCandidate("John", 8.4)
	gary, _ := framework.NewEvaluatedCandidate("Gary", 6.2)

	population[0] = steve
	population[1] = mary
	population[2] = john
	population[3] = gary

	selection := selector.Select(population, true, 2, nil)

	assert.Len(t, selection, 2, "want selection size to be 2, got ", len(selection))
	assert.Contains(t, selection, steve.Candidate(), "best candidate not selected")
	assert.Contains(t, selection, mary.Candidate(), "second best candidate not selected")
}

func TestTruncationSelectionNonNaturalFitness(t *testing.T) {
	selector, err := NewTruncationSelection(WithConstantSelectionRatio(0.5))
	assert.NoError(t, err)
	population := make(framework.EvaluatedPopulation, 4)

	// Lower score is better.
	gary, _ := framework.NewEvaluatedCandidate("Gary", 6.2)
	john, _ := framework.NewEvaluatedCandidate("John", 8.4)
	mary, _ := framework.NewEvaluatedCandidate("Mary", 9.1)
	steve, _ := framework.NewEvaluatedCandidate("Steve", 10.0)

	population[0] = gary
	population[1] = john
	population[2] = mary
	population[3] = steve

	selection := selector.Select(population, false, 2, nil)
	assert.Len(t, selection, 2, "want selection size to be 2, got ", len(selection))

	assert.Contains(t, selection, gary.Candidate(), "best candidate not selected")
	assert.Contains(t, selection, john.Candidate(), "second best candidate not selected")
}

func TestTruncationSelectionZeroRatio(t *testing.T) {
	// The selection ratio must be greater than zero to be useful. This test
	// ensures that an appropriate exception is thrown if the ratio is not
	// positive.  Not throwing an exception is an error because it permits
	// undetected bugs in evolutionary programs.
	_, err := NewTruncationSelection(WithConstantSelectionRatio(0))
	assert.Error(t, err)
}

func TestTruncationSelectionRatioTooHigh(t *testing.T) {
	// The selection ratio must be less than 1 to be useful. This test ensures
	// that an appropriate exception is thrown if the ratio is too high.  Not
	// throwing an exception is an error because it permits undetected bugs in
	// evolutionary programs.
	_, err := NewTruncationSelection(WithConstantSelectionRatio(1))
	assert.Error(t, err)
}
