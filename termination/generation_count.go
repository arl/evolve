package termination

import (
	"fmt"

	"github.com/aurelien-rainone/evolve/pkg/api"
)

// GenerationCount terminates evolution after a set number of generations have
// passed.
type GenerationCount int

// NewGenerationCount creates a GenerationCoun termination condition.
func NewGenerationCount(generationCount int) GenerationCount {
	if generationCount <= 0 {
		panic("Generation count must be positive")
	}
	return GenerationCount(generationCount)
}

// ShouldTerminate reports whether or not evolution should finish at the
// current point.
//
// populationData is the information about the current state of evolution.
// This may be used to determine whether evolution should continue or not.
func (tc GenerationCount) ShouldTerminate(populationData *api.PopulationData) bool {
	return populationData.GenerationNumber()+1 >= int(tc)
}

// String returns the termination condition representation as a string
func (tc GenerationCount) String() string {
	return fmt.Sprintf("Reached %d generations", tc)
}
