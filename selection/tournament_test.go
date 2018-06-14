package selection

import (
	"fmt"
	"testing"
)

func TestTournamentSelection(t *testing.T) {

	var tests = []struct {
		name    string
		natural bool
	}{
		{name: "natural", natural: true},
		{name: "non-natural", natural: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := NewTournament()
			errcheck(t, ts.SetProb(0.7))
			testRandomBasedSelection(t, ts, randomBasedPopNonNatural, tt.natural, 2,
				func(selected []interface{}) error {
					if len(selected) != 2 {
						return fmt.Errorf("want len(selected) == 2, got %v", len(selected))
					}
					return nil
				})
		})
	}
}

func TestTournamentSelectionSetProb(t *testing.T) {
	err := NewTournament().SetProb(0.5)
	if err != ErrInvalidTournamentProb {
		t.Errorf("want ts.SetProb(0.5) = ErrInvalidTournamentProb, got %v", err)
	}
	err = NewTournament().SetProbRange(0.4, 0.6)
	if err != ErrInvalidTournamentProb {
		t.Errorf("want ts.SetProb(0.4, 0.6) = ErrInvalidTournamentProb, got %v", err)
	}
}
