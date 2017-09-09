package factory

import (
	"math/rand"
)

/**
 * Convenient base class for implementations of
 * {@link org.uncommons.watchmaker.framework.CandidateFactory}.
 * @param <T> The type of entity evolved by this engine.
 * @author Daniel Dyer
 */
type AbstractCandidateFactory struct {
	RandomCandidateGenerator
}

/**
 * Randomly, create an initial population of candidates.  If some
 * control is required over the composition of the initial population,
 * consider the overloaded {@link #generateInitialPopulation(int,Collection,Random)}
 * method.
 * @param populationSize The number of candidates to randomly create.
 * @param rng The random number generator to use when creating the random
 * candidates.
 * @return A randomly generated initial population of candidate solutions.
 */
func (f *AbstractCandidateFactory) GenerateInitialPopulation(
	populationSize int,
	rng *rand.Rand) []Candidate {

	population := make([]Candidate, populationSize)
	for i := range population {
		population[i] = f.GenerateRandomCandidate(rng)
	}
	return population
}

/**
 * {@inheritDoc}
 * If the number of seed candidates is less than the required population
 * size, the remainder of the population will be generated randomly via
 * the {@link #generateRandomCandidate(Random)} method.
 */
func (f *AbstractCandidateFactory) SeedInitialPopulation(
	populationSize int,
	seedCandidates []Candidate,
	rng *rand.Rand) []Candidate {

	if len(seedCandidates) > populationSize {
		panic("Too many seed candidates for specified population size.")
	}
	population := make([]Candidate, populationSize)
	for i := range seedCandidates {
		population[i] = seedCandidates[i]
	}
	for i := len(seedCandidates); i < populationSize; i++ {
		population[i] = f.GenerateRandomCandidate(rng)
	}
	return population
}
