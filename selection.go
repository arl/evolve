package evolve

import (
	"fmt"
	"math/rand"
)

// Selection is the interface that wraps the Select method.
//
// Select implements "natural" selection.
type Selection interface {
	fmt.Stringer

	// Select selects the specified number of candidates from the population.
	//
	// - pop must be sorted by descending fitness, i.e the fittest individual of the
	// population should be pop[0].
	// - natural indicates fitter individuals have fitness scores.
	// - size is the number of individual selections to perform (not necessarily the
	// number of distinct candidates to select, since the same individual may
	// potentially be selected more than once).
	//
	// Returns the selected candidates.
	Select(pop Population, natural bool, size int, rng *rand.Rand) []interface{}
}
