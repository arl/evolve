package selection

import (
	"math/rand"
	"sort"

	"github.com/arl/evolve"
)

// RouletteWheel is a selection strategy where the probability of picking a
// candidate is proportional to its fitness score.
//
// This is analogous to each candidate being assigned an area on a roulette
// wheel proportionate to its fitness and the wheel being spun N times.
// Candidates may be selected more than once. In some instances, particularly
// with small population sizes, the randomness of selection may result in
// excessively high occurrences of particular candidates. If this is a problem,
// SUS provides an alternative fitness-proportionate strategy for selection.
type RouletteWheel[T any] struct{}

// Select selects the required number of candidates from the population with the
// probability of selecting any particular candidate being proportional to that
// candidate's fitness score. Selection is with replacement (the same candidate
// may be selected multiple times).
//
// naturalFitnessScores should be true if higher fitness scores indicate fitter
// individuals, false if lower fitness scores indicate fitter individuals.
// selectionSize is the number of selections to make.
func (RouletteWheel[T]) Select(pop *evolve.Population[T], natural bool, n int, rng *rand.Rand) []T {
	// Record the cumulative fitness scores. It doesn't matter whether the
	// population is sorted or not. We will use these cumulative scores to
	// work out an index into the population. The cumulative array itself is
	// implicitly sorted since each element must be greater than the
	// previous one. The numerical difference between an element and the
	// previous one is directly proportional to the probability of the
	// corresponding candidate in the population being selected.
	cumfitness := make([]float64, pop.Len())
	cumfitness[0] = adjustedFitness(pop.Fitness[0], natural)
	for i := 1; i < pop.Len(); i++ {
		fitness := adjustedFitness(pop.Fitness[i], natural)
		cumfitness[i] = cumfitness[i-1] + fitness
	}

	sel := make([]T, n)
	for i := 0; i < n; i++ {
		rand := rng.Float64() * cumfitness[len(cumfitness)-1]
		j := sort.SearchFloat64s(cumfitness, rand)
		if j < 0 {
			// Convert negative insertion point to array index.
			j = abs(j + 1)
		}
		sel[i] = pop.Candidates[j]
	}
	return sel
}

func (RouletteWheel[T]) String() string { return "Roulette Wheel Selection" }

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}
