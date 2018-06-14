package evolve

import "time"

// PopulationData contains statistics about the state of an evolved population
// and a reference to the fittest candidate solution in the population.
type PopulationData struct {

	// BestCandidate is the fittest candidate present in the population.
	// TODO: rename into Best (or: why not having an evluated candidate here, so
	// we would have the best candidate ANd their fitness)
	BestCand interface{}

	// BestFitness is the fitness score for the fittest candidate in the
	// population.
	BestFitness float64

	// Mean is the arithmetic mean of fitness scores for each member of
	// the population.
	Mean float64

	// StdDev is a measure of the variation in fitness scores.
	StdDev float64

	// Natural indicates, if true, that higher fitness is better.
	Natural bool

	// Size is the number of individuals in the population.
	Size int

	// NumElites is the number of candidates preserved via elitism.
	NumElites int

	// GenNumber is the (zero-based) number of the last generation that was
	// processed.
	GenNumber int

	// Elapsed is the duration elapsed since the evolution start.
	Elapsed time.Duration
}
