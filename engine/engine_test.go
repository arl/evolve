package engine

import (
	"testing"

	"github.com/arl/evolve/condition"
)

func TestEngineArgumentErrors(t *testing.T) {
	var (
		eng *Engine[int]
		err error
	)

	t.Run("invalid elite count", func(t *testing.T) {
		for _, nelites := range []int{10, -1} {
			eng, err = New[int](zeroFactory, intEvaluator{}, nil)
			check(t, err)
			_, _, err = eng.Evolve(10, Elites[int](nelites), EndOn[int](condition.GenerationCount[int](10)))
			if err == nil {
				t.Errorf("Evolve(Elites(%v)), want error, got nil", nelites)
			}
		}
	})

	t.Run("no condition condition", func(t *testing.T) {
		eng, err = New[int](zeroFactory, intEvaluator{}, nil)
		check(t, err)
		_, _, err = eng.Evolve(10)
		if err == nil {
			t.Error("Evolve(), want error, got nil")
		}
	})

	t.Run("elite count", func(t *testing.T) {
		eng, err = New[int](zeroFactory, intEvaluator{}, nil)
		check(t, err)
		_, _, err = eng.Evolve(0, EndOn[int](condition.GenerationCount[int](1)))
		if err == nil {
			t.Error("Evolve(), want error, got nil")
		}
	})
}
