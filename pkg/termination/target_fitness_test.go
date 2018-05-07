package termination

import (
	"testing"

	"github.com/aurelien-rainone/evolve/pkg/api"
)

// Unit test for termination condition that checks the best fitness attained so far
// against a pre-determined target.
func TestTargetFitness(t *testing.T) {

	t.Run("natural fitness", func(t *testing.T) {
		cond := TargetFitness{Fitness: 10.0, Natural: true}
		popdata := &api.PopulationData{Natural: true}

		popdata.BestFitness = 5.0
		if cond.ShouldTerminate(popdata) {
			t.Errorf("should not terminate before target fitness has been reached")
		}

		popdata.BestFitness = 10.0
		if !cond.ShouldTerminate(popdata) {
			t.Errorf("should terminate once target fitness has been reached")
		}
	})

	t.Run("non-natural fitness", func(t *testing.T) {
		cond := TargetFitness{Fitness: 1.0, Natural: false}
		popdata := &api.PopulationData{Natural: true}

		popdata.BestFitness = 5.0
		if cond.ShouldTerminate(popdata) {
			t.Errorf("should not terminate before target fitness has been reached")
		}

		popdata.BestFitness = 1.0
		if !cond.ShouldTerminate(popdata) {
			t.Errorf("should terminate once target fitness has been reached")
		}
	})
}
