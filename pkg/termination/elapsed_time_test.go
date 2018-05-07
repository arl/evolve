package termination

import (
	"testing"
	"time"

	"github.com/aurelien-rainone/evolve/pkg/api"
)

func TestElapsedTime(t *testing.T) {
	cond := 1 * ElapsedTime(time.Second)

	if cond.ShouldTerminate(
		&api.PopulationData{struct{}{}, 0, 0, 0, true, 2, 0, 0, 100 * time.Millisecond}) {
		t.Errorf("should not terminate before elapsed time")
	}

	if !cond.ShouldTerminate(
		&api.PopulationData{struct{}{}, 0, 0, 0, true, 2, 0, 0, time.Second}) {
		t.Errorf("should terminate after elapsed time")
	}
}
