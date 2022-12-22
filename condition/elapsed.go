package condition

import (
	"fmt"
	"time"

	"github.com/arl/evolve"
)

// ElapsedTime is satisfied when a time duration has elapsed.
type ElapsedTime[T any] time.Duration

// IsSatisfied returns true if the time duration has elapsed.
// TODO when/if method will accept type parameters we can get rid of T on ElapsedTime[T]
func (dur ElapsedTime[T]) IsSatisfied(stats *evolve.PopulationStats[T]) bool {
	return stats.Elapsed >= time.Duration(dur)
}

// String returns a string representation of this condition.
func (dur ElapsedTime[T]) String() string {
	return fmt.Sprintf("Elapsed Time (%v)", time.Duration(dur))
}
