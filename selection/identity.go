package selection

import (
	"math/rand"

	"github.com/arl/evolve"
)

// Identity is a selection strategy that returns identical candidates
type Identity[T any] struct{}

// Select selects the specified number of candidates from the population.
func (Identity[T]) Select(pop evolve.Population[T], natural bool, size int, rng *rand.Rand) []T {
	sel := make([]T, size)
	for i := 0; i < size; i++ {
		sel[i] = pop[i].Candidate
	}
	return sel
}

func (Identity[T]) String() string { return "Identity Selection" }
