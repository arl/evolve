package selection

import (
	"math"
	"math/rand"

	"github.com/aurelien-rainone/evolve/framework"
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
// NOTE: It is an error to call this method with an empty or null population.
//
// population is the population from which to select.
// naturalFitnessScores indicates whether higher fitness values represent fitter
// individuals or not.
// selectionSize is the number of individual selections to make (not necessarily
// the number of distinct candidates to select, since the same individual may
// potentially be selected more than once).
//
// Returns a slice containing the selected candidates. Some individual
// candidates may potentially have been selected multiple times.
func (sel StochasticUniversalSampling) Select(
	population framework.EvaluatedPopulation,
	naturalFitnessScores bool,
	selectionSize int,
	rng *rand.Rand) []framework.Candidate {

	// Calculate the sum of all fitness values.
	var aggregateFitness float64
	for _, candidate := range population {
		aggregateFitness += adjustedFitness(candidate.Fitness(), naturalFitnessScores)
	}

	selection := make([]framework.Candidate, 0, selectionSize)

	// Pick a random offset between 0 and 1 as the starting point for
	// selection.
	var (
		startOffset           = rng.Float64()
		cumulativeExpectation float64
		index                 int
	)
	for _, candidate := range population {
		// Calculate the number of times this candidate is expected to
		// be selected on average and add it to the cumulative total
		// of expected frequencies.
		cumulativeExpectation += adjustedFitness(candidate.Fitness(),
			naturalFitnessScores) / aggregateFitness * float64(selectionSize)

		// If f is the expected frequency, the candidate will be selected at
		// least as often as floor(f) and at most as often as ceil(f). The
		// actual count depends on the random starting offset.
		for cumulativeExpectation > startOffset+float64(index) {
			selection = append(selection, candidate.Candidate())
			index++
		}
	}
	return selection
}

func (sel StochasticUniversalSampling) String() string {
	return "Stochastic Universal Sampling"
}

func adjustedFitness(rawFitness float64, naturalFitness bool) float64 {
	if naturalFitness {
		return rawFitness
	}
	// If standardised fitness is zero we have found the best possible
	// solution. The evolutionary algorithm should not be continuing
	// after finding it.
	if rawFitness == 0 {
		return math.MaxFloat64
	}
	return 1 / rawFitness
}
