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
type Factory interface {

	// New returns a new random candidate, using the provided pseudo-random
	// number generator.
	New(*rand.Rand) interface{}
}

// The FactoryFunc type is an adapter to allow the use of ordinary
// functions as candidate generators. If f is a function with the appropriate
// signature, FactoryFunc(f) is a Factory that calls f.
type FactoryFunc func(*rand.Rand) interface{}

// New calls f(rng) and returns its return value.
func (f FactoryFunc) New(rng *rand.Rand) interface{} { return f(rng) }

// GeneratePopulation returns a slice of count random candidates, generated
// with the provided Generator.
//
// If some control is required over the composition of the initial population,
// consider using SeedPopulation.
func GeneratePopulation(gen Factory, count int, rng *rand.Rand) []interface{} {
	pop := make([]interface{}, count)
	for i := 0; i < count; i++ {
		pop[i] = gen.New(rng)
	}
	return pop
}

// SeedPopulation seeds all or a part of an initial population with some
// candidates.
//
// Sometimes it is desirable to seed the initial population with some known
// good candidates, or partial solutions, in order to provide some hints for
// the evolution process. If the number of seed candidates is less than the
// required population size, gen will generate the additional candidates to fill
// the remaining spaces in the population.
func SeedPopulation(gen Factory, count int, seeds []interface{}, rng *rand.Rand) ([]interface{}, error) {
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
