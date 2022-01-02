package condition

import (
	"fmt"

	"github.com/arl/evolve"
)

// GenerationCount is a condtion that is met
// when a number of generation has passed.
type GenerationCount[T any] int

// IsSatisfied reports whether or not evolution
// should finish at the current point.
func (n GenerationCount[T]) IsSatisfied(stats *evolve.PopulationStats[T]) bool {
	return stats.GenNumber+1 >= int(n)
}

// String returns a string representation of this condition.
func (n GenerationCount[T]) String() string {
	return fmt.Sprintf("Reached %d generations", n)
}
