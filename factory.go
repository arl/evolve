package evolve

import (
	"errors"
	"math/rand"
)

// ErrTooManySeedCandidates is the error returned by SeedPopulation when the
// number of seed candidates is greater than the population size.
var ErrTooManySeedCandidates = errors.New("too many seed candidates for population size")

// A Factory generates random candidates.
//
// It is used by evolution engine to increase genetic diversity and/or add new
// candidates to a population.
type Factory[T any] interface {
	// New returns a new random candidate, using the provided pseudo-random
	// number generator.
	New(*rand.Rand) T
}

// The FactoryFunc type is an adapter to allow the use of ordinary
// functions as candidate generators. If f is a function with the appropriate
// signature, FactoryFunc(f) is a Factory that calls f.
type FactoryFunc[T any] func(*rand.Rand) T

// New calls f(rng) and returns its return value.
func (f FactoryFunc[T]) New(rng *rand.Rand) T { return f(rng) }

// GeneratePopulation returns a slice of n random candidates, generated
// with the provided Factory.
//
// If some control is required over the composition of the initial population,
// consider using SeedPopulation.
func GeneratePopulation[T any](fac Factory[T], n int, rng *rand.Rand) []T {
	pop := make([]T, 0, n)
	for i := 0; i < n; i++ {
		pop = append(pop, fac.New(rng))
	}
	return pop
}

// SeedPopulation returns a slice of n candidates, where a part of them are
// seeded while the rest is generated randomly using the provided factory.
// Sometimes it is desirable to seed the initial population with some known good
// candidates, providing some hints for the evolution process.
//
func SeedPopulation[T any](gen Factory[T], count int, seeds []T, rng *rand.Rand) ([]T, error) {
	if len(seeds) > count {
		return nil, ErrTooManySeedCandidates
	}

	// directory add the generated candidates to the backing array of seeds,
	// but seeds won't be modified
	for i := len(seeds); i < count; i++ {
		seeds = append(seeds, gen.New(rng))
	}
	return seeds, nil
}
