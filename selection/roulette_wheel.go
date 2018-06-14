package selection

import (
	"math/rand"
	"sort"

	"github.com/arl/evolve"
)

// RouletteWheel implements selection of N candidates from a population
// by selecting N candidates at random where the probability of each candidate
// getting selected is proportional to its fitness score.
//
// This is analogous to each candidate being assigned an area on a roulette
// wheel proportionate to its fitness and the wheel being spun N times.
// Candidates may be selected more than once.  In some instances, particularly
// with small population sizes, the randomness of selection may result in
// excessively high occurrences of particular candidates. If this is a problem,
// StochasticUniversalSampling provides an alternative fitness-proportionate
// strategy for selection.
var RouletteWheel = rouletteWheel{}

type rouletteWheel struct{}

// Select selects the required number of candidates from the population with the
// probability of selecting any particular candidate being proportional to that
// candidate's fitness score.  Selection is with replacement (the same candidate
// may be selected multiple times).
//
// naturalFitnessScores should be true if higher fitness scores indicate fitter
// individuals, false if lower fitness scores indicate fitter individuals.
// selectionSize is the number of selections to make.
func (rouletteWheel) Select(
	pop evolve.Population,
	natural bool,
	size int,
	rng *rand.Rand) []interface{} {

	// Record the cumulative fitness scores. It doesn't matter whether the
	// population is sorted or not. We will use these cumulative scores to
	// work out an index into the population. The cumulative array itself is
	// implicitly sorted since each element must be greater than the
	// previous one. The numerical difference between an element and the
	// previous one is directly proportional to the probability of the
	// corresponding candidate in the population being selected.
	cumfitness := make([]float64, len(pop))
	cumfitness[0] = adjustedFitness(pop[0].Fitness, natural)
	for i := 1; i < len(pop); i++ {
		fitness := adjustedFitness(pop[i].Fitness, natural)
		cumfitness[i] = cumfitness[i-1] + fitness
	}

	sel := make([]interface{}, size)
	for i := 0; i < size; i++ {
		rand := rng.Float64() * cumfitness[len(cumfitness)-1]
		j := sort.SearchFloat64s(cumfitness, rand)
		if j < 0 {
			// Convert negative insertion point to array index.
			j = abs(j + 1)
		}
		sel[i] = pop[j].Candidate
	}
	return sel
}

func (rouletteWheel) String() string { return "Roulette Wheel Selection" }

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}
