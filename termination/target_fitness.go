package termination

import (
	"fmt"

	"github.com/aurelien-rainone/evolve/framework"
)

// TargetFitness terminates evolution once at least one candidate in the
// population has equalled or bettered a pre-determined fitness score.
type TargetFitness struct {
	targetFitness float64
	natural       bool
}

// NewTargetFitness returns a TargetFitness termination condition.
//
// targetFitness represents the fitness score that must be achieved by at least
// one individual in the population in order for this condition to be satisfied.
// natural indicates whether fitness scores are natural or non-natural. If
// fitness is natural, the condition will be satisfied if any individual has a
// fitness that is greater than or equal to the target fitness. If fitness is
// non-natural, the condition will be satisfied in any individual has a fitness
// that is less than or equal to the target fitness.
func NewTargetFitness(targetFitness float64, natural bool) *TargetFitness {
	return &TargetFitness{
		targetFitness: targetFitness,
		natural:       natural,
	}
}

// ShouldTerminate reports whether or not evolution should finish at the
// current point.
//
// populationData is the information about the current state of evolution.
// This may be used to determine whether evolution should continue or not.
func (tc *TargetFitness) ShouldTerminate(populationData *framework.PopulationData) bool {
	if tc.natural {
		// If we're using "natural" fitness scores, higher values are better.
		return populationData.BestCandidateFitness() >= tc.targetFitness
	}
	// If we're using "non-natural" fitness scores, lower values are better.
	return populationData.BestCandidateFitness() <= tc.targetFitness
}

// String returns the termination condition representation as a string
func (tc *TargetFitness) String() string {
	return fmt.Sprintf("Reached target fitness of %f", tc.targetFitness)
}
