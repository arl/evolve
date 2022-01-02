package condition

import (
	"testing"

	"github.com/arl/evolve"
)

func TestUserAbort(t *testing.T) {
	stats := &evolve.PopulationStats[any]{}
	cond := NewUserAbort[any]()

	if cond.IsSatisfied(stats) {
		t.Errorf("should be false before user abort")
	}
	cond.Abort()
	if !cond.IsSatisfied(stats) {
		t.Errorf("should be true before user abort")
	}
	cond.Reset()
	if cond.IsSatisfied(stats) {
		t.Errorf("should be false after reset")
	}
}
