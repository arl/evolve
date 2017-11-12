package selection

import (
	"fmt"
	"math/rand"

	"github.com/aurelien-rainone/evolve/framework"
)

// SigmaScaling is an alternative to straightforward fitness-proportionate
// selection such as that offered by RouletteWheelSelection and
// StochasticUniversalSampling. Uses the mean population fitness and fitness
// standard deviation to adjust individual fitness scores.
//
// Early on in an evolutionary algorithm this helps to avoid premature
// convergence caused by the dominance of one or two relatively fit candidates
// in a population of mostly unfit individuals. It also helps to amplify minor
// fitness differences in a more mature population where the rate of improvement
// has slowed.
type SigmaScaling struct {
	delegate framework.SelectionStrategy
}

// NewSigmaScaling creates a sigma-scaled selection strategy.
//
// delegate is the proprtionate selector that will be delegated to after after
// fitness scores have been adjusted using sigma scaling. The delegate
// parameter. It may be nil, in which case the SigmaScaling is created with a
// default delegate using stochastic universal sampling
func NewSigmaScaling(delegate framework.SelectionStrategy) *SigmaScaling {
	if delegate == nil {
		delegate = StochasticUniversalSampling{}
	}
	return &SigmaScaling{delegate: delegate}
}

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
func (sel *SigmaScaling) Select(
	population framework.EvaluatedPopulation,
	naturalFitnessScores bool,
	selectionSize int,
	rng *rand.Rand) []framework.Candidate {

	statistics := framework.NewDataSet(framework.WithInitialCapacity(len(population)))
	for _, candidate := range population {
		statistics.AddValue(candidate.Fitness())
	}

	scaledPopulation := make(framework.EvaluatedPopulation, len(population))
	var err error
	for i, candidate := range population {
		scaledFitness := sigmaScaledFitness(candidate.Fitness(),
			statistics.ArithmeticMean(),
			statistics.StandardDeviation())
		scaledPopulation[i], err = framework.NewEvaluatedCandidate(candidate.Candidate(),
			scaledFitness)
		if err != nil {
			panic(fmt.Sprintln("couldn't create evaluated candidate: ", err))
		}
	}
	return sel.delegate.Select(scaledPopulation, naturalFitnessScores, selectionSize, rng)
}

func (sel *SigmaScaling) String() string {
	return "Sigma Scaling"
}

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
