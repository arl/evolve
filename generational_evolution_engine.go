package evolve

import (
	"math/rand"

	"github.com/aurelien-rainone/evolve/pkg/api"
)

// GenerationalEvolutionEngine implements a general-purpose generational
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
type GenerationalEvolutionEngine struct {
	op   api.Operator
	eval api.FitnessEvaluator
	sel  api.SelectionStrategy
	eng  *BaseEngine
}

// NewGenerationalEvolutionEngine creates a new evolution engine by specifying
// the various components required by a generational evolutionary algorithm.
//
// f is the factory used to create the initial population that is iteratively
// evolved.
// op is the combination of evolutionary operators used to evolve the population
// at each generation.
// eval is a function for assigning fitness scores to candidate solutions.
// sel is a strategy for selecting which candidates survive to be evolved.
// rng is the source of randomness used by all stochastic processes (including
// evolutionary operators and selection strategies).
func NewGenerationalEvolutionEngine(
	f api.Factory,
	op api.Operator,
	eval api.FitnessEvaluator,
	sel api.SelectionStrategy,
	rng *rand.Rand) *BaseEngine {

	// create the Stepper implementation
	stepper := &GenerationalEvolutionEngine{
		op:   op,
		eval: eval,
		sel:  sel,
	}

	// create the evolution engine implementation
	impl := NewBaseEngine(f, eval, rng, stepper)

	// provide the engine to the stepper for forwarding
	stepper.eng = impl
	return impl
}

// Step performs a single step/iteration of the evolutionary process.
//
// evpop is the population at the beginning of the process.
// nelites is the number of the fittest individuals that must be preserved.
//
// Returns the updated population after the evolutionary process has proceeded
// by one step/iteration.
func (e *GenerationalEvolutionEngine) Step(evpop api.EvaluatedPopulation, nelites int, rng *rand.Rand) api.EvaluatedPopulation {

	pop := make([]api.Candidate, 0, len(evpop))

	// First perform any elitist selection.
	elite := make([]api.Candidate, nelites)
	for i := 0; i < nelites; i++ {
		elite[i] = evpop[i].Candidate()
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
