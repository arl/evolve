package selection

import (
	"fmt"
	"testing"

	"github.com/arl/evolve/generator"
)

func TestTournamentSelection(t *testing.T) {
	tournament := &Tournament[string]{Probability: generator.Const(0.7)}
	check := func(s []string) error {
		if len(s) != 2 {
			return fmt.Errorf("got %d selected elements, want 2", len(s))
		}
		return nil
	}

	t.Run("natural", testRandomBasedSelection(tournament, randomBasedPopNatural, true, 2, check))
	t.Run("non-natural", testRandomBasedSelection(tournament, randomBasedPopNonNatural, false, 2, check))
}
