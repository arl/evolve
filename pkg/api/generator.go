package api

import (
	"errors"
	"math/rand"
)

// ErrTooManySeedCandidates is the error returned by SeedPopulation when the
// number of seed candidates is greater than the population size.
var ErrTooManySeedCandidates = errors.New("Too many seed candidates for population size")

// CandidateGenerator is the interface implemented by objects that
// generate random candidates
type CandidateGenerator interface {

	// GenerateCandidate randomly create a single candidate solution.
	GenerateCandidate(rng *rand.Rand) interface{}
}

// GeneratePopulation returns a slice of count candidates, randomly generated
// with gen.
//
// If some control is required over the composition of the initial population,
// consider using SeedPopulation.
func GeneratePopulation(gen CandidateGenerator, count int, rng *rand.Rand) []interface{} {
	pop := make([]interface{}, count)
	for i := 0; i < count; i++ {
		pop[i] = gen.GenerateCandidate(rng)
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
func SeedPopulation(gen CandidateGenerator, count int, seeds []interface{}, rng *rand.Rand) ([]interface{}, error) {
	if len(seeds) > count {
		return nil, ErrTooManySeedCandidates
	}

	// directory add the generated candidates to the backing array of seeds,
	// but seeds won't be modified
	for i := len(seeds); i < count; i++ {
		seeds = append(seeds, gen.GenerateCandidate(rng))
	}
	return seeds, nil
}
