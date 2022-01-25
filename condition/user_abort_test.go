package condition

import (
	"testing"

	"github.com/arl/evolve"
)

func TestUserAbort(t *testing.T) {
	stats := &evolve.PopulationStats[any]{}
	var cond UserAbort[any]

	if cond.IsSatisfied(stats) {
		t.Errorf("zero-value UserAbort should be a non-satisfied termination condition")
	}
	cond.Abort()
	if !cond.IsSatisfied(stats) {
		t.Errorf("after Abort, termination condition should be satisfied")
	}
	cond.Reset()
	if cond.IsSatisfied(stats) {
		t.Errorf("after Reset, termination condition should not be satisfied")
	}
}
