package selection

import (
	"testing"
)

func TestSUS(t *testing.T) {
	naturalPop := testPopulation{
		{name: "Steve", fitness: 10.0, wantMin: 2, wantMax: 3},
		{name: "John", fitness: 4.5, wantMin: 1, wantMax: 2},
		{name: "Mary", fitness: 1.0, wantMin: 0, wantMax: 1},
		{name: "Gary", fitness: 0.5, wantMin: 0, wantMax: 1},
	}
	t.Run("natural", testFitnessBasedSelection(SUS[string]{}, naturalPop, true))

	nonNaturalPop := testPopulation{
		{name: "Steve", fitness: 0.5, wantMin: 2, wantMax: 3},
		{name: "John", fitness: 1.0, wantMin: 1, wantMax: 2},
		{name: "Mary", fitness: 4.5, wantMin: 0, wantMax: 1},
		{name: "Gary", fitness: 10.0, wantMin: 0, wantMax: 1},
	}
	t.Run("non-natural", testFitnessBasedSelection(SUS[string]{}, nonNaturalPop, false))
}
