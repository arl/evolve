package selection

import (
	"math/rand"
	"testing"
	"time"

	"github.com/arl/evolve"
)

type testPopulation []struct {
	name             string
	fitness          float64
	wantMin, wantMax int
}

// test population and selection results for fitness-based (as opposed to random
// based) selection strategies with natural fitness (higher is better).
var fitnessBasedPopNatural = testPopulation{
	{name: "Steve", fitness: 10.0, wantMin: 1, wantMax: 2},
	{name: "John", fitness: 4.5, wantMin: 1, wantMax: 2},
	{name: "Mary", fitness: 1.0, wantMin: 0, wantMax: 1},
	{name: "Gary", fitness: 0.5, wantMin: 0, wantMax: 1},
}

// test population and selection results for fitness-based (as opposed to random
// based) selection strategies with non natural fitness (lower is better).
var fitnessBasedPopNonNatural = testPopulation{
	{name: "Steve", fitness: 0.5, wantMin: 1, wantMax: 2},
	{name: "John", fitness: 1.0, wantMin: 1, wantMax: 2},
	{name: "Mary", fitness: 4.5, wantMin: 0, wantMax: 1},
	{name: "Gary", fitness: 10.0, wantMin: 0, wantMax: 1},
}

// test population and selection results for fitness-based (as opposed to
// random based) and linear selection strategies
var fitnessBasedPopAllEqual = testPopulation{
	{name: "Steve", fitness: 4.0, wantMin: 1, wantMax: 1},
	{name: "John", fitness: 4.0, wantMin: 1, wantMax: 1},
	{name: "Mary", fitness: 4.0, wantMin: 1, wantMax: 1},
	{name: "Gary", fitness: 4.0, wantMin: 1, wantMax: 1},
}

// test population, sorted with natural fitness
var randomBasedPopNatural = testPopulation{
	{name: "Steve", fitness: 10.0},
	{name: "John", fitness: 9.1},
	{name: "Mary", fitness: 8.4},
	{name: "Gary", fitness: 6.2},
}

// test population, sorted with non-natural fitness
var randomBasedPopNonNatural = testPopulation{
	{name: "Gary", fitness: 6.2},
	{name: "Mary", fitness: 8.4},
	{name: "John", fitness: 9.1},
	{name: "Steve", fitness: 10.0},
}

func testFitnessBasedSelection(ss evolve.Selection[string], tpop testPopulation, natural bool) func(*testing.T) {
	return func(t *testing.T) {
		rng := rand.New(rand.NewSource(99))

		// Create population.
		pop := evolve.NewPopulation[string](len(tpop))
		for i := range tpop {
			pop.Candidates[i] = tpop[i].name
			pop.Fitness[i] = tpop[i].fitness
		}

		// Apply selection.
		sel := ss.Select(pop, natural, 4, rng)
		if len(sel) != 4 {
			t.Fatalf("selected %d individuals, want 4", len(sel))
		}

		// Check that frequencies of candidate selection match expectations.
		for i := range tpop {
			tcand := tpop[i]
			freq := frequency(sel, tcand.name)
			if freq < tcand.wantMin || freq > tcand.wantMax {
				t.Errorf("freq = %v want = %s selected [%v,%v] times", freq, tcand.name, tcand.wantMin, tcand.wantMax)
			}
		}
	}
}

func frequency[T comparable](slice []T, val T) int {
	var count int
	for _, s := range slice {
		if s == val {
			count++
		}
	}
	return count
}

// test a random based selection strategy ss by selecting the n best candidates
// of tpop, running the result to check.
func testRandomBasedSelection(s evolve.Selection[string], tpop testPopulation, natural bool, n int, check func([]string) error) func(*testing.T) {
	return func(t *testing.T) {
		seed := time.Now().UnixNano()
		rng := rand.New(rand.NewSource(seed))

		pop := evolve.NewPopulation[string](len(tpop))
		for i := range tpop {
			pop.Candidates[i] = tpop[i].name
			pop.Fitness[i] = tpop[i].fitness
		}

		// Apply selection.
		selected := s.Select(pop, natural, n, rng)
		msg := check(selected)
		if msg != nil {
			t.Fatalf("%v with seed %v: %v", s.String(), seed, msg)
		}
	}
}
