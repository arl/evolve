package selection

import (
	"math"
	"math/rand"
	"sort"

	"github.com/aurelien-rainone/evolve/framework"
)

// RouletteWheelSelection implements selection of N candidates from a population
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
type RouletteWheelSelection struct{}

// Select selects the required number of candidates from the population with the
// probability of selecting any particular candidate being proportional to that
// candidate's fitness score.  Selection is with replacement (the same candidate
// may be selected multiple times).
//
// - naturalFitnessScores should be true if higher fitness scores indicate
// fitter individuals, false if lower fitness scores indicate fitter
// individuals.
// - selectionSize is the number of selections to make.
func (rws *RouletteWheelSelection) Select(population []*framework.EvaluatedCandidate, naturalFitnessScores bool, selectionSize int, rng *rand.Rand) []framework.Candidate {

	// Record the cumulative fitness scores. It doesn't matter whether the
	// population is sorted or not. We will use these cumulative scores to
	// work out an index into the population. The cumulative array itself is
	// implicitly sorted since each element must be greater than the
	// previous one. The numerical difference between an element and the
	// previous one is directly proportional to the probability of the
	// corresponding candidate in the population being selected.
	cumulativeFitnesses := make([]float64, len(population))
	cumulativeFitnesses[0] = adjustedFitness(population[0].Fitness(), naturalFitnessScores)
	for i := 1; i < len(population); i++ {
		fitness := adjustedFitness(population[i].Fitness(), naturalFitnessScores)
		cumulativeFitnesses[i] = cumulativeFitnesses[i-1] + fitness
	}

	selection := make([]framework.Candidate, selectionSize)
	for i := 0; i < selectionSize; i++ {
		randomFitness := rng.Float64() * cumulativeFitnesses[len(cumulativeFitnesses)-1]
		index := sort.SearchFloat64s(cumulativeFitnesses, randomFitness)
		if index < 0 {
			// Convert negative insertion point to array index.
			index = abs(index + 1)
		}
		selection[i] = population[index].Candidate()
	}
	return selection
}

func (rws *RouletteWheelSelection) String() string {
	return "Roulette Wheel Selection"
}

func adjustedFitness(rawFitness float64, naturalFitness bool) float64 {
	if naturalFitness {
		return rawFitness
	}
	// If standardised fitness is zero we have found the best possible
	// solution.  The evolutionary algorithm should not be continuing
	// after finding it.
	if rawFitness == 0 {
		return math.MaxFloat64
	}
	return 1 / rawFitness
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}
