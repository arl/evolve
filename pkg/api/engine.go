package api

import (
	"math/rand"
	"runtime"
	"time"

	"github.com/aurelien-rainone/evolve/worker"
	"github.com/pkg/errors"
)

// Epocher is the interface implemented by objects having an Epoch method.
type Epocher interface {

	// Epoch performs one epoch (i.e generation) of the evolutionary process.
	//
	// It takes as argument the population to evolve in that step, the elitism
	// count -that is how many of the fittest candidates are preserved and
	// directly inserted into the nexct generation, without selection- and a
	// source of randomess.
	//
	// It returns the next generation.
	Epoch(Population, int, *rand.Rand) Population
}

/*
// Engine is the interface implemented by objects that provide evolution
// operations.
// TODO: does the Engine interface really needs all this methods? wouldn't one
// suffice and the others be derived from it (in engine.Base)?
// TODO: Could AddObserver/RemoveObserver be made an external interface,
// embedded in Engine? How would that go with future island observers?
type Engine interface {

	// Evolve executes the evolutionary algorithm until one of the termination
	// conditions is met, then return the fittest candidate from the final
	// generation.
	//
	// To return the entire population rather than just the fittest candidate,
	// use the EvolvePopulation method instead.
	//
	// size is the number of candidate solutions present in the population at
	// any point in time.
	// nelites is the number of candidates preserved via elitism. In elitism, a
	// sub-set of the population with the best fitness scores are preserved
	// unchanged in the subsequent generation. Candidate solutions that are
	// preserved unchanged through elitism remain eligible for selection for
	// breeding the remainder of the next generation. This value must be
	// non-negative and less than the population size. A value of zero means
	// that no elitism will be applied.
	// conds is a slice of conditions that may cause the evolution to terminate.
	//
	// Returns the fittest solution found by the evolutionary process.
	Evolve(size, nelites int, conds ...TerminationCondition) interface{}

	// EvolveWithSeedCandidates executes the evolutionary algorithm until one of
	// the termination conditions is met, then return the fittest candidate from
	// the final generation. Provide a set of candidates to seed the starting
	// population with.
	//
	// To return the entire population rather than just the fittest candidate,
	// use the EvolvePopulationWithSeedCandidates method instead.
	//
	// size is the number of candidate solutions present in the population at
	// any point in time.
	// nelites is the number of candidates preserved via elitism. In elitism, a
	// sub-set of the population with the best fitness scores are preserved
	// unchanged in the subsequent generation. Candidate solutions that are
	// preserved unchanged through elitism remain eligible for selection for
	// breeding the remainder of the next generation.  This value must be
	// non-negative and less than the population size. A value of zero means
	// that no elitism will be applied.
	// seedcands is a set of candidates to seed the population with. The size of
	// this collection must be no greater than the specified population size.
	// conds is a slice of conditions that may cause the evolution to terminate.
	//
	// Returns the fittest solution found by the evolutionary process.
	EvolveWithSeedCandidates(size, nelites int, seedcands []interface{},
		conds ...TerminationCondition) interface{}

	// EvolvePopulation executes the evolutionary algorithm until one of the
	// termination conditions is met, then return all of the candidates from the
	// final generation.
	//
	// To return just the fittest candidate rather than the entire population,
	// use the Evolve method instead.
	// size is the number of candidate solutions present in the population at
	// any point in time.
	// nelites is the number of candidates preserved via elitism. In elitism, a
	// sub-set of the population with the best fitness scores are preserved
	// unchanged in the subsequent generation. Candidate solutions that are
	// preserved unchanged through elitism remain eligible for selection for
	// breeding the remainder of the next generation.  This value must be
	// non-negative and less than the population size. A value of zero means
	// that no elitism will be applied.
	// conds is a slice of conditions that may cause the evolution to terminate.
	//
	// Returns the fittest solution found by the evolutionary process.
	EvolvePopulation(size, nelites int, conds ...TerminationCondition) Population

	// EvolvePopulationWithSeedCandidates executes the evolutionary algorithm
	// until one of the termination conditions is met, then return all of the
	// candidates from the final generation.
	//
	// To return just the fittest candidate rather than the entire population,
	// use the EvolveWithSeedCandidates method instead.
	// size is the number of candidate solutions present in the population at
	// any point in time.
	// nelites The number of candidates preserved via elitism. In elitism, a
	// sub-set of the population with the best fitness scores are preserved
	// unchanged in the subsequent generation. Candidate solutions that are
	// preserved unchanged through elitism remain eligible for selection for
	// breeding the remainder of the next generation. This value must be
	// non-negative and less than the population size. A value of zero means
	// that no elitism will be applied.
	// seedcands is a set of candidates to seed the population with.The size of
	// this collection must be no greater than the specified population size.
	// conditions One or more conditions that may cause the evolution to
	// terminate.
	//
	// Returns the fittest solution found by the evolutionary process.
	EvolvePopulationWithSeedCandidates(size, nelites int, seedcands []interface{},
		conds ...TerminationCondition) Population

	// AddObserver registers an observer to receive status updates on the
	// evolution progress.
	AddObserver(o Observer)

	// RemoveObserver removes an evolution observer.
	RemoveObserver(o Observer)

	// SatisfiedTerminationConditions returns a slice of all
	// TerminationCondition's that are satisfied by the current state of the
	// evolution engine.
	//
	// Usually this list will contain only one item, but it is possible that
	// multiple termination conditions will become satisfied at the same time.
	// In this case the condition objects in the list will be in the same order
	// that they were specified when passed to the engine.
	//
	// If the evolution has not yet terminated (either because it is still in
	// progress or because it hasn't even been started) then an
	// IllegalStateException will be thrown.
	//
	// If the evolution terminated because the request thread was interrupted
	// before any termination conditions were satisfied then this method will
	// return an empty list.
	//
	// Returns a list of statisfied conditions. The list is guaranteed to be
	// non-null. The list may be empty because it is possible for evolution to
	// terminate without any conditions being matched. The only situation in
	// which this occurs is when the request goroutine is interrupted.
	//
	// May return ErrIllegalState if this method is invoked on an evolution
	// engine before evolution is started or while it is still in progress.
	// TODO: find shorter name 'SatisfiedConditions' ?
	SatisfiedTerminationConditions() ([]TerminationCondition, error)
}*/

