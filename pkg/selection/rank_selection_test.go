package selection

import (
	"testing"
)

func TestRankSelectionNatural(t *testing.T) {
	testFitnessBasedSelection(t, Rank, fitnessBasedPopNatural, true)
}

func TestRankSelectionNonNatural(t *testing.T) {
	testFitnessBasedSelection(t, Rank, fitnessBasedPopNonNatural, false)
}
