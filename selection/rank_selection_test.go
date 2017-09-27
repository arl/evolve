package selection

import (
	"math/rand"
	"testing"

	"github.com/aurelien-rainone/evolve/framework"
	"github.com/stretchr/testify/assert"
)

// Test selection when fitness scoring is natural (higher is better).
func TestRankSelectionNaturalFitness(t *testing.T) {
	rng := rand.New(rand.NewSource(99))
	selector := NewRankSelection(nil)
	population := make(framework.EvaluatedPopulation, 4)

	// Higher score is better.
	steve, _ := framework.NewEvaluatedCandidate("Steve", 10.0)
	john, _ := framework.NewEvaluatedCandidate("John", 4.5)
	mary, _ := framework.NewEvaluatedCandidate("Mary", 1.0)
	gary, _ := framework.NewEvaluatedCandidate("Gary", 0.5)

	population[0] = steve
	population[1] = john
	population[2] = mary
	population[3] = gary

	selection := selector.Select(population, true, 4, rng)
	assert.Len(t, selection, 4, "got selection size:", len(selection), "want 4")

	steveCount := frequency(selection, steve.Candidate())
	johnCount := frequency(selection, john.Candidate())
	garyCount := frequency(selection, gary.Candidate())
	maryCount := frequency(selection, mary.Candidate())

	assert.True(t, steveCount >= 1 && steveCount <= 2, "candidate selected wrong number of times (should be 1 or 2, was ", steveCount, ")")
	assert.True(t, johnCount >= 1 && johnCount <= 2, "candidate selected wrong number of times (should be 1 or 2, was ", johnCount, ")")
	assert.True(t, garyCount <= 1, "Candidate selected wrong number of times (should be 0 or 1, was ", garyCount, ")")
	assert.True(t, maryCount <= 1, "Candidate selected wrong number of times (should be 0 or 1, was ", maryCount, ")")
}

// Test selection when fitness scoring is non-natural (lower is better).
func TestRankSelectionNonNaturalFitness(t *testing.T) {
	rng := rand.New(rand.NewSource(99))
	selector := NewRankSelection(nil)
	population := make(framework.EvaluatedPopulation, 4)

	// Lower score is better.
	gary, _ := framework.NewEvaluatedCandidate("Gary", 0.5)
	mary, _ := framework.NewEvaluatedCandidate("Mary", 1.0)
	john, _ := framework.NewEvaluatedCandidate("John", 4.5)
	steve, _ := framework.NewEvaluatedCandidate("Steve", 10.0)

	population[0] = gary
	population[1] = mary
	population[2] = john
	population[3] = steve

	selection := selector.Select(population, false, 4, rng)
	assert.Len(t, selection, 4, "got selection size:", len(selection), "want 4")

	garyCount := frequency(selection, gary.Candidate())
	maryCount := frequency(selection, mary.Candidate())
	johnCount := frequency(selection, john.Candidate())
	steveCount := frequency(selection, steve.Candidate())

	assert.True(t, garyCount >= 1 && garyCount <= 2, "candidate selected wrong number of times (should be 1 or 2, was ", garyCount, ")")
	assert.True(t, maryCount >= 1 && maryCount <= 2, "candidate selected wrong number of times (should be 1 or 2, was ", maryCount, ")")
	assert.True(t, johnCount <= 1, "candidate selected wrong number of times (should be 0 or 1, was ", johnCount, ")")
	assert.True(t, steveCount <= 1, "candidate selected wrong number of times (should be 0 or 1, was ", steveCount, ")")
}
