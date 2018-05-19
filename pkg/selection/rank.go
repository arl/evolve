package selection

import (
	"math/rand"

	"github.com/aurelien-rainone/evolve/pkg/api"
)

type rank struct{ selector api.Selection }

// NewRank returns a rank selection stragy rank, that is similar to
// fitness-proportionate selection except that is uses relative fitness rather
// than absolute fitness in order to determine the probability of selection for
// a given individual (i.e. the actual numerical fitness values are ignored and
// only the ordering of the sorted population is considered).
//
// selector is the proportionate selector that will be delegated to after
// converting rankings into relative fitness scores.
//
// Rank selection is implemented in terms of a mapping function
// mapRankToScore(int, int) and delegation to a fitness-proportionate selector.
// The mapping function converts ranks into relative fitness scores that are
// used to drive the delegate selector.
func NewRank(selector api.Selection) api.Selection {
	return rank{selector: selector}
}

// Rank is the default rank based selection strategy. It uses
// StochasticUniversalSampling as its selector.
var Rank = NewRank(StochasticUniversalSampling{})

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
func (rs rank) Select(
	pop api.Population,
	natural bool,
	size int,
	rng *rand.Rand) []interface{} {

	ranked := make(api.Population, len(pop))
	for i, cand := range pop {
		ranked[i] = &api.Individual{cand.Candidate,
			mapRankToScore(i+1, len(pop))}
	}
	return rs.selector.Select(ranked, true, size, rng)
}

func (rank) String() string { return "Rank Selection" }

// mapRankToScore maps a population index to a relative pseudo-fitness score
// that can be used for fitness-proportionate selection. The general contract
// for the mapping function is:
//  f(rank) >= f(rank + 1)
// For all legal values of rank, assuming natural scores.
//
// The default mapping function is a simple linear transformation, but this can
// be overridden by composition. Alternative implementations can be linear or
// non-linear and either natural or non-natural. rank is a zero-based index into
// the population (0 <= rank < population size)
//
// Returns population size - rank
func mapRankToScore(rank, size int) float64 {
	return float64(size - rank)
}
