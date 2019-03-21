package condition

import (
	"testing"
	"time"

	"github.com/arl/evolve"
)

func TestElapsedTime(t *testing.T) {
	cond := 1 * ElapsedTime(time.Second)
	stats := &evolve.PopulationStats{}

	stats.Elapsed = 100 * time.Millisecond
	if cond.IsSatisfied(stats) {
		t.Errorf("should not terminate before elapsed time")
	}

	stats.Elapsed = time.Second
	if !cond.IsSatisfied(stats) {
		t.Errorf("should terminate after elapsed time")
	}
}
