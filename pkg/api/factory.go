package api

import "math/rand"

// A Factory creates new populations of candidates.
type Factory interface {

	// GenPopulation creates an initial population of candidates.
	//
	// If more control is required over the composition of the initial
	// population, consider using SeedPopulation.
	//
	// Returns an initial population of candidate solutions.
	GenPopulation(size int, rng *rand.Rand) []interface{}

	// SeedPopulation seeds all or a part of an initial population
	// with some candidates.
	//
	// Sometimes it is desirable to seed the initial population with some known
	// good candidates, or partial solutions, in order to provide some hints for
	// the evolution process. This method generates an initial population,
	// seeded with some initial candidates. If the number of seed candidates is
	// less than the required population size, the factory should generate
	// additional candidates to fill the remaining spaces in the population.
	//
	// size is the size of the initial population.
	// cands is the slice of candidates to seed the population with. The number
	// of candidates must be no greater than the population size.  rng is the
	// random number generator to use when creating additional candidates to
	// fill the population when the number of seed candidates is insufficient.
	// It may be nil if and only if the number of seed candidates provided is
	// sufficient to fully populate the initial population.
	//
	// Returns an initial population of candidate solutions, including the
	// specified seed candidates.
	SeedPopulation(size int, cands []interface{}, rng *rand.Rand) []interface{}
}

// CandidateGenerator is the interface implemented by objects that
// generate random candidates
// TODO: can we come up with a better name?
type CandidateGenerator interface {

	// GenerateCandidate randomly create a single candidate solution.
	GenerateCandidate(rng *rand.Rand) interface{}
}
