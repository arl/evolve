package selection

import (
	"math/rand"

	"github.com/arl/evolve"
)

// Identity is a selection strategy that returns identical candidates.
type Identity[T any] struct{}

// Select selects the n fittest candidates from the poipulation.
func (Identity[T]) Select(pop evolve.Population[T], natural bool, n int, rng *rand.Rand) []T {
	sel := make([]T, n)
	for i := 0; i < n; i++ {
		sel[i] = pop.Candidates[i]
	}
	return sel
}

func (Identity[T]) String() string { return "Identity Selection" }
