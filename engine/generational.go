package engine

import (
	"math/rand"

	"github.com/arl/evolve"
)

// Generational implements a general-purpose engine for generational
// evolutionary algorithm.
type Generational[T any] struct {
	Operator  evolve.Operator[T]
	Evaluator evolve.Evaluator[T]
	Selection evolve.Selection[T]

	// Elites defines the number of candidates preserved via elitism for the
	// engine. By default it is set to 0, no elitism is applied.
	//
	// In elitism, a subset of the population with the best fitness scores is
	// preserved, unchanged, and placed into the successive generation.
	// Candidate solutions that are preserved unchanged through elitism remain
	// eligible for selection for breeding the remainder of the next generation.
	// This value must be non-negative and less than the population size or
	// Evolve will return en error
	Elites int
}

// Epoch performs a single step/iteration of the evolutionary process.
//
// pop is the population to evolve, sorted by fitness, the fittest first.
//
// Returns the updated population after the evolutionary process has proceeded
// by one step/iteration.
func (e *Generational[T]) Epoch(pop evolve.Population[T], rng *rand.Rand) evolve.Population[T] {
	nextpop := make([]T, 0, len(pop))

	// Perform elitism: straightforward copy the n fittest candidates into the
	// next generation, without any kind of selection.
	elite := make([]T, e.Elites)
	for i := 0; i < e.Elites; i++ {
		elite[i] = pop[i].Candidate
	}

	// Select the rest of population through natural selection.
	selected := e.Selection.Select(pop, e.Evaluator.IsNatural(), len(pop)-e.Elites, rng)

	// Apply genetic operators on the selected candidates.
	nextpop = e.Operator.Apply(append(nextpop, selected...), rng)

	// While the elites, if any, are added, untouched, to the next population.
	nextpop = append(nextpop, elite...)
	return evolve.EvaluatePopulation(nextpop, e.Evaluator, true)
}
