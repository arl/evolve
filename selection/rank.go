package selection

import (
	"math/rand"

	"github.com/arl/evolve"
)

// RankBased is selection strategy that is similar to fitness-proportionate
// selection except that is uses relative fitness rather than absolute fitness
// in order to determine the probability of selection for a given individual
// (i.e. the actual numerical fitness values are ignored and only the ordering
// of the sorted population is considered).
//
// RankBased is implemented in terms of a mapping function and delegation to a
// fitness-proportionate selector. The mapping function converts ranks into
// relative fitness scores that are used to drive the delegate selector.
type RankBased struct {
	Selector evolve.Selection
	Map      MappingFunc
}

// Select selects the specified number of candidates from the population.
//
// - pop must be sorted by descending fitness, i.e the fittest individual of the
// population should be pop[0].
// - natural indicates fitter individuals have fitness scores.
// - size is the number of individual selections to perform (not necessarily the
// number of distinct candidates to select, since the same individual may
// potentially be selected more than once).
//
// Returns the selected candidates.
func (rb RankBased) Select(
	pop evolve.Population,
	natural bool,
	size int,
	rng *rand.Rand) []interface{} {

	ranked := make(evolve.Population, len(pop))
	for i, cand := range pop {
		ranked[i] = &evolve.Individual{
			Candidate: cand.Candidate,
			// use candidate 1-based index
			Fitness: rb.Map(i+1, len(pop)),
		}
	}
	return rb.Selector.Select(ranked, true, size, rng)
}

func (RankBased) String() string { return "Rank-Based Selection" }

// MapRankToScore maps a population index to a relative pseudo-fitness score
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
// Returns size - rank
func MapRankToScore(rank, size int) float64 { return float64(size - rank) }

// MappingFunc is the type of functions that maps a population index to a
// relative pseudo-fitness score that can be used for fitness-proportionate
// selection. The general contract for the mapping function is:
//
//  f(rank) >= f(rank + 1) for all legal values of rank and assuming natural
//  scores.
type MappingFunc func(rank, size int) float64

// Rank is a preconfigured all-round rank-based selection strategy. It uses
// StochasticUniversalSampling as selector and MapRankToScore as mapping
// function.
var Rank = RankBased{
	Selector: StochasticUniversalSampling{},
	Map:      MapRankToScore,
}
