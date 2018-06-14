package condition

import (
	"testing"

	"github.com/arl/evolve"
)

func TestGenerationCount(t *testing.T) {
	cond := GenerationCount(5)
	popdata := &evolve.PopulationData{}

	popdata.GenNumber = 3
	if cond.IsSatisfied(popdata) {
		t.Errorf("should not terminate after 4th generation")
	}

	popdata.GenNumber = 4
	if !cond.IsSatisfied(popdata) {
		t.Errorf("should terminate after 5th generation")
	}
}
