package selection

import (
	"math/rand"
	"testing"

	"github.com/aurelien-rainone/evolve/framework"
	"github.com/stretchr/testify/assert"
)

func errcheck(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("want error = nil, got %v", err)
	}
}

func TestTournamentSelectionNaturalFitness(t *testing.T) {
	rng := rand.New(rand.NewSource(99))

	tournament := NewTournamentSelection()
	errcheck(tournament.SetProb(0.7))

	steve, _ := framework.NewEvaluatedCandidate("Steve", 10.0)
	mary, _ := framework.NewEvaluatedCandidate("Mary", 9.1)
	john, _ := framework.NewEvaluatedCandidate("John", 8.4)
	gary, _ := framework.NewEvaluatedCandidate("Gary", 6.2)
	pop := framework.EvaluatedPopulation{steve, mary, john, gary}

	// Run several iterations so that we get different tournament outcomes.
	for i := 0; i < 20; i++ {
		selection := tournament.Select(pop, true, 2, rng)
		assert.Len(t, selection, 2, "want len(selection) = 2, got", len(selection))
	}
}

func TestTournamentSelectionNonNaturalFitness(t *testing.T) {
	rng := rand.New(rand.NewSource(99))

	tournament := NewTournamentSelection()
	errcheck(tournament.SetProb(0.7))

	gary, _ := framework.NewEvaluatedCandidate("Gary", 6.2)
	john, _ := framework.NewEvaluatedCandidate("John", 8.4)
	mary, _ := framework.NewEvaluatedCandidate("Mary", 9.1)
	steve, _ := framework.NewEvaluatedCandidate("Steve", 10.0)
	pop := framework.EvaluatedPopulation{gary, john, mary, steve}

	// Run several iterations so that we get different tournament outcomes.
	for i := 0; i < 20; i++ {
		selection := selector.Select(pop, false, 2, rng)
		assert.Len(t, selection, 2, "want len(selection) = 2, got", len(selection))
	}
}

// This test ensures that an error is returned if the probability is 0.5 or
// less.
func TestTournamentSelectionProbabilityTooLow(t *testing.T) {
	ts := NewTournamentSelection()
	err := ts.SetProb(0.5)
	if err != ErrInvalidTournamentProb {
		t.Errorf("want ts.SetProb(0.5) = ErrInvalidTournamentProb, got %v", err)
	}
	err := ts.SetProbRange(0.4, 0.6)
	if err != ErrInvalidTournamentProb {
		t.Errorf("want ts.SetProb(0.4, 0.6) = ErrInvalidTournamentProb, got %v", err)
	}
}
