package selection

import (
	"math/rand"
	"testing"
	"time"

	"github.com/aurelien-rainone/evolve/pkg/api"
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

func testFitnessBasedSelection(t *testing.T, ss api.Selection, tpop testPopulation, natural bool) {
	rng := rand.New(rand.NewSource(99))

	// create the population
	pop := api.EvaluatedPopulation{}
	for i := range tpop {
		cand, err := api.NewEvaluatedCandidate(tpop[i].name, tpop[i].fitness)
		if err != nil {
			t.Errorf("couldn't create evaluated candidate: %v", err)
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

func frequency(slice []interface{}, val interface{}) int {
	var count int
	for _, s := range slice {
		if s.(string) == val {
			count++
		}
	}
	return count
}

// function to check the selected candidates (returns nil of test fail message)
type popCheckFunc func(selected []interface{}) error

// test a random based selection strategy ss by selecting the n best candidates
// of tpop, running the result to f
func testRandomBasedSelection(t *testing.T, ss api.Selection, tpop testPopulation, natural bool, n int, f popCheckFunc) {
	seed := time.Now().UnixNano()
	rng := rand.New(rand.NewSource(seed))

	// create the population
	pop := api.EvaluatedPopulation{}
	for i := range tpop {
		cand, err := api.NewEvaluatedCandidate(tpop[i].name, tpop[i].fitness)
		if err != nil {
			t.Errorf("couldn't create evaluated candidate: %v", err)
		}
		pop = append(pop, cand)
	}

	// apply selection
	selected := ss.Select(pop, natural, n, rng)
	msg := f(selected)
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
