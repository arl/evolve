package termination

import (
	"testing"

	"github.com/aurelien-rainone/evolve/pkg/api"
)

func TestUserAbort(t *testing.T) {
	// population data is irrelevant
	popdata := &api.PopulationData{}
	cond := NewUserAbort()

	if cond.ShouldTerminate(popdata) {
		t.Errorf("should not terminate before user abort")
	}

	if cond.IsAborted() {
		t.Errorf("should not have aborted without user intervention")
	}

	cond.Abort()

	if !cond.ShouldTerminate(popdata) {
		t.Errorf("should not terminate after user abort")
	}

	if !cond.IsAborted() {
		t.Errorf("should have aborted after user intervention")
	}
}
