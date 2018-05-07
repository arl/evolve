package termination

import (
	"testing"
	"time"

	"github.com/aurelien-rainone/evolve/pkg/api"
)

func TestGenerationCount(t *testing.T) {
	cond := GenerationCount(5)

	if cond.ShouldTerminate(
		&api.PopulationData{struct{}{}, 0, 0, 0, true, 2, 0, 3, 100 * time.Millisecond}) {
		t.Errorf("should not terminate after 4th generation")
	}

	if !cond.ShouldTerminate(
		&api.PopulationData{struct{}{}, 0, 0, 0, true, 2, 0, 4, 100 * time.Millisecond}) {
		t.Errorf("should terminate after 5th generation")
	}
}
