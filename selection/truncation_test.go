package selection

import (
	"fmt"
	"testing"

	"github.com/arl/evolve/generator"
	"github.com/stretchr/testify/assert"
)

// Unit test for truncation selection strategy ensures the 2 best candidates are
// selected.
func testTruncationSelection(t *testing.T, tpop testPopulation, natural bool) {
	ts := &Truncation[string]{SelectionRatio: generator.Const(0.5)}
	testRandomBasedSelection(t, ts, tpop, natural, 2,
		func(selected []string) error {
			if len(selected) != 2 {
				return fmt.Errorf("want len(selected) == 2, got %v", len(selected))
			}

			sstr := []string{}
			for _, c := range selected {
				sstr = append(sstr, c)
			}

			assert.Contains(t, sstr, tpop[0].name, "best candidate not selected")
			assert.Contains(t, sstr, tpop[1].name, "second best candidate not selected")
			return nil
		})
}

func TestTruncationSelectionNatural(t *testing.T) {
	testTruncationSelection(t, randomBasedPopNatural, true)
}

func TestTruncationSelectionNonNatural(t *testing.T) {
	testTruncationSelection(t, randomBasedPopNonNatural, false)
}
