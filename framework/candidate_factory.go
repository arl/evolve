package framework

import (
	"math/rand"
)

// A CandidateFactory creates new populations of candidates.
type CandidateFactory interface {

	// GenerateInitialPopulation creates an initial population of candidates.
	//
	// If more control is required over the composition of the initial
	// population, consider the SeedInitialPopulation method.
	//
	// Returns an initial population of candidate solutions.
	GenerateInitialPopulation(
		populationSize int,
		rng *rand.Rand) []Candidate

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
	// - populationSize is the size of the initial population.
	// - seedCandidates is the slice of candidates to seed the population with.
	// Number of candidates must be no bigger than the population size.
	// - rng is the random number generator to use when creating additional
	// candidates to fill the population when the number of seed candidates is
	// insufficient. This can be null if and only if the number of seed
	// candidates provided is sufficient to fully populate the initial
	// population.
	//
	// Return an initial population of candidate solutions, including the
	// specified seed candidates.
	SeedInitialPopulation(
		populationSize int,
		seedCandidates []Candidate,
		rng *rand.Rand) []Candidate
}

// RandomCandidateGenerator is the interface implemented by objects that
// generate random candidates
type RandomCandidateGenerator interface {

	// GenerateRandomCandidate randomly create a single candidate solution.
	GenerateRandomCandidate(rng *rand.Rand) Candidate
}
