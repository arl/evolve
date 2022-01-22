package selection

import (
	"math"
	"math/rand"

	"github.com/arl/evolve"
)

// SUS is stochastic univeral sampling selection, an alternative to roulette
// wheel as a fitness-proportionate selection strategy. Ensures that the
// frequency of selection for each candidate is consistent with its expected
// frequency of selection.
type SUS[T any] struct{}

// Select selects a given number of candidates from a population.
func (SUS[T]) Select(pop evolve.Population[T], natural bool, n int, rng *rand.Rand) []T {
	// Calculate the sum of all fitness values.
	var sum float64
	for _, cand := range pop {
		sum += adjustedFitness(cand.Fitness, natural)
	}

	sel := make([]T, 0, n)

	// Pick a random offset between 0 and 1 as the starting point for
	// selection.
	var (
		off    = rng.Float64()
		expect float64
		i      int
	)
	for _, cand := range pop {
		// Calculate the number of times this candidate is expected to
		// be selected on average and add it to the cumulative total
		// of expected frequencies.
		expect += adjustedFitness(cand.Fitness, natural) / sum * float64(n)

		// If f is the expected frequency, the candidate will be selected at
		// least as often as floor(f) and at most as often as ceil(f). The
		// actual count depends on the random starting offset.
		for expect > off+float64(i) {
			sel = append(sel, cand.Candidate)
			i++
		}
	}
	return sel
}

func (SUS[T]) String() string { return "Stochastic Universal Sampling" }

func adjustedFitness(fitness float64, natural bool) float64 {
	if natural {
		return fitness
	}
	// If standardised fitness is zero we have found the best possible
	// solution. The evolutionary algorithm should not be continuing
	// after finding it.
	if fitness == 0 {
		return math.MaxFloat64
	}
	return 1 / fitness
}
