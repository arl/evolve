package base

import (
	"time"
)

// Population is an immutable data object containing statistics about the state
// of an evolved population and a reference to the fittest candidate solution in
// the population.
type PopulationData struct {
	// The fittest candidate present in the population.
	bestCandidate Candidate

	// The fitness score for the fittest candidate in the population.
	bestCandidateFitness float64

	// The arithmetic mean of fitness scores for each member of the population.
	meanFitness float64

	//  A measure of the variation in fitness scores.
	fitnessStandardDeviation float64

	// True if higher fitness scores are better, false otherwise.
	naturalFitness bool

	// The number of individuals in the population.
	populationSize int

	// The number of candidates preserved via elitism.
	eliteCount int

	// The (zero-based) number of the last generation that was processed.
	generationNumber int

	// The number of milliseconds since the start of the algorithm.
	elapsedTime time.Duration
}
