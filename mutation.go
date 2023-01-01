package evolve

import (
	"math/rand"
)

// Mutation implements the mutation evolutionnary operator. It modifies the
// genetic content of individuals in order to maintain diversity from one
// population to the next.
//
// At individual level, mutation is applied through a Mutater, which performs
// modification on a single element at once.
type Mutation[T any] struct {
	Mutater[T]
}

// NewMutation returns a operator that mutates chromosomes with the given mutater.
func NewMutation[T any](mutater Mutater[T]) *Mutation[T] {
	return &Mutation[T]{Mutater: mutater}
}

// Apply applies the mutation operator to all individuals in the provided population.
func (op *Mutation[T]) Apply(population []T, rng *rand.Rand) []T {
	muted := make([]T, len(population))
	for i, cand := range population {
		muted[i] = op.Mutate(cand, rng)
	}
	return muted
}

// A Mutater mutates individuals.
type Mutater[T any] interface {
	// Mutate performs mutation on an individual.
	//
	// The original individual must be let untouched while the mutant is
	// returned, unless no mutation is performed, in which case the original
	// individual can be returned.
	Mutate(T, *rand.Rand) T
}
