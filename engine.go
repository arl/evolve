package evolve

import (
	"math/rand"
)

// An Epocher defines and implements an epoch, or generation, of an evolutionary
// algorithm.
type Epocher[T any] interface {
	// Epoch performs one epoch of the evolutionary process.
	//
	// It receives the population to evolve in that step, and returns another,
	// possibly evolved, population: the next generation.
	Epoch(*Population[T], *rand.Rand) *Population[T]
}

// EpochFunc is an adapter to allow the use of ordinary functions as Epocher. If
// f is a function with the appropriate signature, EpochFunc returns an object
// satisfying the Epocher interface, for which the Epoch method calls f.
type EpochFunc[T any] func(*Population[T], *rand.Rand) *Population[T]

func (f EpochFunc[T]) Epoch(pop *Population[T], rng *rand.Rand) *Population[T] {
	return f(pop, rng)
}
