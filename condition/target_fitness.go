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

// IsSatisfied returns true if the specified fitness
// score has been reached or exceeded.
func (tf TargetFitness) IsSatisfied(stats *evolve.PopulationStats) bool {
	if tf.Natural {
		return stats.BestFitness >= tf.Fitness
	}
	return stats.BestFitness <= tf.Fitness
}

// String returns a string representation of this condition.
func (tf TargetFitness) String() string {
	return fmt.Sprintf("Reached target fitness of %f", tf.Fitness)
}
