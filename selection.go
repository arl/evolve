package evolve

import (
	"fmt"
	"math/rand"
)

// Selection is the interface that wraps the Select method.
type Selection[T any] interface {
	fmt.Stringer

	// Select selects a given number of candidates from a population.
	//
	// The population must be sorted by descending fitness (i.e. pop[0] is the
	// fittest). A natural fitness means that the higher fitness, the better. n
	// is the number of candidates to select and return. Note that since a same
	// individual can be selected more than once, they may not be all distinct.
	Select(pop Population[T], natural bool, n int, rng *rand.Rand) []T
}
