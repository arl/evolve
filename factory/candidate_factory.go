package factory

import (
	"math/rand"

	"github.com/aurelien-rainone/evolve/base"
)

// A CandidateFactory creates new populations of candidates.
type CandidateFactory interface {

	// Creates an initial population of candidates.  If more control is required
	// over the composition of the initial population, consider the overloaded
	// {@link #generateInitialPopulation(int,Collection,Random)} method.
	// @param populationSize The number of candidates to create.
	// @param rng The random number generator to use when creating the initial
	// candidates.
	// @return An initial population of candidate solutions.
	GenerateInitialPopulation(
		populationSize int,
		rng *rand.Rand) []base.Candidate

	// Sometimes it is desirable to seed the initial population with some
	// known good candidates, or partial solutions, in order to provide some
	// hints for the evolution process.  This method generates an initial
	// population, seeded with some initial candidates.  If the number of seed
	// candidates is less than the required population size, the factory should
	// generate additional candidates to fill the remaining spaces in the
	// population.
	// @param populationSize The size of the initial population.
	// @param seedCandidates Candidates to seed the population with.  Number
	// of candidates must be no bigger than the population size.
	// @param rng The random number generator to use when creating additional
	// candidates to fill the population when the number of seed candidates is
	// insufficient.  This can be null if and only if the number of seed
	// candidates provided is sufficient to fully populate the initial population.
	// @return An initial population of candidate solutions, including the
	// specified seed candidates.
	SeedInitialPopulation(
		populationSize int,
		seedCandidates []base.Candidate,
		rng *rand.Rand) []base.Candidate
}

type RandomCandidateGenerator interface {
	// Randomly create a single candidate solution.
	// @param rng The random number generator to use when creating the random
	// candidate.
	// @return A randomly-initialised candidate.
	GenerateRandomCandidate(rng *rand.Rand) base.Candidate
}
