package selection

import (
	"constraints"
	"math"
	"math/rand"

	"github.com/arl/evolve"
	"github.com/arl/evolve/generator"
)

// Truncation implements the selection of n candidates from a population by
// simply selecting the n candidates with the highest fitness scores (the rest
// is discarded). The same candidate is never selected more than once.
type Truncation[T any] struct {
	// TODO: document
	SelectionRatio generator.Float
}

// Select selects the fittest candidates. If the SelectionRatio results in fewer
// selected candidates than required, then these candidates are selected
// multiple times to make up the shortfall.
//
// pop is the population of evolved and evaluated candidates from which to
// select. natural indicates whether higher fitness values represent fitter
// individuals or not. size is the number of candidates to select from the
// evolved population.
//
// Returns the selected candidates.
func (ts *Truncation[T]) Select(pop evolve.Population[T], natural bool, size int, rng *rand.Rand) []T {
	sel := make([]T, 0, size)

	// get a random value to decide wether to select the fitter individual
	// or the weaker one.
	eligible := int(math.Round(ts.SelectionRatio.Next() * float64(len(pop))))
	if eligible > size {
		eligible = size
	}

	for {
		count := min(eligible, size-len(sel))
		for i := 0; i < count; i++ {
			sel = append(sel, pop[i].Candidate)
		}
		if len(sel) >= size {
			break
		}
	}
	return sel
}

func (ts *Truncation[T]) String() string {
	return "Truncation Selection"
}

func min[T constraints.Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}
