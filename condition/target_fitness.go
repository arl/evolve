package condition

import (
	"fmt"

	"github.com/arl/evolve"
)

// TargetFitness is a termination condition that is met when at least one
// candidate in the population has reached or exceeded a given fitness score.
type TargetFitness[T any] struct {
	Fitness float64
	Natural bool
}

// IsSatisfied returns true if the specified fitness
// score has been reached or exceeded.
func (tf TargetFitness[T]) IsSatisfied(stats *evolve.PopulationStats[T]) bool {
	if tf.Natural {
		return stats.BestFitness >= tf.Fitness
	}
	return stats.BestFitness <= tf.Fitness
}

// String returns a string representation of this condition.
func (tf TargetFitness[T]) String() string {
	return fmt.Sprintf("Reached fitness of %f", tf.Fitness)
}
