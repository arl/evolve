package condition

import (
	"testing"

	"github.com/arl/evolve"
)

func TestGenerationCount(t *testing.T) {
	cond := GenerationCount[any](5)
	stats := &evolve.PopulationStats[any]{}

	stats.GenNumber = 3
	if cond.IsSatisfied(stats) {
		t.Errorf("should not terminate after 4th generation")
	}

	stats.GenNumber = 4
	if !cond.IsSatisfied(stats) {
		t.Errorf("should terminate after 5th generation")
	}
}
