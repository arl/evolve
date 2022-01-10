package selection

import (
	"fmt"
	"testing"

	"github.com/arl/evolve/generator"
)

func TestTournamentSelection(t *testing.T) {
	tests := []struct {
		name    string
		natural bool
	}{
		{name: "natural", natural: true},
		{name: "non-natural", natural: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := &Tournament[string]{Probability: generator.Const(0.7)}
			testRandomBasedSelection(t, ts, randomBasedPopNonNatural, tt.natural, 2,
				func(selected []string) error {
					if len(selected) != 2 {
						return fmt.Errorf("want len(selected) == 2, got %v", len(selected))
					}
					return nil
				})
		})
	}
}
