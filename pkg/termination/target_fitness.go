package termination

import (
	"fmt"

	"github.com/aurelien-rainone/evolve/pkg/api"
)

// TargetFitness terminates evolution once at least one candidate in the
// population has equalled or bettered a pre-determined fitness score.
type TargetFitness struct {
	Fitness float64
	Natural bool
}

// ShouldTerminate reports whether or not evolution should finish at the
// current point.
//
// populationData is the information about the current state of evolution.
// This may be used to determine whether evolution should continue or not.
func (tc TargetFitness) ShouldTerminate(populationData *api.PopulationData) bool {
	if tc.Natural {
		// If we're using "natural" fitness scores, higher values are better.
		return populationData.BestFitness >= tc.Fitness
	}
	// If we're using "non-natural" fitness scores, lower values are better.
	return populationData.BestFitness <= tc.Fitness
}

// String returns the termination condition representation as a string
func (tc TargetFitness) String() string {
	return fmt.Sprintf("Reached target fitness of %f", tc.Fitness)
}
