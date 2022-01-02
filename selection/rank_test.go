package selection

import (
	"testing"
)

func TestRankSelectionNatural(t *testing.T) {
	testFitnessBasedSelection(t, Rank[string](), fitnessBasedPopNatural, true)
}

func TestRankSelectionNonNatural(t *testing.T) {
	testFitnessBasedSelection(t, Rank[string](), fitnessBasedPopNonNatural, false)
}
