package selection

import (
	"fmt"
	"testing"
)

func testRouletteWheelSelection(tpop testPopulation, natural bool) func(t *testing.T) {
	return func(t *testing.T) {
		// We can't easily test that the correct candidates are returned because of
		// the random aspect of the selection, but we can at least make sure the
		// right number of candidates are selected.
		check := func(s []string) error {
			if len(s) != 2 {
				return fmt.Errorf("got %d selected elements, want 2", len(s))
			}
			return nil
		}

		for i := 0; i < 20; i++ {
			t.Run(fmt.Sprintf("run_%d", i), testRandomBasedSelection(RouletteWheel[string]{}, tpop, natural, 2, check))
		}
	}
}

func TestRouletteWheelSelection(t *testing.T) {
	t.Run("natural", testRouletteWheelSelection(randomBasedPopNatural, true))
	t.Run("non-natural", testRouletteWheelSelection(randomBasedPopNonNatural, false))

	perfect := testPopulation{
		{name: "Gary", fitness: 0.0},
		{name: "Mary", fitness: 8.4},
		{name: "John", fitness: 9.1},
		{name: "Steve", fitness: 10.0},
	}
	t.Run("natural-perfect", testRouletteWheelSelection(perfect, true))
}
