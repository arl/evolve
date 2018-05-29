package engine

import (
	"math/rand"
	"sort"
	"time"

	"github.com/aurelien-rainone/evolve/pkg/api"
	"github.com/aurelien-rainone/evolve/random"
	"github.com/pkg/errors"
)

var (
	ErrEnginePopSize     = errors.New("population size must be at least 1")
	ErrEngineElite       = errors.New("elite count must be positive and less than population size")
	ErrEngineTermination = errors.New("at least one termination condition must be specified")
)

type Engine struct {
	obs            map[api.Observer]struct{}
	rng            *rand.Rand
	gen            api.Generator
	eval           api.Evaluator
	ep             api.Epocher
	stats          *api.Dataset
	singleThreaded bool

	nelites int
	seeds   []interface{}
	conds   []api.TerminationCondition
	size    int
}

func New(g api.Generator, ev api.Evaluator, ep api.Epocher, options ...func(*Engine) error) (*Engine, error) {
	eng := Engine{
		obs:  make(map[api.Observer]struct{}),
		gen:  g,
		eval: ev,
		ep:   ep,
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
// be non-negative and less than the population size or Evolve will return
// ErrEngineElite
func Elites(n int) func(*Engine) error {
	return func(eng *Engine) error {
		if n < 0 || n >= eng.size {
			return ErrEngineElite
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
// with Seeds. size must be at least 1 or Evolve will return ErrEnginePopSize.
//
// At least one termination condition must be defined with EndOn, or Evolve will
// return ErrEngineTermination.
func (e *Engine) Evolve(popsize int, options ...func(*Engine) error) (api.Population, []api.TerminationCondition, error) {
	e.size = popsize
	for _, opt := range options {
		if err := opt(e); err != nil {
			return nil, nil, err
		}
	}

	if popsize <= 0 {
		return nil, nil, ErrEnginePopSize
	}
	if len(e.conds) == 0 {
		return nil, nil, ErrEngineTermination
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

	// Evaluate fitness of the whole initial population
	evpop := api.EvaluatePopulation(pop, e.eval, true)

	for {

		ngen++

		// Sort population according to fitness.
		if e.eval.IsNatural() {
			sort.Sort(sort.Reverse(evpop))
		} else {
			sort.Sort(evpop)
		}

		e.stats.Clear()
		for _, cand := range evpop {
			e.stats.AddValue(cand.Fitness)
		}

		// Notify observers with the population state
		data := api.PopulationData{
			BestCand:    evpop[0].Candidate,
			BestFitness: evpop[0].Fitness,
			Mean:        e.stats.ArithmeticMean(),
			StdDev:      e.stats.StandardDeviation(),
			Natural:     e.eval.IsNatural(),
			Size:        e.stats.Len(),
			NumElites:   e.nelites,
			GenNumber:   ngen,
			Elapsed:     time.Since(start),
		}
		for o := range e.obs {
			o.PopulationUpdate(&data)
		}

		// TODO: all the functions in api like ShouldContinue,
		// SortEvaluatedPopulation, etc. really are useful there?
		// Why aren't they in engine package, unexported?
		satisfied = shouldContinue(&data, e.conds...)
		if satisfied != nil {
			break
		}

		evpop = e.ep.Epoch(evpop, e.nelites, e.rng)
	}
	return evpop, satisfied, nil
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
