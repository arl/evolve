package selection

import (
	"fmt"
	"testing"

	"github.com/arl/evolve/generator"
)

func TestTruncationSelection(t *testing.T) {
	test := func(tpop testPopulation, natural bool) func(*testing.T) {
		check := func(s []string) error {
			if len(s) != 2 {
				return fmt.Errorf("got %d selected elements, want 2", len(s))
			}

			// Have we selected the 2 fittest individuals?
			i0, i1 := false, false
			for _, c := range s {
				i0 = i0 || (c == tpop[0].name)
				i1 = i1 || (c == tpop[1].name)
			}
			if !i0 {
				t.Errorf("best candidate %q not selected", tpop[0].name)
			}
			if !i1 {
				t.Errorf("second best candidate %q not selected", tpop[1].name)
			}
			return nil
		}

		truncation := &Truncation[string]{SelectionRatio: generator.Const(0.5)}
		return testRandomBasedSelection(truncation, tpop, natural, 2, check)
	}

	t.Run("natural", test(randomBasedPopNatural, true))
	t.Run("non-natural", test(randomBasedPopNonNatural, false))
}
