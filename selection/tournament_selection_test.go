package selection

import (
	"fmt"
	"testing"

	"github.com/aurelien-rainone/evolve/framework"
)

func TestTournamentSelectionNatural(t *testing.T) {
	ts := NewTournamentSelection()
	errcheck(t, ts.SetProb(0.7))
	testRandomBasedSelection(t, ts, randomBasedPopNatural, true, 2,
		func(selected []framework.Candidate) error {
			if len(selected) != 2 {
				return fmt.Errorf("want len(selected) == 2, got %v", len(selected))
			}
			return nil
		})
}

func TestTournamentSelectionNonNatural(t *testing.T) {
	ts := NewTournamentSelection()
	errcheck(t, ts.SetProb(0.7))
	testRandomBasedSelection(t, ts, randomBasedPopNonNatural, false, 2,
		func(selected []framework.Candidate) error {
			if len(selected) != 2 {
				return fmt.Errorf("want len(selected) == 2, got %v", len(selected))
			}
			return nil
		})
}

func TestTournamentSelectionSetProb(t *testing.T) {
	err := NewTournamentSelection().SetProb(0.5)
	if err != ErrInvalidTournamentProb {
		t.Errorf("want ts.SetProb(0.5) = ErrInvalidTournamentProb, got %v", err)
	}
	err = NewTournamentSelection().SetProbRange(0.4, 0.6)
	if err != ErrInvalidTournamentProb {
		t.Errorf("want ts.SetProb(0.4, 0.6) = ErrInvalidTournamentProb, got %v", err)
	}
}
