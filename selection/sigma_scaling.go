package selection

import (
	"math/rand"

	"github.com/arl/evolve"
	"github.com/arl/evolve/pkg/dataset"
)

// SigmaScaling creates a sigma-scaled selection strategy. This is an
// alternative to straightforward fitness-proportionate selection such as that
// offered by RouletteWheel and SUS. Uses the mean population fitness and
// fitness standard deviation to adjust individual fitness scores.
//
// Early on in an evolutionary algorithm this may help to avoid premature
// convergence caused by the dominance of one or two relatively fit candidates
// in a population of mostly unfit individuals. It also helps to amplify minor
// fitness differences in a more mature population where the rate of improvement
// has slowed.
type SigmaScaling[T any] struct {
	// Selector is the proportionate selector that will be delegated to after
	// fitness scores have been adjusted using sigma scaling. If selector is
	// nil, then it's set to the SUS selector.
	Selector evolve.Selection[T]
}

// Select select candidates from the population using a sigma scaling strategy.
func (sel *SigmaScaling[T]) Select(pop *evolve.Population[T], natural bool, n int, rng *rand.Rand) []T {
	if sel.Selector == nil {
		sel.Selector = SUS[T]{}
	}

	stats := dataset.New(pop.Len())
	for i := 0; i < pop.Len(); i++ {
		stats.AddValue(pop.Fitness[i])
	}

	scaledPop := evolve.NewPopulation(pop.Len(), pop.Evaluator)
	for i := 0; i < pop.Len(); i++ {
		scaledPop.Candidates[i] = pop.Candidates[i]
		scaledPop.Fitness[i] = sigmaScaledFitness(pop.Fitness[i], stats.ArithmeticMean(), stats.StandardDeviation())
	}
	return sel.Selector.Select(scaledPop, natural, n, rng)
}

func (SigmaScaling[T]) String() string { return "Sigma Scaling" }

func sigmaScaledFitness(fitness, mean, stddev float64) float64 {
	if stddev == 0 {
		return 1
	}
	scaled := 1 + (fitness-mean)/(2*stddev)
	// Don't allow negative expected frequencies, use an arbitrary low but still
	// positive frequency of 1 time in 10 for extremely unfit individuals
	// (relative to the remainder of the population).
	if scaled > 0 {
		return scaled
	}
	return 0.1
}
