package islands

import (
	"fmt"
	"math/rand"
	"runtime"
	"time"

	"github.com/aurelien-rainone/evolve"
	"github.com/aurelien-rainone/evolve/framework"
	"github.com/aurelien-rainone/evolve/termination"
	"github.com/aurelien-rainone/evolve/worker"
)

// IslandEvolution is an implementation of island evolution in which multiple
// independent populations are evolved in parallel with periodic migration of
// individuals between islands.
type IslandEvolution struct {
	islands                        []framework.EvolutionEngine
	migration                      Migration
	naturalFitness                 bool
	rng                            *rand.Rand
	observers                      map[IslandEvolutionObserver]struct{}
	satisfiedTerminationConditions []framework.TerminationCondition
}

// NewIslandEvolution creates an island system with the specified number of
// identically-configured islands.
//
// If you want more fine-grained control over the configuration of each island,
// use NewIslandEvolutionWithPreconfiguredIslands function, which accepts a list
// of pre-created islands (each is an instance of EvolutionEngine).
//
// - islandCount is the number of separate islands that will be part of the
// system.
// - migration is a migration strategy for moving individuals between islands at
// the end of an epoch.
// - candidateFactory generates the initial population for each island.
// - evolutionScheme is the evolutionary operator, or combination of
// evolutionary operators, used on each island.
// - fitnessEvaluator is the fitness function used on each island.
// - selectionStrategy is the selection strategy used on each island.
// - rng is a source of randomness, used by all islands.
func NewIslandEvolution(islandCount int,
	migration Migration,
	candidateFactory framework.CandidateFactory,
	evolutionScheme framework.EvolutionaryOperator,
	fitnessEvaluator framework.FitnessEvaluator,
	selectionStrategy framework.SelectionStrategy,
	rng *rand.Rand) *IslandEvolution {

	ie := NewIslandEvolutionWithPreconfiguredIslands(
		createIslands(
			islandCount,
			candidateFactory,
			evolutionScheme,
			fitnessEvaluator,
			selectionStrategy,
			rng),
		migration,
		fitnessEvaluator.IsNatural(),
		rng)
	return ie
}

type islandPopulationUpdater struct {
	observers   *map[IslandEvolutionObserver]struct{}
	islandIndex int
}

func (upd *islandPopulationUpdater) PopulationUpdate(data *framework.PopulationData) {
	for islandObserver := range *upd.observers {
		islandObserver.IslandPopulationUpdate(upd.islandIndex, data)
	}
}

// NewIslandEvolutionWithPreconfiguredIslands creates an island evolution system
// from a list of pre-configured islands.
//
// This function gives more control over the configuration of individual islands
// than the alternative constructor. The other construction function should be
// used where possible to avoid having to explicitly create each island.
//
// - islands is a list of pre-configured islands.
// - migration is a migration strategy for moving individuals between islands at
// the end of an epoch.
// - naturalFitness indicates, if true, that higher fitness values mean fitter
// individuals. If false, indicates that fitter individuals will have lower
// scores.
// - rng A source of randomness, used by all islands.
func NewIslandEvolutionWithPreconfiguredIslands(
	islands []framework.EvolutionEngine,
	migration Migration,
	naturalFitness bool,
	rng *rand.Rand) *IslandEvolution {

	ie := &IslandEvolution{
		islands:        islands,
		migration:      migration,
		naturalFitness: naturalFitness,
		rng:            rng,
		observers:      make(map[IslandEvolutionObserver]struct{}),
	}

	for i, island := range islands {
		island.AddEvolutionObserver(
			&islandPopulationUpdater{
				observers:   &ie.observers,
				islandIndex: i,
			})
	}
	return ie
}

// createIslands is an helper method used by NewIslandEvolution to create the
// individual islands if they haven't been provided already (via
// NewIslandEvolutionWithPreconfiguredIslands).
func createIslands(islandCount int,
	candidateFactory framework.CandidateFactory,
	evolutionScheme framework.EvolutionaryOperator,
	fitnessEvaluator framework.FitnessEvaluator,
	selectionStrategy framework.SelectionStrategy,
	rng *rand.Rand) []framework.EvolutionEngine {

	islands := make([]framework.EvolutionEngine, islandCount)
	for i := 0; i < islandCount; i++ {
		island := evolve.NewGenerationalEvolutionEngine(
			candidateFactory,
			evolutionScheme,
			fitnessEvaluator,
			selectionStrategy,
			rng)

		// Don't need fine-grained concurrency when
		// each island is on a separate thread.
		island.SetSingleThreaded(true)
		islands[i] = island
	}
	return islands
}

