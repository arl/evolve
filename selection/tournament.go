package selection

import (
	"math/rand"

	"github.com/arl/evolve"
	"github.com/arl/evolve/generator"
)

// Tournament is a selection strategy that picks a pair of candidates at random
// and then selects the fitter of the two candidates with probability p, where p
// is the selection probability (i.e the probability of the less fit candidate
// being selected is 1 - p).
type Tournament[T any] struct {
	Probability generator.Float
}

// Select selects the specified number of candidates from the population.
func (ts *Tournament[T]) Select(pop *evolve.Population[T], natural bool, n int, rng *rand.Rand) []T {
	sel := make([]T, n)
	selidx := 0
	for i := 0; i < n; i++ {
		// Pick two candidates at random.
		idx1 := rng.Intn(pop.Len())
		idx2 := rng.Intn(pop.Len())

		// Get probability for this selection.
		if natural && rng.Float64() < ts.Probability.Next() {
			// Select the fitter candidate.
			if pop.Fitness[idx2] > pop.Fitness[idx1] {
				selidx = idx2
			} else {
				selidx = idx1
			}
		} else {
			// Select the less fit candidate.
			if pop.Fitness[idx2] > pop.Fitness[idx1] {
				selidx = idx1
			} else {
				selidx = idx2
			}
		}

		sel[i] = pop.Candidates[selidx]
	}
	return sel
}

func (ts *Tournament[T]) String() string {
	return "Tournament Selection"
}
