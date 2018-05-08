package selection

import (
	"fmt"
	"testing"
)

// Unit test for roulette selection strategy. We cannot easily test
// that the correct candidates are returned because of the random aspect
// of the selection, but we can at least make sure the right number of
// candidates are selected.
func testRouletteWheelSelection(t *testing.T, tpop testPopulation, natural bool) {
	for i := 0; i < 20; i++ {
		testRandomBasedSelection(t, RouletteWheelSelection, tpop, natural, 2,
			func(selected []interface{}) error {
				if len(selected) != 2 {
					return fmt.Errorf("want len(selected) == 2, got %v", len(selected))
				}
				return nil
			})
	}
}

func TestRouletteWheelSelectionNatural(t *testing.T) {
	testRouletteWheelSelection(t, randomBasedPopNatural, true)
}

func TestRouletteWheelSelectionNonNatural(t *testing.T) {
	testRouletteWheelSelection(t, randomBasedPopNonNatural, false)
}

func TestRouletteWheelSelectionNaturalPerfect(t *testing.T) {
	var testPop = testPopulation{
		{name: "Gary", fitness: 0.0},
		{name: "Mary", fitness: 8.4},
		{name: "John", fitness: 9.1},
		{name: "Steve", fitness: 10.0},
	}

	testRouletteWheelSelection(t, testPop, true)
}
