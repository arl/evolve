package condition

import (
	"testing"
	"time"

	"github.com/arl/evolve"
)

func TestElapsedTime(t *testing.T) {
	cond := 1 * ElapsedTime(time.Second)
	popdata := &evolve.PopulationData{}

	popdata.Elapsed = 100 * time.Millisecond
	if cond.IsSatisfied(popdata) {
		t.Errorf("should not terminate before elapsed time")
	}

	popdata.Elapsed = time.Second
	if !cond.IsSatisfied(popdata) {
		t.Errorf("should terminate after elapsed time")
	}
}
