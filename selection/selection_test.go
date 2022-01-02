package selection

import (
	"math/rand"
	"testing"
	"time"

	"github.com/arl/evolve"
	"github.com/stretchr/testify/assert"
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

func testFitnessBasedSelection(t *testing.T, ss evolve.Selection[string], tpop testPopulation, natural bool) {
	rng := rand.New(rand.NewSource(99))

	// create the population
	pop := evolve.Population[string]{}
	for i := range tpop {
		cand := &evolve.Individual[string]{
			Candidate: tpop[i].name,
			Fitness:   tpop[i].fitness,
		}
		pop = append(pop, cand)
	}

	// apply selection
	sel := ss.Select(pop, natural, 4, rng)
	assert.Len(t, sel, 4, "got selection size:", len(sel), "want 4")

	// check candidate frequencies match expected results
	for i := range tpop {
		tcand := tpop[i]
		freq := frequency(sel, tcand.name)
		if freq < tcand.wantMin || freq > tcand.wantMax {
			t.Errorf("want %s selected [%v,%v] times, got %v", tcand.name, tcand.wantMin, tcand.wantMax, freq)
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
// of tpop, running the result to checkfn (returns nil or test fail message)
func testRandomBasedSelection(t *testing.T, ss evolve.Selection[string], tpop testPopulation, natural bool, n int, checkfn func([]string) error) {
	seed := time.Now().UnixNano()
	rng := rand.New(rand.NewSource(seed))

	// create the population
	pop := evolve.Population[string]{}
	for i := range tpop {
		cand := &evolve.Individual[string]{
			Candidate: tpop[i].name,
			Fitness:   tpop[i].fitness,
		}
		pop = append(pop, cand)
	}

	// apply selection
	selected := ss.Select(pop, natural, n, rng)
	msg := checkfn(selected)
	if msg != nil {
		t.Fatalf("%v with seed %v: %v", ss.String(), seed, msg)
	}
}

func errcheck(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("want error = nil, got %v", err)
	}
}
