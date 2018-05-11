package engine

import (
	"fmt"
	"testing"

	"github.com/aurelien-rainone/evolve/pkg/api"
	"github.com/pkg/errors"
)

func TestEngineBaseIllegalState(t *testing.T) {
	base := NewBase(nil, nil, nil, nil)
	_, err := base.SatisfiedTerminationConditions()
	if err == nil {
		t.Errorf("engine not started, want an error, got nil")
		return
	}
	fmt.Println("err.Cause", errors.Cause(err))
	if errors.Cause(err) != api.ErrIllegalState {
		t.Errorf("error cause should be api.ErrIllegalState, got %v", err)
	}
}
