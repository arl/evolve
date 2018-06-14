package selection

import (
	"math"
	"math/rand"

	"github.com/arl/evolve"
)

// StochasticUniversalSampling is an alternative to RouletteWheelSelection as a
// fitness-proportionate selection strategy. Ensures that the frequency of
// selection for each candidate is consistent with its expected frequency of
// selection.
type StochasticUniversalSampling struct{}

// Select selects the specified number of candidates from the population.
//
// Implementations may assume that the population is sorted in descending
// order according to fitness (so the fittest individual is the first item
// in the list).
// NOTE: It is an error to call this method with an empty or nil population.
//
// pop is the population from which to select.
// natural indicates whether higher fitness values represent fitter individuals
// or not.
// size is the number of individual selections to make (not necessarily the
// number of distinct candidates to select, since the same individual may
// potentially be selected more than once).
//
// Returns a slice containing the selected candidates. Some individual
// candidates may potentially have been selected multiple times.
func (StochasticUniversalSampling) Select(
	pop evolve.Population,
	natural bool,
	size int,
	rng *rand.Rand) []interface{} {

	// Calculate the sum of all fitness values.
	var sum float64
	for _, cand := range pop {
		sum += adjustedFitness(cand.Fitness, natural)
	}

	sel := make([]interface{}, 0, size)

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
		expect += adjustedFitness(cand.Fitness, natural) / sum * float64(size)

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

func (StochasticUniversalSampling) String() string { return "Stochastic Universal Sampling" }

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
