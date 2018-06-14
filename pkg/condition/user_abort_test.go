package condition

import (
	"testing"

	"github.com/arl/evolve"
)

func TestUserAbort(t *testing.T) {
	popdata := &evolve.PopulationData{}
	cond := NewUserAbort()

	if cond.IsSatisfied(popdata) {
		t.Errorf("should be false before user abort")
	}
	cond.Abort()
	if !cond.IsSatisfied(popdata) {
		t.Errorf("should be true before user abort")
	}
	cond.Reset()
	if cond.IsSatisfied(popdata) {
		t.Errorf("should be false after reset")
	}
}
