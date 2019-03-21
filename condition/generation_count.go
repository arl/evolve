package condition

import (
	"fmt"

	"github.com/arl/evolve"
)

// GenerationCount is a condtion that is met
// when a number of generation has passed.
type GenerationCount int

// IsSatisfied reports whether or not evolution
// should finish at the current point.
func (n GenerationCount) IsSatisfied(stats *evolve.PopulationStats) bool {
	return stats.GenNumber+1 >= int(n)
}

// String returns a string representation of this condition.
func (n GenerationCount) String() string {
	return fmt.Sprintf("Reached %d generations", n)
}
