package engine

import (
	"testing"

	"github.com/aurelien-rainone/evolve/pkg/termination"
)

func TestEngineArgumentErrors(t *testing.T) {
	var (
		eng *Engine
		err error
	)

	t.Run("invalid elite count", func(t *testing.T) {
		for _, nelites := range []int{10, -1} {
			eng, err = New(zeroGenerator{}, intEvaluator{}, nil, Seed(99))
			check(t, err)
			_, _, err = eng.Evolve(10, Elites(nelites), EndOn(termination.GenerationCount(10)))
			if err != ErrEngineElite {
				t.Errorf("Evolve(Elites(%v)), wantErr %v, got %v", nelites, ErrEngineElite, err)
			}
		}
	})

	t.Run("no termination condition", func(t *testing.T) {
		eng, err = New(zeroGenerator{}, intEvaluator{}, nil, Seed(99))
		check(t, err)
		_, _, err = eng.Evolve(10)
		if err != ErrEngineTermination {
			t.Errorf("Evolve(), wantErr %v, got %v", ErrEngineTermination, err)
		}
	})

	t.Run("elite count", func(t *testing.T) {
		eng, err = New(zeroGenerator{}, intEvaluator{}, nil, Seed(99))
		check(t, err)
		_, _, err = eng.Evolve(0, EndOn(termination.GenerationCount(1)))
		if err != ErrEnginePopSize {
			t.Errorf("Evolve(), wantErr %v, got %v", ErrEnginePopSize, err)
		}
	})
}
