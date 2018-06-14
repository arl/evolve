package condition

import (
	"fmt"

	"github.com/arl/evolve"
)

// TargetFitness is a condition that is met when at least one candidate in the
// population has reached or exceeded a given fitness score.
type TargetFitness struct {
	Fitness float64
	Natural bool
}

// IsSatisfied returns true if the time duration has elapsed.
func (tc TargetFitness) IsSatisfied(data *evolve.PopulationData) bool {
	if tc.Natural {
		return data.BestFitness >= tc.Fitness
	}
	return data.BestFitness <= tc.Fitness
}

// String returns a string representation of this condition.
func (tc TargetFitness) String() string {
	return fmt.Sprintf("Reached target fitness of %f", tc.Fitness)
}
