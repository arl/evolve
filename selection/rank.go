package selection

import (
	"math/rand"

	"github.com/arl/evolve"
)

// A MappingFunc maps ranks (index of a candidate in a fitness-sorted
// population) to fitness scores. The general contract for the mapping function
// is:
//
//	f(rank) >= f(rank + 1) for all legal values of rank and assuming natural scores.
type MappingFunc func(rank, size int) float64

// RankBased is a configurable selection strategy that ignores absolute
// fitnesses and only consider the relative ordering of the candidates, or rank
// (in a sorted population).
//
// A mapping function, Map converts ranks into fitness scores. Actual selection
// is delegated to another selector, that uses the mapped fitness scores to
// drive the selection.
type RankBased[T any] struct {
	Selector evolve.Selection[T]
	Map      MappingFunc
}

// Select selects a given number of candidates from a population.
func (rb RankBased[T]) Select(pop *evolve.Population[T], natural bool, n int, rng *rand.Rand) []T {
	ranked := evolve.NewPopulation(pop.Len(), pop.Evaluator)
	for i := range pop.Candidates {
		ranked.Candidates[i] = pop.Candidates[i]
		// use candidate 1-based index
		ranked.Fitness[i] = rb.Map(i+1, pop.Len())
	}
	return rb.Selector.Select(ranked, true, n, rng)
}

func (RankBased[T]) String() string { return "Rank-Based Selection" }

// MapRankToScore maps ranks (sorted population index) to a relative
// pseudo-fitness score that can be used for fitness-proportionate selection.
//
// Returns size - rank
func MapRankToScore(rank, size int) float64 { return float64(size - rank) }

// Rank returns a preconfigured all-round rank-based selection strategy, using
// SUS (stochastic universal sampling) selection and MapRankToScore as mapping
// function.
func Rank[T any]() RankBased[T] {
	return RankBased[T]{
		Selector: SUS[T]{},
		Map:      MapRankToScore,
	}
}