// Evolve starts the evolutionary process on each island and return the fittest
// candidate so far at the point any of the termination conditions is satisfied.
//
// If you interrupt the calling goroutine before this method returns, the method
// will return prematurely (with the best individual found so far).
//
// TODO: REWRITE THIS PART OF THE OCUMENTATION
// After returning in this way, the current goroutine's interrupted flag will be
// set.  It is preferable to use an appropritate framework.TerminationCondition
// rather than interrupting the evolution in this way.
//
// - populationSize is the population size "for each island". Therefore, if you
// have 5 islands, setting this parameter to 200 will result in 1000 individuals
// overall, 200 on each island.
// - eliteCount is the number of candidates preserved via elitism "on each
// island".  In elitism, a sub-set of the population with the best fitness
// scores are preserved unchanged in the subsequent generation. Candidate
// solutions that are preserved unchanged through elitism remain eligible for
// selection for breeding the remainder of the next generation.  This value must
// be non-negative and less than the population size. A value of zero means that
// no elitism will be applied.
// - epochLength is the number of generations that make up an epoch. Islands
// evolve independently for this number of generations and then migration occurs
// at the end of the epoch and the next epoch starts.
// - migrantCount is the number of individuals that will be migrated from each
// island at the end of each
// epoch.
// - conditions are one or more conditions that may cause the evolution to
// terminate.
//
// Returns the fittest solution found by the evolutionary process on any of the
// islands.
func (ie *IslandEvolution) Evolve(populationSize int,
	eliteCount int,
	epochLength int,
	migrantCount int,
	conditions ...framework.TerminationCondition) framework.Candidate {

	//threadPool := Executors.newFixedThreadPool(len(ie.islands))
	islandPopulations := make([][]framework.Candidate, len(ie.islands))
	evaluatedCombinedPopulation := make(framework.EvaluatedPopulation, 0)

	var (
		data                *framework.PopulationData
		satisfiedConditions []framework.TerminationCondition
		currentEpochIndex   int
		startTime           = time.Now()
	)

	for satisfiedConditions == nil {
		islandEpochs := ie.createEpochTasks(populationSize,
			eliteCount,
			epochLength,
			islandPopulations)
		defer func() {
			// TODO: catch statement here
			// done waitgroup or channel
		}()
		pool := worker.NewPool(runtime.NumCPU())

		results, err := pool.Submit(islandEpochs)
		if err != nil {
			panic(fmt.Sprintf("errors during island fitness evaluation island: %v", err))
		}
		evaluatedCombinedPopulation = nil
		evaluatedPopulations := make([]framework.EvaluatedPopulation, len(ie.islands))

		for i, result := range results {
			evaluatedIslandPopulation, ok := result.(framework.EvaluatedPopulation)
			if !ok {
				panic(fmt.Sprintf("result is not of the expected type, got %T", evaluatedIslandPopulation))
			}
			//evaluatedCombinedPopulation.addAll(evaluatedIslandPopulation)
			evaluatedCombinedPopulation = append(evaluatedCombinedPopulation, evaluatedIslandPopulation...)

			evaluatedPopulations[i] = evaluatedIslandPopulation
		}

		ie.migration.Migrate(evaluatedPopulations, migrantCount, ie.rng)

		evolve.SortEvaluatedPopulation(evaluatedCombinedPopulation, ie.naturalFitness)
		data = evolve.ComputePopulationData(evaluatedCombinedPopulation,
			ie.naturalFitness,
			eliteCount,
			currentEpochIndex,
			startTime)
		ie.notifyPopulationChange(data)

		// TODO 2 dimensions clear
		//islandPopulations.clear();
		//islandPopulations = nil
		for i, evaluatedPopulation := range evaluatedPopulations {
			islandPopulations[i] = ie.candidateSlice(evaluatedPopulation)
		}
		currentEpochIndex++
		//}
		//catch (InterruptedException ex)
		//{
		//Thread.currentThread().interrupt();
		//}
		//catch (ExecutionException ex)
		//{
		//throw new IllegalStateException(ex);
		//}
		satisfiedConditions = evolve.ShouldContinue(data, conditions...)
	}
	//threadPool.shutdownNow()

	ie.satisfiedTerminationConditions = satisfiedConditions
	return evaluatedCombinedPopulation[0].Candidate()
}

