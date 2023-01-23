package engine

import (
	"math/rand"
	"runtime"

	"github.com/arl/evolve"
)

// Generational implements a general-purpose engine for generational
// evolutionary algorithm.
type Generational[T any] struct {
	Operator  evolve.Operator[T]
	Evaluator evolve.Evaluator[T]
	Selection evolve.Selection[T]

	// NumElites defines the number of candidates preserved via elitism for the
	// engine. By default it is set to 0, no elitism is applied.
	//
	// In elitism, a subset of the population with the best fitness scores is
	// preserved, unchanged, and placed into the successive generation.
	// Candidate solutions that are preserved unchanged through elitism remain
	// eligible for selection for breeding the remainder of the next generation.
	// This value must be non-negative and less than the population size or
	// Evolve will return en error
	NumElites int

	// Number of concurrent processes to use (defaults to the number of cores).
	Concurrency int
}

// Epoch performs a single step/iteration of the evolutionary process.
//
// pop is the population to evolve, sorted by fitness, the fittest first.
//
// Returns the updated population after the evolutionary process has proceeded
// by one step/iteration.
func (e *Generational[T]) Epoch(pop *evolve.Population[T], rng *rand.Rand) *evolve.Population[T] {
	if e.Concurrency == 0 {
		e.Concurrency = runtime.NumCPU()
	}

	// nextpop := make([]T, 0, pop.Len())
	nextpop := evolve.NewPopulationWithCapacity(0, pop.Len(), e.Evaluator)

	// Perform elitism: straightforward copy the n fittest candidates into the
	// next generation, without any kind of selection.
	elite := make([]T, e.NumElites)
	for i := 0; i < e.NumElites; i++ {
		elite[i] = pop.Candidates[i]
	}

	// Select the rest of population through natural selection.
	selected := e.Selection.Select(pop, e.Evaluator.IsNatural(), pop.Len()-e.NumElites, rng)

	// Add selected candidates to the next population
	nextpop.Candidates = append(nextpop.Candidates, selected...)
	// Reslice the other slices to the same length of Candidates.
	nextpop.Fitness = nextpop.Fitness[0:pop.Len()]
	nextpop.Evaluated = nextpop.Evaluated[0:pop.Len()]

	// Apply genetic operators on the selected candidates.
	e.Operator.Apply(nextpop, rng)

	// Finally, add elites (if any), untouched, to the next population.
	nextpop.Candidates = append(nextpop.Candidates, elite...)
	// Reslice again the 2 other slices..
	nextpop.Fitness = nextpop.Fitness[0:pop.Len()]
	nextpop.Evaluated = nextpop.Evaluated[0:pop.Len()]

	nextpop.Evaluate(e.Concurrency)
	return nextpop
}
