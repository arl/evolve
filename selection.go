package evolve

import (
	"fmt"
	"math/rand"
)

// Selection is the interface that wraps the Select method.
type Selection[T any] interface {
	fmt.Stringer

	// Select selects a given number of candidates from a population. pop must
	// be sorted by descending fitness (i.e pop[0] is the fittest), natural
	// indicates whether candidates have natural fitness (if true, the higher
	// the better). Size if the number of selections to perform (not necessarily
	// the number of distinct candidates to select, since the same individual
	// may potentially be selected more than once).
	Select(pop Population[T], natural bool, size int, rng *rand.Rand) []T
}
