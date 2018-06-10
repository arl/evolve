package condition

import (
	"fmt"
	"time"

	"github.com/arl/evolve/pkg/api"
)

// ElapsedTime is satisfied when a time duration has elapsed.
type ElapsedTime time.Duration

// IsSatisfied returns true if the time duration has elapsed.
func (dur ElapsedTime) IsSatisfied(popdata *api.PopulationData) bool {
	return popdata.Elapsed >= time.Duration(dur)
}

// String returns a string representation of this condition.
func (dur ElapsedTime) String() string {
	return fmt.Sprintf("Elapsed Time (%v)", time.Duration(dur))
}
