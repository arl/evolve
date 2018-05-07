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
	evolutionScheme   api.EvolutionaryOperator
	fitnessEvaluator  api.FitnessEvaluator
	selectionStrategy api.SelectionStrategy
	engine            *AbstractEvolutionEngine
}

// NewGenerationalEvolutionEngine creates a new evolution engine by specifying
// the various components required by a generational evolutionary algorithm.
//
// candidateFactory is the factory used to create the initial population that is
// iteratively evolved.
// evolutionScheme is the combination of evolutionary operators used to evolve
// the population at each generation.
// fitnessEvaluator is a function for assigning fitness scores to candidate
// solutions.
// selectionStrategy is a strategy for selecting which candidates survive to be
// evolved.
// rng is the source of randomness used by all stochastic processes (including
// evolutionary operators and selection strategies).
func NewGenerationalEvolutionEngine(
	candidateFactory api.Factory,
	evolutionScheme api.EvolutionaryOperator,
	fitnessEvaluator api.FitnessEvaluator,
	selectionStrategy api.SelectionStrategy,
	rng *rand.Rand) *AbstractEvolutionEngine {

	// create the Stepper implementation
	stepper := &GenerationalEvolutionEngine{
		evolutionScheme:   evolutionScheme,
		fitnessEvaluator:  fitnessEvaluator,
		selectionStrategy: selectionStrategy,
	}

	// create the evolution engine implementation
	engineImpl := NewAbstractEvolutionEngine(
		candidateFactory,
		fitnessEvaluator,
		rng,
		stepper,
	)

	// provide the engine to the stepper for forwarding
	stepper.engine = engineImpl
	return engineImpl
}

// NextEvolutionStep performs a single step/iteration of the evolutionary process.
//
// evaluatedPopulation is the population at the beginning of the process.
// eliteCount is the number of the fittest individuals that must be preserved.
//
// Returns the updated population after the evolutionary process has proceeded
// by one step/iteration.
func (e *GenerationalEvolutionEngine) NextEvolutionStep(
	evaluatedPopulation api.EvaluatedPopulation,
	eliteCount int,
	rng *rand.Rand) api.EvaluatedPopulation {

	population := make([]api.Candidate, 0, len(evaluatedPopulation))

	// First perform any elitist selection.
	elite := make([]api.Candidate, eliteCount)
	for i := 0; i < eliteCount; i++ {
		elite[i] = evaluatedPopulation[i].Candidate()
	}

	// Then select candidates that will be operated on to create the evolved
	// portion of the next generation.
	population = append(population, e.selectionStrategy.Select(evaluatedPopulation,
		e.fitnessEvaluator.IsNatural(),
		len(evaluatedPopulation)-eliteCount,
		rng)...)

	// Then evolve the population.
	population = e.evolutionScheme.Apply(population, rng)
	// When the evolution is finished, add the elite to the population.
	population = append(population, elite...)

	return e.engine.evaluatePopulation(population)
}
