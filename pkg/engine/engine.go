package engine

import (
	"math/rand"
	"sort"
	"time"

	"github.com/arl/evolve/pkg/api"
	"github.com/arl/evolve/random"
	"github.com/pkg/errors"
)

// Engine runs an evolutionary algorithm, following all the steps of evolution,
// from the creation of the initial population to the end of evolution.
type Engine struct {
	obs     map[api.Observer]struct{}
	rng     *rand.Rand
	gen     api.Generator
	eval    api.Evaluator
	epoch   api.Epocher
	stats   *api.Dataset
	nelites int
	seeds   []interface{}
	conds   []api.TerminationCondition
	size    int
}

// New returns a new evolution Engine.
//
// gen generates new random candidates solutions.
// eval evaluates fitness scores of candidates.
// epoch transforms a whole population into the next generation.
// rng is the source of randomness.
func New(gen api.Generator, eval api.Evaluator, epoch api.Epocher, options ...func(*Engine) error) (*Engine, error) {
	eng := Engine{
		obs:   make(map[api.Observer]struct{}),
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
		eng.rng = rand.New(random.NewMT19937(seed))
	}
	return &eng, nil
}

// AddObserver adds an observer of the evolution process.
func (e *Engine) AddObserver(o api.Observer) {
	e.obs[o] = struct{}{}
}

// RemoveObserver removes an observer of the evolution process.
func (e *Engine) RemoveObserver(o api.Observer) {
	delete(e.obs, o)
}

// Rand sets the global random number generator for the engine.
func Rand(rng *rand.Rand) func(*Engine) error {
	return func(eng *Engine) error {
		eng.rng = rng
		return nil
	}
}

// Observer adds an observer of the evolution process for the engine.
func Observer(o api.Observer) func(*Engine) error {
	return func(eng *Engine) error {
		eng.obs[o] = struct{}{}
		return nil
	}
}

// Elites defines the number of candidates preserved via elitism. By default it
// is set to 0, no elitism is applied.
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

// Seeds provides a set of candidates to seed the starting population with.
// A second (and successive) call to Seeds will replace the set of seed
// candidates defined by the previous call.
func Seeds(seeds []interface{}) func(*Engine) error {
	return func(eng *Engine) error {
		eng.seeds = seeds
		return nil
	}
}

// EndOn adds a termination condition to the engine. The engine stops
// after one or more termination conditions are met.
func EndOn(cond api.TerminationCondition) func(*Engine) error {
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
func (e *Engine) Evolve(popsize int, options ...func(*Engine) error) (api.Population, []api.TerminationCondition, error) {
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
	e.stats = api.NewDataset(popsize)

	var ngen int
	start := time.Now()

	pop, err := api.SeedPopulation(e.gen, popsize, e.seeds, e.rng)
	if err != nil {
		return nil, nil, errors.Wrap(err, "can't seed population")
	}

	var satisfied []api.TerminationCondition

	// Evaluate initial population fitness
	evpop := api.EvaluatePopulation(pop, e.eval, true)

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

func (e *Engine) updateStats(pop api.Population, ngen int, elapsed time.Duration) *api.PopulationData {

	e.stats.Clear()
	for _, cand := range pop {
		e.stats.AddValue(cand.Fitness)
	}

	// Notify observers with the population state
	data := api.PopulationData{
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
func shouldContinue(data *api.PopulationData, conds ...api.TerminationCondition) []api.TerminationCondition {
	satisfied := make([]api.TerminationCondition, 0)
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
