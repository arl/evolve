package factory

import (
	"math/rand"

	"github.com/aurelien-rainone/evolve/pkg/api"
)

// BaseFactory is a base for implementations of
// the CandidateFactory interface.
type BaseFactory struct{ api.RandomCandidateGenerator }

// GenerateInitialPopulation randomly creates an initial population of
// candidates.
//
// If some control is required over the composition of the initial population,
// consider the SeedInitialPopulation method.
//
// Returns a randomly generated initial population of candidate solutions.
func (f *BaseFactory) GenerateInitialPopulation(
	size int,
	rng *rand.Rand) []api.Candidate {

	pop := make([]api.Candidate, size)
	for i := range pop {
		pop[i] = f.GenerateRandomCandidate(rng)
	}
	return pop
}

// SeedInitialPopulation seeds all or a part of an initial population
// with some candidates.
//
// Sometimes it is desirable to seed the initial population with some known
// good candidates, or partial solutions, in order to provide some hints for
// the evolution process. This method generates an initial population,
// seeded with some initial candidates. If the number of seed candidates is
// less than the required population size, the factory should generate
// additional candidates to fill the remaining spaces in the population.
//
// populationSize is the size of the initial population.
// seedCandidates is the slice of candidates to seed the population with.
// Number of candidates must be no bigger than the population size.
// rng is the random number generator to use when creating additional candidates
// to fill the population when the number of seed candidates is insufficient.
// This can be null if and only if the number of seed candidates provided is
// sufficient to fully populate the initial population.
//
// Returns an initial population of candidate solutions, including the
// specified seed candidates.
func (f *BaseFactory) SeedInitialPopulation(
	size int,
	cands []api.Candidate,
	rng *rand.Rand) []api.Candidate {

	if len(cands) > size {
		panic("Too many seed candidates for specified population size.")
	}
	pop := make([]api.Candidate, size)
	copy(pop, cands)

	for i := len(cands); i < size; i++ {
		pop[i] = f.GenerateRandomCandidate(rng)
	}
	return pop
}
