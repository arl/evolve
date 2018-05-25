package api

import (
	"testing"

	"github.com/pkg/errors"
)

func TestEngineIllegalState(t *testing.T) {
	eng := &Engine{}
	_, err := eng.SatisfiedTerminationConditions()
	if errors.Cause(err) != ErrIllegalState {
		t.Errorf("engine not started, want cause %v, got %v", ErrIllegalState, err)
		return
	}
}
