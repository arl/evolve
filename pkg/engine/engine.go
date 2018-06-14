package engine

import (
	"math/rand"
	"sort"
	"time"

	"github.com/arl/evolve"
	"github.com/arl/evolve/pkg/mt19937"
	"github.com/pkg/errors"
)

// Engine runs an evolutionary algorithm, following all the steps of evolution,
// from the creation of the initial population to the end of evolution.
type Engine struct {
	obs     map[Observer]struct{}
	rng     *rand.Rand
	gen     evolve.Generator
	eval    evolve.Evaluator
	epoch   evolve.Epocher
	stats   *evolve.Dataset
	nelites int
	seeds   []interface{}
	conds   []evolve.Condition
	size    int
}

// New creates an evolution engine.
//
// gen generates new random candidates solutions.
// eval evaluates fitness scores of candidates.
// epoch transforms a whole population into the next generation.
func New(gen evolve.Generator, eval evolve.Evaluator, epoch evolve.Epocher, options ...func(*Engine) error) (*Engine, error) {
	eng := Engine{
		obs:   make(map[Observer]struct{}),
		gen:   gen,
		eval:  eval,
		epoch: epoch,
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
func (e *Engine) AddObserver(o Observer) {
	e.obs[o] = struct{}{}
}

// RemoveObserver removes an observer of the evolution process.
func (e *Engine) RemoveObserver(o Observer) {
	delete(e.obs, o)
}

// Rand sets rng as the source of randomness of the engine.
func Rand(rng *rand.Rand) func(*Engine) error {
	return func(eng *Engine) error {
		eng.rng = rng
		return nil
	}
}

// Observe adds an observer of the evolution process.
func Observe(o Observer) func(*Engine) error {
	return func(eng *Engine) error {
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
func Elites(n int) func(*Engine) error {
	return func(eng *Engine) error {
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
func Seeds(seeds []interface{}) func(*Engine) error {
	return func(eng *Engine) error {
		eng.seeds = seeds
		return nil
	}
}

// EndOn adds a termination condition to the engine. The engine stops
// after one or more condition is met.
func EndOn(cond evolve.Condition) func(*Engine) error {
	return func(eng *Engine) error {
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
func (e *Engine) Evolve(popsize int, options ...func(*Engine) error) (evolve.Population, []evolve.Condition, error) {
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

	pop, err := evolve.SeedPopulation(e.gen, popsize, e.seeds, e.rng)
	if err != nil {
		return nil, nil, errors.Wrap(err, "can't seed population")
	}

	var satisfied []evolve.Condition

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

func (e *Engine) updateStats(pop evolve.Population, ngen int, elapsed time.Duration) *evolve.PopulationData {

	e.stats.Clear()
	for _, cand := range pop {
		e.stats.AddValue(cand.Fitness)
	}

	// Notify observers with the population state
	data := evolve.PopulationData{
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
		o.PopulationUpdate(&data)
	}
	return &data
}

// shouldContinue determines whether or not the evolution should continue.
func shouldContinue(data *evolve.PopulationData, conds ...evolve.Condition) []evolve.Condition {
	satisfied := make([]evolve.Condition, 0)
	for _, cond := range conds {
		if cond.IsSatisfied(data) {
			satisfied = append(satisfied, cond)
		}
	}
	if len(satisfied) == 0 {
		return nil
	}
	return satisfied
}
