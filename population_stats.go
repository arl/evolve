package evolve

import "time"

// PopulationStats contains statistics about the state of an evolved population
// and a reference to the fittest candidate solution in the population.
type PopulationStats[T any] struct {
	// Best is the fittest individual of the population.
	Best Individual[T]

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

	// Generation is the 0-based index of the generation that was processed.
	Generation int

	// Elapsed is the duration elapsed since the evolution start.
	Elapsed time.Duration
}
