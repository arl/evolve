package xover

import (
	"constraints"
	"math/rand"

	"github.com/arl/evolve/generator"
)

// Mater is the interface implemented by objects defining the Mate function.
type Mater[T any] interface {
	// Mate performs crossover on a pair of parents to generate a pair of
	// offspring.
	//
	// parent1 and parent2 are the two individuals that provides the source
	// material for generating offspring.
	Mate(parent1, parent2 T, nxpts int, rng *rand.Rand) (T, T)
}

// Crossover implements a standard crossover operator.
//
// It supports all crossover processes that operate on a pair of parent
// candidates.
// Both the number of crossovers points and the crossover probability are
// configurable. Crossover is applied to a proportion of selected parent pairs,
// with the remainder copied unchanged into the output population. The size of
// this evolved proportion is controlled by the code crossoverProbability
// parameter.
type Crossover[T any] struct {
	Mater[T]
	Probability generator.Float
	Points      generator.Generator[int]
}

// New creates a Crossover operator based off the provided Mater.
//
// The returned Crossover performs a one point crossover with 1.0 (i.e always)
// probability.
func New[T any](mater Mater[T]) *Crossover[T] {
	return &Crossover[T]{Mater: mater}
}

// Apply applies the crossover operation to the selected candidates.
//
// Pairs of candidates are chosen randomly from the selected candidates and
// subjected to crossover to produce a pair of offspring candidates. The
// selected candidates, sel, are the evolved individuals that have survived to
// be eligible to reproduce.
//
// Returns the combined set of evolved offsprings generated by applying
// crossover to the selected candidates.
func (op *Crossover[T]) Apply(sel []T, rng *rand.Rand) []T {
	// Generate a slice of 0..n indices (n=len(sel) and pair candidates using
	// shuffled indices so that the evolution is not influenced by any ordering
	// artifacts from previous operations.
	idx := seq[int](len(sel))
	rand.Shuffle(len(sel), func(i, j int) {
		idx[i], idx[j] = idx[j], idx[i]
	})

	res := make([]T, 0, len(sel))
	for i := 0; i < len(sel); {
		p1 := sel[idx[i]]
		i++
		if i < len(sel) {
			p2 := sel[idx[i]]
			i++

			// Probability for this pair to be mated
			p := op.Probability.Next()

			npts := 0
			if rng.Float64() < p {
				// we got a crossover to perform, get/decide the number of
				// crossover points
				npts = int(op.Points.Next())
			}
			if npts > 0 {
				off1, off2 := op.Mate(p1, p2, npts, rng)
				res = append(res, off1, off2)
			} else {
				// If there is no crossover to perform, just add the parents to the
				// results unaltered.
				res = append(res, p1, p2)
			}
		} else {
			// If we have an odd number of selected candidates, we can't pair up
			// the last one so just leave it unmodified.
			res = append(res, p1)
		}
	}
	return res
}

// seq returns a slice containing the sequence of consecutive numbers from 0 to n.
func seq[T constraints.Integer | constraints.Float](n int) []T {
	s := make([]T, n)
	for i := 0; i < n; i++ {
		s[i] = T(i)
	}
	return s
}
