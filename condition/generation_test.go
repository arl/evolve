package condition

import (
	"testing"

	"github.com/arl/evolve"
)

func TestGenerationCount(t *testing.T) {
	cond := GenerationCount[any](5)
	stats := &evolve.PopulationStats[any]{}

	stats.Generation = 3
	if cond.IsSatisfied(stats) {
		t.Errorf("generation = %v, termination condition should not be satisfied", stats.Generation)
	}

	stats.Generation = 4
	if !cond.IsSatisfied(stats) {
		t.Errorf("generation = %v, termination condition should be satisfied", stats.Generation)
	}
}
