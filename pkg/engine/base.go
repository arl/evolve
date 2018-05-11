package engine

import (
	"fmt"
	"math/rand"
	"runtime"
	"time"

	"github.com/aurelien-rainone/evolve/pkg/api"
	"github.com/aurelien-rainone/evolve/worker"
	"github.com/pkg/errors"
)

// Stepper is the interface implemented by objects having a NextEvolutionStep
// method.
type Stepper interface {

	// Step performs a single step/iteration of the evolutionary process.
	//
	// evpop is the population at the beginning of the process.
	// nelites is the number of the fittest individuals that must be preserved.
	//
	// Returns the updated population after the evolutionary process has
	// proceeded by one step/iteration.
	Step(evpop api.EvaluatedPopulation, nelites int, rng *rand.Rand) api.EvaluatedPopulation
}

// Base is a base struct for EvolutionEngine implementations.
type Base struct {
	pool           *worker.Pool // shared concurrent worker
	obs            map[api.Observer]struct{}
	rng            *rand.Rand
	f              api.Factory
	eval           api.Evaluator
	singleThreaded bool
	satisfied      []api.TerminationCondition
	Stepper
}

// NewBase creates a new evolution engine by specifying the various
// components required by an evolutionary algorithm.
//
// candidateFactory is the factory used to create the initial population that is
// iteratively evolved.
// fitnessEvaluator is a function for assigning fitness scores to candidate
// solutions.
// rng is the source of randomness used by all stochastic processes (including
// evolutionary operators and selection strategies).
func NewBase(f api.Factory, eval api.Evaluator, rng *rand.Rand, stepper Stepper) *Base {
	return &Base{
		f:       f,
		eval:    eval,
		rng:     rng,
		obs:     make(map[api.Observer]struct{}),
		Stepper: stepper,
	}
}

// Evolve executes the evolutionary algorithm until one of the termination
// conditions is met, then return the fittest candidate from the final
// generation.
//
// To return the entire population rather than just the fittest candidate,
// use the EvolvePopulation method instead.
//
// populationSize is the number of candidate solutions present in the population
// at any point in time.
// eliteCount is the number of candidates preserved via elitism. In elitism, a
// sub-set of the population with the best fitness scores are preserved
// unchanged in the subsequent generation. Candidate solutions that are
// preserved unchanged through elitism remain eligible for selection for
// breeding the remainder of the next generation. This value must be
// non-negative and less than the population size. A value of zero means that no
// elitism will be applied.
// conditions is a slice of conditions that may cause the evolution to
// terminate.
//
// Returns the fittest solution found by the evolutionary process.
func (e *Base) Evolve(size, nelites int, conds ...api.TerminationCondition) interface{} {
	return e.EvolveWithSeedCandidates(size, nelites, []interface{}{}, conds...)
}

// EvolveWithSeedCandidates executes the evolutionary algorithm until one of
// the termination conditions is met, then return the fittest candidate from
// the final generation. Provide a set of candidates to seed the starting
// population with.
//
// To return the entire population rather than just the fittest candidate,
// use the EvolvePopulationWithSeedCandidates method instead.
// populationSize is the number of candidate solutions present in the
// population at any point in time.
// eliteCount is the number of candidates preserved via elitism. In elitism, a
// sub-set of the population with the best fitness scores are preserved
// unchanged in the subsequent generation. Candidate solutions that are
// preserved unchanged through elitism remain eligible for selection for
// breeding the remainder of the next generation.  This value must be
// non-negative and less than the population size. A value of zero means that no
// elitism will be applied.
// seedCandidates is a set of candidates to seed the population with. The size
// of this collection must be no greater than the specified population size.
// conditions is a slice of conditions that may cause the evolution to
// terminate.
//
// Returns the fittest solution found by the evolutionary process.
func (e *Base) EvolveWithSeedCandidates(size, nelites int, seedcands []interface{}, conds ...api.TerminationCondition) interface{} {
	return e.EvolvePopulationWithSeedCandidates(size, nelites, seedcands, conds...)[0].Candidate()
}

