package termination

import (
	"time"

	"github.com/aurelien-rainone/evolve/base"
)

// ElapsedTime terminates evolution after a pre-determined period of time has
// elapsed.
type ElapsedTime struct {
	maxDuration time.Duration
}

// NewElapsedTime creates an ElapsedTime termination condition.
func NewElapsedTime(maxDuration time.Duration) ElapsedTime {
	if maxDuration <= 0 {
		panic("Duration must be positive")
	}
	return ElapsedTime{maxDuration: maxDuration}
}

// ShouldTerminate returns true if the pre-configured maximum permitted time
// has elapsed.
func (tc ElapsedTime) ShouldTerminate(populationData base.PopulationData) bool {
	return populationData.ElapsedTime() >= tc.maxDuration
}
