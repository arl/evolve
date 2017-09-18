package termination

import "github.com/aurelien-rainone/evolve/framework"

// GenerationCount terminates evolution after a set number of generations have
// passed.
type GenerationCount struct {
	generationCount int
}

// NewGenerationCount creates a GenerationCoun termination condition.
func NewGenerationCount(generationCount int) GenerationCount {
	if generationCount <= 0 {
		panic("Generation count must be positive")
	}
	return GenerationCount{generationCount: generationCount}
}

// ShouldTerminate reports whether or not evolution should finish at the
// current point.
//
// populationData is the information about the current state of evolution.
// This may be used to determine whether evolution should continue or not.
func (tc GenerationCount) ShouldTerminate(populationData *framework.PopulationData) bool {
	return populationData.GenerationNumber()+1 >= tc.generationCount
}