// Engine bla
// TODO: documentation
type Engine struct {
	pool           *worker.Pool // shared concurrent worker
	obs            map[Observer]struct{}
	rng            *rand.Rand
	gen            Generator
	eval           Evaluator
	singleThreaded bool
	satisfied      []TerminationCondition
	Epocher
}

// NewEngine creates a new evolution engine, injecting the various components
// required by an evolutionary algorithm.
//
// gen is the generator used to create the initial population that is
// iteratively evolved.
// eval evaluates fitness scores of candidate solutions.
// rng is the source of randomness used by all stochastic processes.
// ep is the Epocher
func NewEngine(gen Generator, eval Evaluator, rng *rand.Rand, ep Epocher) *Engine {
	return &Engine{
		gen:     gen,
		eval:    eval,
		rng:     rng,
		obs:     make(map[Observer]struct{}),
		Epocher: ep,
	}
}

// Evolve runs the evolutionary algorithm until one of the termination
// conditions is met, then return the fittest candidate from the final
// generation.
//
// To return the entire population rather than just the fittest candidate,
// use the EvolvePopulation method instead.
//
// size is the number of candidate solutions present in the population at any
// point in time.
// nelites is the number of candidates preserved via elitism. In elitism, a
// sub-set of the population with the best fitness scores are preserved
// unchanged in the subsequent generation. Candidate solutions that are
// preserved unchanged through elitism remain eligible for selection for
// breeding the remainder of the next generation. This value must be
// non-negative and less than the population size. A value of zero means that no
// elitism will be applied.
// conds is a slice of conditions that may cause the evolution to terminate.
//
// Returns the fittest solution found by the evolutionary process.
func (e *Engine) Evolve(size, nelites int, conds ...TerminationCondition) interface{} {
	return e.EvolveWithSeedCandidates(size, nelites, []interface{}{}, conds...)
}

// EvolveWithSeedCandidates runs the evolutionary algorithm until one of the
// termination conditions is met, then return the fittest candidate from the
// final generation. Provide a set of candidates to seed the starting population
// with.
//
// To return the entire population rather than just the fittest candidate,
// use the EvolvePopulationWithSeedCandidates method instead.
// size is the number of candidate solutions present in the population at any
// point in time.
// nelites is the number of candidates preserved via elitism. In elitism, a
// sub-set of the population with the best fitness scores are preserved
// unchanged in the subsequent generation. Candidate solutions that are
// preserved unchanged through elitism remain eligible for selection for
// breeding the remainder of the next generation. This value must be
// non-negative and less than the population size. A value of zero means that no
// elitism will be applied.
// seeds is a set of candidates to seed the population with. The size of this
// collection must be no greater than the specified population size.
// conds is a slice of conditions that may cause the evolution to terminate.
//
// Returns the fittest solution found by the evolutionary process.
func (e *Engine) EvolveWithSeedCandidates(size, nelites int, seeds []interface{}, conds ...TerminationCondition) interface{} {
	return e.EvolvePopulationWithSeedCandidates(size, nelites, seeds, conds...)[0].Candidate
}

