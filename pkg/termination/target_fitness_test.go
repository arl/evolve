package termination

import (
	"testing"

	"github.com/aurelien-rainone/evolve/pkg/api"
)

// Unit test for termination condition that checks the best fitness attained so far
// against a pre-determined target.
func TestTargetFitness(t *testing.T) {

	t.Run("natural fitness", func(t *testing.T) {
		cond := NewTargetFitness(10.0, true)

		if cond.ShouldTerminate(
			&api.PopulationData{struct{}{}, 5.0, 4.0, 0, true, 2, 0, 0, 100}) {
			t.Errorf("should not terminate before target fitness has been reached")
		}

		if !cond.ShouldTerminate(
			&api.PopulationData{struct{}{}, 10.0, 8.0, 0, true, 2, 0, 0, 100}) {
			t.Errorf("should terminate once target fitness has been reached")
		}
	})

	t.Run("non-natural fitness", func(t *testing.T) {
		cond := NewTargetFitness(1.0, false)

		if cond.ShouldTerminate(
			&api.PopulationData{struct{}{}, 5.0, 4.0, 0, true, 2, 0, 0, 100}) {
			t.Errorf("should not terminate before target fitness has been reached")
		}

		if !cond.ShouldTerminate(
			&api.PopulationData{struct{}{}, 1.0, 3.1, 0, true, 2, 0, 0, 100}) {
			t.Errorf("should terminate once target fitness has been reached")
		}
	})
}
