package condition

import (
	"testing"

	"github.com/arl/evolve"
)

// Unit test for termination condition that checks the best fitness attained so far
// against a pre-determined target.
func TestTargetFitness(t *testing.T) {
	t.Run("natural fitness", func(t *testing.T) {
		cond := TargetFitness[any]{Fitness: 10.0, Natural: true}
		stats := &evolve.PopulationStats[any]{Natural: true}

		stats.BestFitness = 5.0
		if cond.IsSatisfied(stats) {
			t.Errorf("fitness = %v, termination condition should not be satisfied", stats.BestFitness)
		}

		stats.BestFitness = 10.0
		if !cond.IsSatisfied(stats) {
			t.Errorf("fitness = %v, termination condition should be satisfied", stats.BestFitness)
		}
	})

	t.Run("non-natural fitness", func(t *testing.T) {
		cond := TargetFitness[any]{Fitness: 1.0, Natural: false}
		stats := &evolve.PopulationStats[any]{Natural: false}

		stats.BestFitness = 5.0
		if cond.IsSatisfied(stats) {
			t.Errorf("fitness = %v, termination condition should not be satisfied", stats.BestFitness)
		}

		stats.BestFitness = 1.0
		if !cond.IsSatisfied(stats) {
			t.Errorf("fitness = %v, termination condition should be satisfied", stats.BestFitness)
		}
	})
}
