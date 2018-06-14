package selection

import (
	"math/rand"

	"github.com/arl/evolve"
)

// Identity is a selection strategy that returns identical candidates
type Identity struct{}

// Select selects the specified number of candidates from the population.
func (Identity) Select(
	pop evolve.Population,
	natural bool,
	size int,
	rng *rand.Rand) []interface{} {

	sel := make([]interface{}, size)
	for i := 0; i < size; i++ {
		sel[i] = pop[i].Candidate
	}
	return sel
}

func (Identity) String() string { return "Identity Selection" }
