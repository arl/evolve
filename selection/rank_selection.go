package selection

import (
	"fmt"
	"math/rand"

	"github.com/aurelien-rainone/evolve/framework"
)

// RankSelection is a selection strategy that is similar to
// fitness-proportionate selection except that is uses relative fitness rather
// than absolute fitness in order to determine the probability of selection for
// a given individual (i.e. the actual numerical fitness values are ignored and
// only the ordering of the sorted population is considered).
//
// Rank selection is implemented in terms of a mapping function
// mapRankToScore(int, int) and delegation to a fitness-proportionate selector.
// The mapping function converts ranks into relative fitness scores that are
// used to drive the delegate selector.
type RankSelection struct {
	delegate framework.SelectionStrategy
}

// NewRankSelection creates a rank-based selector with a linear mapping
// function.
//
// - delegate is the proprtionate selector that will be delegated to after
// converting rankings into relative fitness scores. The delegate parameter
// may be nil, in which case the RankSelection is created with a default
// rank-based selector and selection frequencies that correspond to expected
// values. T
func NewRankSelection(delegate framework.SelectionStrategy) *RankSelection {
	if delegate == nil {
		delegate = StochasticUniversalSampling{}
	}
	return &RankSelection{delegate: delegate}
}

// Select selects the specified number of candidates from the population.
//
// Implementations may assume that the population is sorted in descending
// order according to fitness (so the fittest individual is the first item
// in the list).
// NOTE: It is an error to call this method with an empty or null population.
//
// - population is the population from which to select.
// naturalFitnessScores indicates whether higher fitness values represent
// - fitter individuals or not.
// - selectionSize is the number of individual selections to make (not
// necessarily the number of distinct candidates to select, since the same
// individual may potentially be selected more than once).
//
// Returns a slice containing the selected candidates. Some individual
// candidates may potentially have been selected multiple times.
func (sel *RankSelection) Select(
	population framework.EvaluatedPopulation,
	naturalFitnessScores bool,
	selectionSize int,
	rng *rand.Rand) []framework.Candidate {

	rankedPopulation := make(framework.EvaluatedPopulation, len(population))
	var err error
	for index, candidate := range population {
		rankedPopulation[index], err = framework.NewEvaluatedCandidate(candidate.Candidate(),
			mapRankToScore(index+1, len(population)))
		if err != nil {
			panic(fmt.Sprintln("couldn't create evaluated candidate: ", err))
		}
	}
	return sel.delegate.Select(rankedPopulation, true, selectionSize, rng)
}

func (sel *RankSelection) String() string {
	return "Rank Selection"
}

// mapRankToScore maps a population index to a relative pseudo-fitness score
// that can be used for fitness-proportionate selection. The general
// contract for the mapping function
// is:
//  f(rank) >= f(rank + 1)
// for all legal values of rank, assuming natural scores.
// The default mapping function is a simple linear transformation, but this
// can be over-ridden by composition. Alternative implementations can be
// linear or non-linear and either natural or non-natural.
// - rank is a zero-based index into the population
//  (0 <= rank < populationSize)
// return populationSize - rank
func mapRankToScore(rank, populationSize int) float64 {
	return float64(populationSize - rank)
}
