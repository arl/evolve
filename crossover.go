package evolve

import (
	"math/rand"

	"github.com/arl/evolve/generator"
	"golang.org/x/exp/constraints"
)

// Mater is the interface implemented by objects defining the Mate function.
type Mater[T any] interface {
	// Mate performs crossover on a pair of parents to generate a pair of
	// offspring.
	//
	// parent1 and parent2 are the two individuals that provides the source
	// material for generating offspring.
	Mate(parent1, parent2 T, rng *rand.Rand) (T, T)
}

// Crossover implements a standard crossover operator.
//
// It supports all crossover processes that operate on a pair of parent
// candidates via the Mater interface.
// Both the number of crossovers points and the crossover probability are
// configurable. Crossover is applied to a proportion of selected parent pairs,
// with the remainder copied unchanged into the output population. The size of
// this evolved proportion is controlled by the code crossoverProbability
// parameter.
type Crossover[T any] struct {
	Mater[T]
	Probability generator.Float
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
func (op *Crossover[T]) Apply(pop *Population[T], rng *rand.Rand) {
	// Shuffle candidates so that evolution is not biased by previous
	// operations.
	rng.Shuffle(pop.Len(), pop.Swap)

	for i := 0; i < pop.Len()-1; i += 2 {
		j := i + 1
		p1 := pop.Candidates[i]
		p2 := pop.Candidates[j]

		// Probability for this pair to be mated
		p := op.Probability.Next()
		if rng.Float64() >= p {
			continue // Nothing to do
		}

		// Delegate actual crossover to the mater.
		off1, off2 := op.Mate(p1, p2, rng)
		pop.Candidates[i] = off1
		pop.Candidates[j] = off2
	}
}

// seq returns a slice containing the sequence of consecutive numbers from 0 to n.
// TODO(arl) remove
func seq[T constraints.Integer | constraints.Float](n int) []T {
	s := make([]T, n)
	for i := 0; i < n; i++ {
		s[i] = T(i)
	}
	return s
}
