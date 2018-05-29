package engine

import (
	"math/rand"
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
	seed           int64
	obs            map[api.Observer]struct{}
	rng            *rand.Rand
	gen            api.Generator
	eval           api.Evaluator
	ep             api.Epocher
	singleThreaded bool
}

func New(g api.Generator, ev api.Evaluator, ep api.Epocher, options ...func(*Engine) error) (*Engine, error) {
	eng := Engine{
		seed: time.Now().UnixNano(),
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
		eng.rng = rand.New(random.NewMT19937(0))
	}
	eng.rng.Seed(eng.seed)
	return &eng, nil
}

func (e *Engine) AddObserver(o api.Observer) {
	e.obs[o] = struct{}{}
}

func (e *Engine) RemoveObserver(o api.Observer) {
	delete(e.obs, o)
}

// Rand sets the random number generator for the engine.
func Rand(rng *rand.Rand) func(*Engine) error {
	return func(eng *Engine) error {
		eng.rng = rng
		return nil
	}
}

// Seed sets the seed of the random number generator for the engine.
func Seed(s int64) func(*Engine) error {
	return func(eng *Engine) error {
		eng.seed = s
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

type config struct {
	nelites int
	seeds   []interface{}
	conds   []api.TerminationCondition
}

func Elites(n int) func(*config) error {
	return func(cfg *config) error {
		cfg.nelites = n
		return nil
	}
}

func Seeds(seeds []interface{}) func(*config) error {
	return func(cfg *config) error {
		cfg.seeds = seeds
		return nil
	}
}

func EndOn(cond api.TerminationCondition) func(*config) error {
	return func(cfg *config) error {
		cfg.conds = append(cfg.conds, cond)
		return nil
	}
}

// TODO: rename opts in options if accepted
func (e *Engine) Evolve(popsize int, options ...func(*config) error) (api.Population, []api.TerminationCondition, error) {
	var cfg config
	for _, opt := range options {
		if err := opt(&cfg); err != nil {
			return nil, nil, err
		}
	}

	if popsize <= 0 {
		return nil, nil, ErrEnginePopSize
	}
	if cfg.nelites < 0 || cfg.nelites >= popsize {
		return nil, nil, ErrEngineElite
	}
	if len(cfg.conds) == 0 {
		return nil, nil, ErrEngineTermination
	}

	var curgen int
	startTime := time.Now()

	pop, err := api.SeedPopulation(e.gen, popsize, cfg.seeds, e.rng)
	if err != nil {
		return nil, nil, errors.Wrap(err, "can't seed population")
	}

	// Calculate the fitness scores for each member of the initial population.
	evpop := api.EvaluatePopulation(pop, e.eval, true)

	api.SortEvaluatedPopulation(evpop, e.eval.IsNatural())
	data := api.ComputePopulationData(evpop, e.eval.IsNatural(), cfg.nelites, curgen, startTime)

	// Notify observers of the state of the population.
	for o := range e.obs {
		o.PopulationUpdate(data)
	}

	// TODO: all the functions in api like ShouldContinue,
	// SortEvaluatedPopulation, etc. really are useful there?
	// Why aren't they in engine package, unexported?
	satisfied := api.ShouldContinue(data, cfg.conds...)
	for satisfied == nil {
		curgen++
		evpop = e.ep.Epoch(evpop, cfg.nelites, e.rng)
		api.SortEvaluatedPopulation(evpop, e.eval.IsNatural())
		data = api.ComputePopulationData(evpop, e.eval.IsNatural(), cfg.nelites, curgen, startTime)
		for o := range e.obs {
			o.PopulationUpdate(data)
		}
		satisfied = api.ShouldContinue(data, cfg.conds...)
	}
	return evpop, satisfied, nil
}
