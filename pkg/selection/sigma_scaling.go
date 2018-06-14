package selection

import (
	"math/rand"

	"github.com/arl/evolve"
)

type sigmaScaling struct{ selector evolve.Selection }

// NewSigmaScaling creates a sigma-scaled selection strategy. This is an
// alternative to straightforward fitness-proportionate
// selection such as that offered by RouletteWheelSelection and
// StochasticUniversalSampling. It uses the mean population fitness and fitness
// standard deviation to adjust individual fitness scores.
//
// selector is the proportionate selector that will be delegated to after
// fitness scores have been adjusted using sigma scaling.
//
// Early on in an evolutionary algorithm this helps to avoid premature
// convergence caused by the dominance of one or two relatively fit candidates
// in a population of mostly unfit individuals. It also helps to amplify minor
// fitness differences in a more mature population where the rate of improvement
// has slowed.
func NewSigmaScaling(selector evolve.Selection) evolve.Selection {
	return &sigmaScaling{selector: selector}
}

// SigmaScaling is the default sigma scaling selection strategy. It uses
// StochasticUniversalSampling as its selector.
var SigmaScaling = NewSigmaScaling(StochasticUniversalSampling{})

// Select selects the specified number of candidates from the population.
//
// Implementations may assume that the population is sorted in descending
// order according to fitness (so the fittest individual is the first item
// in the list).
// NOTE: It is an error to call this method with an empty or null population.
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
func (sel *sigmaScaling) Select(
	pop evolve.Population,
	natural bool,
	size int,
	rng *rand.Rand) []interface{} {

	stats := evolve.NewDataset(len(pop))
	for _, cand := range pop {
		stats.AddValue(cand.Fitness)
	}

	scaledPop := make(evolve.Population, len(pop))
	for i, cand := range pop {
		scaledFitness := sigmaScaledFitness(cand.Fitness,
			stats.ArithmeticMean(),
			stats.StandardDeviation())

		scaledPop[i] = &evolve.Individual{
			Candidate: cand.Candidate,
			Fitness:   scaledFitness,
		}
	}
	return sel.selector.Select(scaledPop, natural, size, rng)
}

func (sigmaScaling) String() string { return "Sigma Scaling" }

func sigmaScaledFitness(candidateFitness, populationMeanFitness, fitnessStandardDeviation float64) float64 {
	if fitnessStandardDeviation == 0 {
		return 1
	}
	scaledFitness := 1 + (candidateFitness-populationMeanFitness)/(2*fitnessStandardDeviation)
	// Don't allow negative expected frequencies, use an arbitrary low
	// but still positive frequency of 1 time in 10 for extremely unfit
	// individuals (relative to the remainder of the population).
	if scaledFitness > 0 {
		return scaledFitness
	}
	return 0.1
}
