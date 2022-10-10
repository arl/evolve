package engine

import (
	"errors"
	"math/rand"
	"runtime"
	"sort"
	"time"

	"github.com/arl/evolve"
	"github.com/arl/evolve/pkg/mt19937"
)

// Engine runs an evolutionary algorithm, following all the steps of evolution,
// from the creation of the initial population to the end of evolution.
type Engine[T any] struct {
	// Factory creates candidate solutions of type T.
	Factory evolve.Factory[T]

	// Evaluator evaluates candidate solutions fitness.
	Evaluator evolve.Evaluator[T]

	Epocher evolve.Epocher[T]

	EndConditions []evolve.Condition[T]

	// Observers of the evolution process.
	Observers []Observer[T]

	// Seeds provides the engine with a set of candidates to seed the starting
	// population with. Successive calls to Seeds will replace the set of seed
	// candidates set in the previous call.
	Seeds []T

	// RNG is the source of randomness of the engine. If nil, it's set to a
	// mt19937 pseudo random number generator.
	RNG *rand.Rand

	// Concurrency is the number of concurrent workers (defaults to the number of cores).
	Concurrency int

	stats *evolve.Dataset
}

// AddObserver adds an observer of the evolution process.
func (e *Engine[T]) AddObserver(o Observer[T]) {
	e.Observers = append(e.Observers, o)
}

// RemoveObserver removes an observer of the evolution process.
func (e *Engine[T]) RemoveObserver(o Observer[T]) {
	for i := range e.Observers {
		if e.Observers[i] == o {
			e.Observers = append(e.Observers[:i], e.Observers[i+1:]...)
			return
		}
	}
}

// Evolve runs the evolutionary algorithm until one of the termination
// conditions is met, then return the entire population present during the final
// generation.
//
// size is the number of candidate in the population. They whole population is
// generated for the first generation, unless some seed candidates are provided
// with Seeds. size must be at least 1 or Evolve will return en error.
//
// At least one termination condition must be defined with EndOn, or Evolve will
// return an error.
func (e *Engine[T]) Evolve(popsize int) (*evolve.Population[T], []evolve.Condition[T], error) {
	if popsize <= 0 {
		return nil, nil, errors.New("invalid population size")
	}
	if len(e.EndConditions) == 0 {
		return nil, nil, errors.New("no termination condition specified")
	}

	if e.Concurrency == 0 {
		e.Concurrency = runtime.NumCPU()
	}

	if e.RNG == nil {
		seed := time.Now().UnixNano()
		e.RNG = rand.New(mt19937.New(seed))
	}

	// Track down evolution stats in a dataset.
	e.stats = evolve.NewDataset(popsize)

	var ngen int
	start := time.Now()

	pop := evolve.SeedPopulation(e.Factory, popsize, e.Seeds, e.RNG)

	var satisfied []evolve.Condition[T]

	// Evaluate initial population fitness
	evpop := evolve.EvaluatePopulation(pop, e.Evaluator, e.Concurrency)
	for {
		// Sort population according to fitness.
		if e.Evaluator.IsNatural() {
			sort.Sort(sort.Reverse(evpop))
		} else {
			sort.Sort(evpop)
		}

		// compute population stats
		data := e.updateStats(evpop, ngen, time.Since(start))

		// check for termination conditions
		satisfied = satisfiedConditions(data, e.EndConditions)
		if satisfied != nil {
			break
		}

		// perform evolution
		evpop = e.Epocher.Epoch(evpop, e.RNG)

		ngen++
	}
	return evpop, satisfied, nil
}

func (e *Engine[T]) updateStats(pop *evolve.Population[T], ngen int, elapsed time.Duration) *evolve.PopulationStats[T] {
	e.stats.Clear()
	for i := 0; i < pop.Len(); i++ {
		e.stats.AddValue(pop.Fitness[i])
	}

	// Notify observers with the population state
	stats := evolve.PopulationStats[T]{
		Best:        pop.Candidates[0],
		BestFitness: pop.Fitness[0],
		Mean:        e.stats.ArithmeticMean(),
		StdDev:      e.stats.StandardDeviation(),
		Natural:     e.Evaluator.IsNatural(),
		Size:        e.stats.Len(),
		Generation:  ngen,
		Elapsed:     elapsed,
	}

	for _, o := range e.Observers {
		o.Observe(&stats)
	}
	return &stats
}

// satisfiedConditions returns the satisfied conditions, or nil if none of them are.
func satisfiedConditions[T any](stats *evolve.PopulationStats[T], conds []evolve.Condition[T]) []evolve.Condition[T] {
	var c []evolve.Condition[T]
	for _, cond := range conds {
		if cond.IsSatisfied(stats) {
			c = append(c, cond)
		}
	}
	return c
}
