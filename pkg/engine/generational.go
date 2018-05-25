package engine

import (
	"math/rand"

	"github.com/aurelien-rainone/evolve/pkg/api"
)

// Generational implements a general-purpose engine for generational
// evolutionary algorithm.
//
// It supports optional concurrent fitness evaluations to take full advantage of
// multi-processor, multi-core and hyper-threaded machines through the
// concurrent evaluation of candidate's fitness.
//
// If multi-threading is enabled, evolution (mutation, crossover, etc.) occurs
// on the request goroutine but fitness evaluations are delegated to a pool of
// worker threads. All of the host's available processing units are used (i.e.
// on a quad-core machine there will be four fitness evaluation worker threads).
//
// If multi-threading is disabled, all work is performed synchronously on the
// request thread. This strategy is suitable for restricted/managed environments
// where it is not permitted for applications to manage their own threads. If
// there are no restrictions on concurrency, applications should enable
// multi-threading for improved performance.
type Generational struct {
	op   api.Operator
	eval api.Evaluator
	sel  api.Selection
	eng  *Base
}

// NewGenerational creates a new generational evolution engine, injecting the
// various components required by an evolutionary algorithm.
//
// gen is the generator used to create the initial population that is iteratively
// evolved.
// op is the evolutionary operator applied at each generation to evolve the
// population.
// eval evaluates fitness scores of candidate solutions.
// sel is a strategy for selecting which candidates survive an epoch.
// rng is the source of randomness used by all stochastic processes.
func NewGenerational(gen api.Generator, op api.Operator, eval api.Evaluator, sel api.Selection, rng *rand.Rand) *Base {

	// create the Epocher implementation
	ep := &Generational{op: op, eval: eval, sel: sel}

	// create the evolution engine implementation
	impl := NewBase(gen, eval, rng, ep)

	// provide the engine to the epocher for forwarding
	ep.eng = impl
	return impl
}

// Epoch performs a single step/iteration of the evolutionary process.
//
// evpop is the population at the beginning of the process.
// nelites is the number of the fittest individuals that must be preserved.
//
// Returns the updated population after the evolutionary process has proceeded
// by one step/iteration.
func (e *Generational) Epoch(evpop api.Population, nelites int, rng *rand.Rand) api.Population {

	pop := make([]interface{}, 0, len(evpop))

	// First perform any elitist selection.
	elite := make([]interface{}, nelites)
	for i := 0; i < nelites; i++ {
		elite[i] = evpop[i].Candidate
	}

	// Then select candidates that will be operated on to create the evolved
	// portion of the next generation.
	pop = append(pop, e.sel.Select(evpop,
		e.eval.IsNatural(),
		len(evpop)-nelites,
		rng)...)

	// Then evolve the population.
	pop = e.op.Apply(pop, rng)
	// When the evolution is finished, add the elite to the population.
	pop = append(pop, elite...)

	return e.eng.evaluatePopulation(pop)
}
