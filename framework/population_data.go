package framework

import (
	"time"
)

// PopulationData is an immutable data object containing statistics about the
// state of an evolved population and a reference to the fittest candidate
// solution in the population.
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

// NewPopulationData creates a new immutable PopulationData object.
func NewPopulationData(bestCandidate Candidate,
	bestCandidateFitness, meanFitness, fitnessStandardDeviation float64,
	naturalFitness bool,
	populationSize, eliteCount, generationNumber int,
	elapsedTime time.Duration) *PopulationData {
	return &PopulationData{

		bestCandidate:            bestCandidate,
		bestCandidateFitness:     bestCandidateFitness,
		meanFitness:              meanFitness,
		fitnessStandardDeviation: fitnessStandardDeviation,
		naturalFitness:           naturalFitness,
		populationSize:           populationSize,
		eliteCount:               eliteCount,
		generationNumber:         generationNumber,
		elapsedTime:              elapsedTime,
	}
}

// BestCandidate returns the fittest candidate present in the population.
func (pd *PopulationData) BestCandidate() Candidate {
	return pd.bestCandidate
}

// BestCandidateFitness returns the fitness score of the fittest candidate.
func (pd *PopulationData) BestCandidateFitness() float64 {
	return pd.bestCandidateFitness
}

// MeanFitness returns the arithmetic mean fitness of individual candidates.
func (pd *PopulationData) MeanFitness() float64 {
	return pd.meanFitness
}

// FitnessStandardDeviation returns a statistical measure of variation in
// fitness scores within the population.
func (pd *PopulationData) FitnessStandardDeviation() float64 {
	return pd.fitnessStandardDeviation
}

// IsNaturalFitness indicates whether the fitness scores are natural or
// non-natural.
//
// Returns true if higher fitness scores indicate fitter individuals, false
// otherwise.
func (pd *PopulationData) IsNaturalFitness() bool {
	return pd.naturalFitness
}

// PopulationSize returns the number of individuals in the current population.
func (pd *PopulationData) PopulationSize() int {
	return pd.populationSize
}

// EliteCount return the number of candidates preserved via elitism.
func (pd *PopulationData) EliteCount() int {
	return pd.eliteCount
}

// GenerationNumber returns the number of this generation (zero-based).
func (pd *PopulationData) GenerationNumber() int {
	return pd.generationNumber
}

// ElapsedTime returns the amount of time (in milliseconds) since the start of
// the evolutionary algorithm's execution.
func (pd *PopulationData) ElapsedTime() time.Duration {
	return pd.elapsedTime
}