// EvolvePopulation runs the evolutionary algorithm until one of the termination
// conditions is met, then return all of the candidates from the final
// generation.
//
// To return just the fittest candidate rather than the entire population,
// use the Evolve method instead.
// size is the number of candidate solutions present in the population at any
// point in time.
// nelites is the number of candidates preserved via elitism. In elitism, a
// sub-set of the population with the best fitness scores are preserved
// unchanged in the subsequent generation. Candidate solutions that are
// preserved unchanged through elitism remain eligible for selection for
// breeding the remainder of the next generation. This value must be
// non-negative and less than the population size. A value of zero means that no
// elitism will be applied.
// conds is a slice of conditions that may cause the evolution to terminate.
//
// Returns the fittest solution found by the evolutionary process.
func (e *Engine) EvolvePopulation(size, nelites int, conds ...TerminationCondition) Population {
	return e.EvolvePopulationWithSeedCandidates(size, nelites, []interface{}{}, conds...)
}

// EvolvePopulationWithSeedCandidates runs the evolutionary algorithm until one
// of the termination conditions is met, then return all of the candidates from
// the final generation.
//
// To return just the fittest candidate rather than the entire population, use
// the EvolveWithSeedCandidates method instead.
// size is the number of candidate solutions present in the population at any
// point in time.
// nelites The number of candidates preserved via elitism. In elitism, a sub-set
// of the population with the best fitness scores are preserved unchanged in the
// subsequent generation. Candidate solutions that are preserved unchanged
// through elitism remain eligible for selection for breeding the remainder of
// the next generation. This value must be non-negative and less than the
// population size. A value of zero means that no elitism will be applied.
// seeds is a set of candidates to seed the population with. The size of this
// collection must be no greater than the specified population size.
// conds is a slice of conditions that may cause the evolution to terminate.
//
// Returns the fittest solution found by the evolutionary process.
func (e *Engine) EvolvePopulationWithSeedCandidates(size, nelites int, seeds []interface{}, conds ...TerminationCondition) Population {

	if nelites < 0 || nelites >= size {
		panic("Elite count must be non-negative and less than population size.")
	}
	if len(conds) == 0 {
		panic("At least one TerminationCondition must be specified.")
	}

	e.satisfied = nil
	var curgen int
	startTime := time.Now()

	pop, err := SeedPopulation(e.gen, size, seeds, e.rng)
	// TODO: this method should return an error
	if err != nil {
		panic(err)
	}

	// Calculate the fitness scores for each member of the initial population.
	evpop := EvaluatePopulation(pop, e.eval, true)

	SortEvaluatedPopulation(evpop, e.eval.IsNatural())
	data := ComputePopulationData(evpop, e.eval.IsNatural(), nelites, curgen, startTime)

	// Notify observers of the state of the population.
	e.notifyPopulationChange(data)

	satisfied := ShouldContinue(data, conds...)
	for satisfied == nil {
		curgen++
		evpop = e.Epoch(evpop, nelites, e.rng)
		SortEvaluatedPopulation(evpop, e.eval.IsNatural())
		data = ComputePopulationData(evpop, e.eval.IsNatural(), nelites, curgen, startTime)
		// Notify observers of the state of the population.
		e.notifyPopulationChange(data)
		satisfied = ShouldContinue(data, conds...)
	}
	e.satisfied = satisfied
	return evpop
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
// ErrIllegalState is returned.
//
// If the evolution terminated because the request thread was interrupted before
// any termination conditions were satisfied then this method will return an
// empty slice.
//
// The slice is guaranteed to be non-null. The slice may be empty because it is
// possible for evolution to terminate without any conditions being matched.
// The only situation in which this occurs is when the request thread is
// interrupted.
func (e *Engine) SatisfiedTerminationConditions() ([]TerminationCondition, error) {
	if e.satisfied == nil {
		return nil, errors.Wrap(ErrIllegalState, "evolution engine has not terminated")
	}
	conds := make([]TerminationCondition, len(e.satisfied))
	copy(conds, e.satisfied)
	return conds, nil
}

// AddObserver adds a listener to receive status updates on the
// evolution progress.
//
// Updates are dispatched synchronously on the request thread. Observers should
// complete their processing and return in a timely manner to avoid holding up
// the evolution.
func (e *Engine) AddObserver(observer Observer) {
	e.obs[observer] = struct{}{}
}

// RemoveObserver removes an evolution progress listener.
func (e *Engine) RemoveObserver(observer Observer) {
	delete(e.obs, observer)
}

// notifyPopulationChange sends the population data to all registered observers.
func (e *Engine) notifyPopulationChange(data *PopulationData) {
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
func (e *Engine) SetSingleThreaded(singleThreaded bool) {
	e.singleThreaded = singleThreaded
}

// workerPool lazily creates the fitness evaluations goroutine pool.
func (e *Engine) workerPool() *worker.Pool {
	if e.pool == nil {
		// create a worker pool and set the maximum number of concurrent
		// goroutines to the number of logical CPUs usable by the current
		// process.
		e.pool = worker.NewPool(runtime.NumCPU())
	}
	return e.pool
}
