package selection

import (
	"testing"
)

// Unit test for fitness proportionate selection where observed selection
// frequencies correspond to expected frequencies.

func TestStochasticUniversalSamplingNatural(t *testing.T) {
	var tpop = testPopulation{
		{name: "Steve", fitness: 10.0, wantMin: 2, wantMax: 3},
		{name: "John", fitness: 4.5, wantMin: 1, wantMax: 2},
		{name: "Mary", fitness: 1.0, wantMin: 0, wantMax: 1},
		{name: "Gary", fitness: 0.5, wantMin: 0, wantMax: 1},
	}
	testFitnessBasedSelection(t, Rank, tpop, true)
}

func TestStochasticUniversalSamplingNonNatural(t *testing.T) {
	var tpop = testPopulation{
		{name: "Steve", fitness: 0.5, wantMin: 2, wantMax: 3},
		{name: "John", fitness: 1.0, wantMin: 1, wantMax: 2},
		{name: "Mary", fitness: 4.5, wantMin: 0, wantMax: 1},
		{name: "Gary", fitness: 10.0, wantMin: 0, wantMax: 1},
	}
	testFitnessBasedSelection(t, Rank, tpop, false)
}