// EvolvePopulation executes the evolutionary algorithm until one of the
// termination conditions is met, then return all of the candidates from the
// final generation.
//
// To return just the fittest candidate rather than the entire population,
// use the Evolve method instead.
// populationSize is the number of candidate solutions present in the population
// at any point in time.
// eliteCount is the number of candidates preserved via elitism. In elitism, a
// sub-set of the population with the best fitness scores are preserved
// unchanged in the subsequent generation. Candidate solutions that are
// preserved unchanged through elitism remain eligible for selection for
// breeding the remainder of the next generation.  This value must be
// non-negative and less than the population size. A value of zero means that no
// elitism will be applied.
// conditions is a slice of conditions that may cause the evolution to
// terminate.
//
// Returns the fittest solution found by the evolutionary process.
func (e *Base) EvolvePopulation(size, nelites int, conds ...api.TerminationCondition) api.EvaluatedPopulation {
	return e.EvolvePopulationWithSeedCandidates(size, nelites, []interface{}{}, conds...)
}

// EvolvePopulationWithSeedCandidates executes the evolutionary algorithm
// until one of the termination conditions is met, then return all of the
// candidates from the final generation.
//
// To return just the fittest candidate rather than the entire population,
// use the EvolveWithSeedCandidates method instead.
// populationSize is the number of candidate solutions present in the population
// at any point in time.
// eliteCount The number of candidates preserved via elitism.  In elitism, a
// sub-set of the population with the best fitness scores are preserved
// unchanged in the subsequent generation.  Candidate solutions that are
// preserved unchanged through elitism remain eligible for selection for
// breeding the remainder of the next generation.  This value must be
// non-negative and less than the population size.  A value of zero means that
// no elitism will be applied.
// seedCandidates A set of candidates to seed the population with. The size of
// this collection must be no greater than the specified population size.
// conditions One or more conditions that may cause the evolution to terminate.
//
// Returns the fittest solution found by the evolutionary process.
func (e *Base) EvolvePopulationWithSeedCandidates(size, nelites int, seedcands []interface{}, conds ...api.TerminationCondition) api.EvaluatedPopulation {

	if nelites < 0 || nelites >= size {
		panic("Elite count must be non-negative and less than population size.")
	}
	if len(conds) == 0 {
		panic("At least one TerminationCondition must be specified.")
	}

	e.satisfied = nil
	var curgen int
	startTime := time.Now()

	pop := e.f.SeedPopulation(size,
		seedcands,
		e.rng)

	// Calculate the fitness scores for each member of the initial population.
	evpop := e.evaluatePopulation(pop)

	api.SortEvaluatedPopulation(evpop, e.eval.IsNatural())
	data := api.ComputePopulationData(evpop, e.eval.IsNatural(), nelites, curgen, startTime)

	// Notify observers of the state of the population.
	e.notifyPopulationChange(data)

	satisfied := api.ShouldContinue(data, conds...)
	for satisfied == nil {
		curgen++
		evpop = e.Step(evpop, nelites, e.rng)
		api.SortEvaluatedPopulation(evpop, e.eval.IsNatural())
		data = api.ComputePopulationData(evpop, e.eval.IsNatural(), nelites, curgen, startTime)
		// Notify observers of the state of the population.
		e.notifyPopulationChange(data)
		satisfied = api.ShouldContinue(data, conds...)
	}
	e.satisfied = satisfied
	return evpop
}

