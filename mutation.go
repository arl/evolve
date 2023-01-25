package evolve

import "math/rand"

// Mutation implements the mutation evolutionnary operator. It modifies the
// genetic content of individuals in order to maintain diversity from one
// population to the next.
//
// At individual level, mutation is applied through a Mutater, which performs
// modification on a single element at once.
type Mutation[T any] struct {
	Mutater[T]
}

// Apply applies the mutation operator to all individuals in the provided
// population.
func (op *Mutation[T]) Apply(pop *Population[T], rng *rand.Rand) {
	for i := 0; i < pop.Len(); i++ {
		op.Mutate(&pop.Candidates[i], rng)
	}
}

// A Mutater mutates individuals.
type Mutater[T any] interface {
	// Mutate performs mutation on an individual.
	//
	// The fact that an individual is mutated or not depends on the particular
	// Mutater implementation.
	Mutate(*T, *rand.Rand)
}
