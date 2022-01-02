package engine

import (
	"errors"
	"fmt"
	"math/rand"
	"sort"
	"time"

	"github.com/arl/evolve"
	"github.com/arl/evolve/pkg/mt19937"
)

// Engine runs an evolutionary algorithm, following all the steps of evolution,
// from the creation of the initial population to the end of evolution.
type Engine[T any] struct {
	obs     map[Observer[T]]struct{}
	rng     *rand.Rand
	factory evolve.Factory[T]
	eval    evolve.Evaluator[T]
	epoch   evolve.Epocher[T]
	stats   *evolve.Dataset
	nelites int
	seeds   []T
	conds   []evolve.Condition[T]
	size    int
}

// New creates an evolution engine.
//
// gen generates new random candidates solutions.
// eval evaluates fitness scores of candidates.
// epoch transforms a whole population into the next generation.
func New[T any](factory evolve.Factory[T], eval evolve.Evaluator[T], epoch evolve.Epocher[T], options ...func(*Engine[T]) error) (*Engine[T], error) {
	eng := Engine[T]{
		obs:     make(map[Observer[T]]struct{}),
		factory: factory,
		eval:    eval,
		epoch:   epoch,
	}
	for _, opt := range options {
		if err := opt(&eng); err != nil {
			return nil, err
		}
	}

	if eng.rng == nil {
		seed := time.Now().UnixNano()
		eng.rng = rand.New(mt19937.New(seed))
	}
	return &eng, nil
}

// AddObserver adds an observer of the evolution process.
func (e *Engine[T]) AddObserver(o Observer[T]) {
	e.obs[o] = struct{}{}
}

// RemoveObserver removes an observer of the evolution process.
func (e *Engine[T]) RemoveObserver(o Observer[T]) {
	delete(e.obs, o)
}

// Rand sets rng as the source of randomness of the engine.
func Rand[T any](rng *rand.Rand) func(*Engine[T]) error {
	return func(eng *Engine[T]) error {
		eng.rng = rng
		return nil
	}
}

// Observe adds an observer of the evolution process.
func Observe[T any](o Observer[T]) func(*Engine[T]) error {
	return func(eng *Engine[T]) error {
		eng.obs[o] = struct{}{}
		return nil
	}
}

// Elites defines the number of candidates preserved via elitism for the
// engine. By default it is set to 0, no elitism is applied.
//
// In elitism, a subset of the population with the best fitness scores is
// preserved, unchanged, and placed into the successive generation. Candidate
// solutions that are preserved unchanged through elitism remain eligible for
// selection for breeding the remainder of the next generation. This value must
// be non-negative and less than the population size or Evolve will return en
// error
func Elites[T any](n int) func(*Engine[T]) error {
	return func(eng *Engine[T]) error {
		if n < 0 || n >= eng.size {
			return errors.New("invalid number of elites")
		}
		eng.nelites = n
		return nil
	}
}

// Seeds provides the engine with a set of candidates to seed the starting
// population with. Successive calls to Seeds will replace the set of seed
// candidates set in the previous call.
func Seeds[T any](seeds []T) func(*Engine[T]) error {
	return func(eng *Engine[T]) error {
		eng.seeds = seeds
		return nil
	}
}

// EndOn adds a termination condition to the engine. The engine stops
// after one or more condition is met.
func EndOn[T any](cond evolve.Condition[T]) func(*Engine[T]) error {
	return func(eng *Engine[T]) error {
		eng.conds = append(eng.conds, cond)
		return nil
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
func (e *Engine[T]) Evolve(popsize int, options ...func(*Engine[T]) error) (evolve.Population[T], []evolve.Condition[T], error) {
	e.size = popsize
	for _, opt := range options {
		if err := opt(e); err != nil {
			return nil, nil, err
		}
	}

	if popsize <= 0 {
		return nil, nil, errors.New("invalid population size")
	}
	if len(e.conds) == 0 {
		return nil, nil, errors.New("no termination condition specified")
	}

	// create the dataset
	e.stats = evolve.NewDataset(popsize)

	var ngen int
	start := time.Now()

	pop, err := evolve.SeedPopulation(e.factory, popsize, e.seeds, e.rng)
	if err != nil {
		return nil, nil, fmt.Errorf("can't seed population: %v", err)
	}

	var satisfied []evolve.Condition[T]

	// Evaluate initial population fitness
	evpop := evolve.EvaluatePopulation(pop, e.eval, true)

	for {
		// Sort population according to fitness.
		if e.eval.IsNatural() {
			sort.Sort(sort.Reverse(evpop))
		} else {
			sort.Sort(evpop)
		}

		// compute population stats
		data := e.updateStats(evpop, ngen, time.Since(start))

		// check for termination conditions
		satisfied = shouldContinue(data, e.conds...)
		if satisfied != nil {
			break
		}

		// perform evolution
		evpop = e.epoch.Epoch(evpop, e.nelites, e.rng)

		ngen++
	}
	return evpop, satisfied, nil
}

func (e *Engine[T]) updateStats(pop evolve.Population[T], ngen int, elapsed time.Duration) *evolve.PopulationStats[T] {
	e.stats.Clear()
	for _, cand := range pop {
		e.stats.AddValue(cand.Fitness)
	}

	// Notify observers with the population state
	stats := evolve.PopulationStats[T]{
		BestCand:    pop[0].Candidate,
		BestFitness: pop[0].Fitness,
		Mean:        e.stats.ArithmeticMean(),
		StdDev:      e.stats.StandardDeviation(),
		Natural:     e.eval.IsNatural(),
		Size:        e.stats.Len(),
		NumElites:   e.nelites,
		GenNumber:   ngen,
		Elapsed:     elapsed,
	}

	for o := range e.obs {
		o.Observe(&stats)
	}
	return &stats
}

// shouldContinue determines whether or not the evolution should continue.
func shouldContinue[T any](stats *evolve.PopulationStats[T], conds ...evolve.Condition[T]) []evolve.Condition[T] {
	satisfied := make([]evolve.Condition[T], 0)
	for _, cond := range conds {
		if cond.IsSatisfied(stats) {
			satisfied = append(satisfied, cond)
		}
	}
	if len(satisfied) == 0 {
		return nil
	}
	return satisfied
}
