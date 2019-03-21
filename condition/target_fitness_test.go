package condition

import (
	"testing"

	"github.com/arl/evolve"
)

// Unit test for termination condition that checks the best fitness attained so far
// against a pre-determined target.
func TestTargetFitness(t *testing.T) {

	t.Run("natural fitness", func(t *testing.T) {
		cond := TargetFitness{Fitness: 10.0, Natural: true}
		stats := &evolve.PopulationStats{Natural: true}

		stats.BestFitness = 5.0
		if cond.IsSatisfied(stats) {
			t.Errorf("should not terminate before target fitness has been reached")
		}

		stats.BestFitness = 10.0
		if !cond.IsSatisfied(stats) {
			t.Errorf("should terminate once target fitness has been reached")
		}
	})

	t.Run("non-natural fitness", func(t *testing.T) {
		cond := TargetFitness{Fitness: 1.0, Natural: false}
		stats := &evolve.PopulationStats{Natural: false}

		stats.BestFitness = 5.0
		if cond.IsSatisfied(stats) {
			t.Errorf("should not terminate before target fitness has been reached")
		}

		stats.BestFitness = 1.0
		if !cond.IsSatisfied(stats) {
			t.Errorf("should terminate once target fitness has been reached")
		}
	})
}
