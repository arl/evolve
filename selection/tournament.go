package selection

import (
	"math/rand"

	"github.com/arl/evolve"
	"github.com/arl/evolve/generator"
)

// Tournament is a selection strategy that picks a pair of candidates at random
// and then selects the fitter of the two candidates with probability p, where p
// is the selection probability (therefore the probability of the less fit
// candidate being selected is 1 - p).
type Tournament[T any] struct {
	Probability generator.Float
}

// Select selects the specified number of candidates from the population.
func (ts *Tournament[T]) Select(pop evolve.Population[T], natural bool, size int, rng *rand.Rand) []T {
	sel := make([]T, size)
	for i := 0; i < size; i++ {
		// Pick two candidates at random.
		cand1 := pop[rng.Intn(len(pop))]
		cand2 := pop[rng.Intn(len(pop))]

		// Get probability for this selection.
		if natural && rng.Float64() < ts.Probability.Next() { // Select the fitter candidate.
			if cand2.Fitness > cand1.Fitness {
				sel[i] = cand2.Candidate
			} else {
				sel[i] = cand1.Candidate
			}
		} else { // Select the less fit candidate.
			if cand2.Fitness > cand1.Fitness {
				sel[i] = cand1.Candidate
			} else {
				sel[i] = cand2.Candidate
			}
		}
	}
	return sel
}

func (ts *Tournament[T]) String() string {
	return "Tournament Selection"
}
