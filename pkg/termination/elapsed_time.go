package termination

import (
	"fmt"
	"time"

	"github.com/aurelien-rainone/evolve/pkg/api"
)

// ElapsedTime terminates evolution after a pre-determined period of time has
// elapsed.
type ElapsedTime time.Duration

// ShouldTerminate returns true if the pre-configured maximum permitted time
// has elapsed.
func (tc ElapsedTime) ShouldTerminate(popdata *api.PopulationData) bool {
	return popdata.Elapsed >= time.Duration(tc)
}

// String returns the termination condition representation as a string
func (tc ElapsedTime) String() string {
	return fmt.Sprintf("Elapsed Time (duration: %v)", time.Duration(tc))
}
