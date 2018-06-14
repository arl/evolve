package condition

import (
	"fmt"

	"github.com/arl/evolve"
)

// GenerationCount is a condtion that is met when a number of generation has
// passed.
type GenerationCount int

// IsSatisfied reports whether or not evolution should finish at the
// current point.
func (num GenerationCount) IsSatisfied(popdata *evolve.PopulationData) bool {
	return popdata.GenNumber+1 >= int(num)
}

// String returns a string representation of this condition.
func (num GenerationCount) String() string {
	return fmt.Sprintf("Reached %d generations", num)
}
