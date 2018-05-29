package termination

import (
	"fmt"

	"github.com/aurelien-rainone/evolve/pkg/api"
)

// TargetFitness terminates evolution once at least one candidate in the
// population has reached or exceeded a given fitness score.
type TargetFitness struct {
	Fitness float64
	Natural bool
}

// IsSatisfied reports whether or not evolution should finish at the
// current point.
func (tc TargetFitness) IsSatisfied(data *api.PopulationData) bool {
	if tc.Natural {
		return data.BestFitness >= tc.Fitness
	}
	return data.BestFitness <= tc.Fitness
}

// String returns the termination condition representation as a string
func (tc TargetFitness) String() string {
	return fmt.Sprintf("Reached target fitness of %f", tc.Fitness)
}