// Create the concurrently-executed tasks that perform evolution on each island.
func (ie *IslandEvolution) createEpochTasks(
	populationSize,
	eliteCount,
	epochLength int,
	islandPopulations [][]framework.Candidate) []worker.Worker {

	islandEpochs := make([]worker.Worker, len(ie.islands))
	for i := 0; i < len(ie.islands); i++ {

		var pop []framework.Candidate
		if len(islandPopulations) == 0 {
			pop = make([]framework.Candidate, 0)
		} else {
			pop = islandPopulations[i]
		}

		islandEpochs[i] = newEpoch(
			ie.islands[i],
			populationSize,
			eliteCount,
			pop,
			termination.NewGenerationCount(epochLength),
		)
	}
	return islandEpochs
}

// candidateList converts a slice of framework.EvaluatedCandidate's into a
// simple list of candidates.
//
// - evaluatedCandidates is the population of candidate objects to relieve of
// their evaluation wrappers.
//
// Returns the candidates, stripped of their fitness scores.
func (ie *IslandEvolution) candidateSlice(evaluatedCandidates framework.EvaluatedPopulation) []framework.Candidate {
	candidates := make([]framework.Candidate, len(evaluatedCandidates))
	for i, evaluatedCandidate := range evaluatedCandidates {
		candidates[i] = evaluatedCandidate.Candidate()
	}
	return candidates
}

// SatisfiedTerminationConditions returns a list of all
// framework.TerminationCondition's that are satisfied by the current state of
// the island evolution.
//
// Usually this list will contain only one item, but it is possible that
// mutliple termination conditions will become satisfied at the same time. In
// this case the condition objects in the list will be in the same order that
// they were specified when passed to the engine.
//
// If the evolution has not yet terminated (either because it is still in
// progress or because it hasn't even been started) then an
// IllegalStateException will be thrown.
//
// If the evolution terminated because the request thread was interrupted before
// any termination conditions were satisfied then this method will return an
// empty list.
//
// May return framework.ErrIllegalState if this method is invoked on an island
// system before evolution is started or while it is still in progress.
//
// Returns a list of statisfied conditions. The list is guaranteed to be
// non-null. The list may be empty because it is possible for evolution to
// terminate without any conditions being matched. The only situation in which
// this occurs is when the request thread is interrupted.
func (ie *IslandEvolution) SatisfiedTerminationConditions() ([]framework.TerminationCondition, error) {

	if ie.satisfiedTerminationConditions == nil {
		return nil, framework.ErrIllegalState("evolution engine has not terminated")
	}
	satisfiedTerminationConditions := make([]framework.TerminationCondition, len(ie.satisfiedTerminationConditions))
	copy(satisfiedTerminationConditions, ie.satisfiedTerminationConditions)
	return satisfiedTerminationConditions, nil
}

// AddEvolutionObserver adds an observer to the evolution.
//
// Observers will receives two types of updates: updates from each individual
// island at the end of each generation, and updates for the combined global
// population at the end of each epoch.
//
// Updates are dispatched synchronously on the request thread. Observers
// should complete their processing and return in a timely manner to avoid
// holding up the evolution.
func (ie *IslandEvolution) AddEvolutionObserver(observer IslandEvolutionObserver) {
	ie.observers[observer] = struct{}{}
}

// RemoveEvolutionObserver removes the specified observer.
func (ie *IslandEvolution) RemoveEvolutionObserver(observer IslandEvolutionObserver) {
	delete(ie.observers, observer)
}

// notifyPopulationChange send the population data to all registered
// observers.
func (ie *IslandEvolution) notifyPopulationChange(data *framework.PopulationData) {
	for observer := range ie.observers {
		observer.PopulationUpdate(data)
	}
}
