package engine

import (
	"testing"

	"github.com/arl/evolve/pkg/termination"
)

func TestEngineArgumentErrors(t *testing.T) {
	var (
		eng *Engine
		err error
	)

	t.Run("invalid elite count", func(t *testing.T) {
		for _, nelites := range []int{10, -1} {
			eng, err = New(zeroGenerator{}, intEvaluator{}, nil)
			check(t, err)
			_, _, err = eng.Evolve(10, Elites(nelites), EndOn(termination.GenerationCount(10)))
			if err == nil {
				t.Errorf("Evolve(Elites(%v)), want error, got nil", nelites)
			}
		}
	})

	t.Run("no termination condition", func(t *testing.T) {
		eng, err = New(zeroGenerator{}, intEvaluator{}, nil)
		check(t, err)
		_, _, err = eng.Evolve(10)
		if err == nil {
			t.Error("Evolve(), want error, got nil")
		}
	})

	t.Run("elite count", func(t *testing.T) {
		eng, err = New(zeroGenerator{}, intEvaluator{}, nil)
		check(t, err)
		_, _, err = eng.Evolve(0, EndOn(termination.GenerationCount(1)))
		if err == nil {
			t.Error("Evolve(), want error, got nil")
		}
	})
}
