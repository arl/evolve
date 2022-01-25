package condition

import (
	"testing"
	"time"

	"github.com/arl/evolve"
)

func TestElapsedTime(t *testing.T) {
	cond := 1 * ElapsedTime[any](time.Second)
	stats := &evolve.PopulationStats[any]{}

	stats.Elapsed = 100 * time.Millisecond
	if cond.IsSatisfied(stats) {
		t.Errorf("elapsed = %v, termination condition should not be satisfied", stats.Elapsed)
	}

	stats.Elapsed = time.Second
	if !cond.IsSatisfied(stats) {
		t.Errorf("elapsed = %v, termination condition should be satisfied", stats.Elapsed)
	}
}
