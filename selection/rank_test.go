package selection

import (
	"testing"
)

func TestRank(t *testing.T) {
	t.Run("natural", testFitnessBasedSelection(Rank[string](), fitnessBasedPopNatural, true))
	t.Run("non-natural", testFitnessBasedSelection(Rank[string](), fitnessBasedPopNonNatural, false))
}