// evaluatePopulation takes a population, assigns a fitness score to each member
// and returns the members with their scores attached, sorted in descending
// order of fitness (descending order of fitness score for natural scores,
// ascending order of scores for non-natural scores).
// population is the population to evaluate (each candidate is assigned a
// fitness score).
//
// Returns the evaluated population (a list of candidates with attached fitness
// scores).
func (e *Base) evaluatePopulation(
	pop []interface{}) api.EvaluatedPopulation {

	// Do fitness evaluations
	evpop := make(api.EvaluatedPopulation, len(pop))

	if e.singleThreaded {

		var err error
		for i, candidate := range pop {
			evpop[i], err = api.NewEvaluatedCandidate(candidate, e.eval.Fitness(candidate, pop))
			if err != nil {
				panic(fmt.Sprintf("Can't evaluate candidate %v: %v", candidate, err))
			}
		}

	} else {

		// Create a worker pool that will divides the required number of fitness
		// evaluations equally among the available goroutines and coordinate
		// them so that we do not proceed until all of them have finished
		// processing.
		workers := make([]worker.Worker, len(pop))
		for i := range pop {
			workers[i] = &fitnessEvaluationWorker{
				idx:       i,
				pop:       pop,
				evaluator: e.eval,
			}
		}

		results, err := e.workerPool().Submit(workers)
		if err != nil {
			panic(fmt.Sprintf("Error while submitting workers to the pool: %v", err))
		}

		for i, result := range results {
			evpop[i] = result.(*api.EvaluatedCandidate)
		}
		// TODO: handle goroutine termination
		/*
		   catch (InterruptedException ex)
		   {
		       // Restore the interrupted status, allows methods further up the call-stack
		       // to abort processing if appropriate.
		       Thread.currentThread().interrupt();
		   }
		*/
	}

	return evpop
}

type fitnessEvaluationWorker struct {
	idx       int           // index of candidate to evaluate
	pop       []interface{} // full population
	evaluator api.Evaluator
}

func (w *fitnessEvaluationWorker) Work() (interface{}, error) {
	return api.NewEvaluatedCandidate(w.pop[w.idx],
		w.evaluator.Fitness(w.pop[w.idx], w.pop))
}

// SatisfiedTerminationConditions returns a slice of all TerminationCondition's
// that are satisfied by the current state of the evolution engine.
//
// Usually this slice will contain only one item, but it is possible that
// multiple termination conditions will become satisfied at the same time. In
// this case the condition objects in the slice will be in the same order that
// they were specified when passed to the engine.
//
// If the evolution has not yet terminated (either because it is still in
// progress or because it hasn't even been started) then
// api.ErrIllegalState is returned.
//
// If the evolution terminated because the request thread was interrupted before
// any termination conditions were satisfied then this method will return an
// empty slice.
//
// The slice is guaranteed to be non-null. The slice may be empty because it is
// possible for evolution to terminate without any conditions being matched.
// The only situation in which this occurs is when the request thread is
// interrupted.
func (e *Base) SatisfiedTerminationConditions() ([]api.TerminationCondition, error) {
	if e.satisfied == nil {
		return nil, errors.Wrap(api.ErrIllegalState, "evolution engine has not terminated")
	}
	conds := make([]api.TerminationCondition, len(e.satisfied))
	copy(conds, e.satisfied)
	return conds, nil
}

// AddObserver adds a listener to receive status updates on the
// evolution progress.
//
// Updates are dispatched synchronously on the request thread. Observers should
// complete their processing and return in a timely manner to avoid holding up
// the evolution.
func (e *Base) AddObserver(observer api.Observer) {
	e.obs[observer] = struct{}{}
}

// RemoveObserver removes an evolution progress listener.
func (e *Base) RemoveObserver(observer api.Observer) {
	delete(e.obs, observer)
}

// notifyPopulationChange sends the population data to all registered observers.
func (e *Base) notifyPopulationChange(data *api.PopulationData) {
	for observer := range e.obs {
		observer.PopulationUpdate(data)
	}
}

// SetSingleThreaded forces evaluation to occur synchronously on the request
// goroutine.
//
// By default, fitness evaluations are performed on separate goroutines (as many
// as there are available cores/processors). This is useful in restricted
// environments where programs are not permitted to start or control threads. It
// might also lead to better performance for programs that have extremely
// lightweight/trivial fitness evaluations.
func (e *Base) SetSingleThreaded(singleThreaded bool) {
	e.singleThreaded = singleThreaded
}

// workerPool lazily creates the fitness evaluations goroutine pool.
func (e *Base) workerPool() *worker.Pool {
	if e.pool == nil {
		// create a worker pool and set the maximum number of concurrent
		// goroutines to the number of logical CPUs usable by the current
		// process.
		e.pool = worker.NewPool(runtime.NumCPU())
	}
	return e.pool
}
