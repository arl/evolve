package base

// EvolutionEngine is the interface implented by objects provide evolution
// operations.
type EvolutionEngine interface {

	// Evolve executes the evolutionary algorithm until one of the termination
	// conditions is met, then return the fittest candidate from the final
	// generation.
	//
	// To return the entire population rather than just the fittest candidate,
	// use the EvolvePopulation method instead.
	//
	// - populationSize is the number of candidate solutions present in the
	// population at any point in time.
	// - eliteCount is the number of candidates preserved via elitism. In
	// elitism, a sub-set of the population with the best fitness scores are
	// preserved unchanged in the subsequent generation. Candidate solutions
	// that are preserved unchanged through elitism remain eligible for
	// selection for breeding the remainder of the next generation. This value
	// must be non-negative and less than the population size. A value of zero
	// means that no elitism will be applied.
	// - conditions is a slice of conditions that may cause the evolution to
	// terminate.
	//
	// Return the fittest solution found by the evolutionary process.
	Evolve(populationSize, eliteCount int, conditions []TerminationCondition) Candidate

	// EvolveWithSeedCandidates executes the evolutionary algorithm until one of
	// the termination conditions is met, then return the fittest candidate from
	// the final generation. Provide a set of candidates to seed the starting
	// population with.
	//
	// To return the entire population rather than just the fittest candidate,
	// use the EvolvePopulationWithSeedCandidates method instead.
	// - populationSize is the number of candidate solutions present in the
	// population at any point in time.
	// - eliteCount is the number of candidates preserved via elitism. In
	// elitism, a sub-set of the population with the best fitness scores are
	// preserved unchanged in the subsequent generation. Candidate solutions
	// that are preserved unchanged through elitism remain eligible for
	// selection for breeding the remainder of the next generation.  This value
	// must be non-negative and less than the population size. A value of zero
	// means that no elitism will be applied.
	// - seedCandidates is a set of candidates to seed the population with. The
	// size of this collection must be no greater than the specified population
	// size.
	// - conditions is a slice of conditions that may cause the evolution to
	// terminate.
	//
	// Returns the fittest solution found by the evolutionary process.
	EvolveWithSeedCandidates(populationSize, eliteCount int,
		seedCandidates []Candidate,
		conditions []TerminationCondition) Candidate

	// EvolvePopulation executes the evolutionary algorithm until one of the
	// termination conditions is met, then return all of the candidates from the
	// final generation.
	//
	// To return just the fittest candidate rather than the entire population,
	// use the Evolve method instead.
	// - populationSize is the number of candidate solutions present in the
	// population at any point in time.
	// - eliteCount is the number of candidates preserved via elitism. In
	// elitism, a sub-set of the population with the best fitness scores are
	// preserved unchanged in the subsequent generation. Candidate solutions
	// that are preserved unchanged through elitism remain eligible for
	// selection for breeding the remainder of the next generation.  This value
	// must be non-negative and less than the population size. A value of zero
	// means that no elitism will be applied.
	// -  conditions is a slice of conditions that may cause the evolution to
	// terminate.
	//
	// Return the fittest solution found by the evolutionary process.
	EvolvePopulation(populationSize, eliteCount int,
		conditions []TerminationCondition) []*EvaluatedCandidate

	// EvolvePopulationWithSeedCandidates executes the evolutionary algorithm
	// until one of the termination conditions is met, then return all of the
	// candidates from the final generation.
	//
	// To return just the fittest candidate rather than the entire population,
	// use the EvolveWithSeedCandidates method instead.
	// - populationSize is the number of candidate solutions present in the
	// population at any point in time.
	// - eliteCount The number of candidates preserved via elitism.  In elitism,
	// a sub-set of the population with the best fitness scores are preserved
	// unchanged in the subsequent generation.  Candidate solutions that are
	// preserved unchanged through elitism remain eligible for selection for
	// breeding the remainder of the next generation.  This value must be
	// non-negative and less than the population size.  A value of zero means
	// that no elitism will be applied.
	// - seedCandidates A set of candidates to seed the population with.  The
	// size of this collection must be no greater than the specified population
	// size.
	// - conditions One or more conditions that may cause the evolution to
	// terminate.
	//
	// Return the fittest solution found by the evolutionary process.
	EvolvePopulationWithSeedCandidates(populationSize, eliteCount int,
		seedCandidates []Candidate,
		conditions []TerminationCondition) []*EvaluatedCandidate

	// AddEvolutionObserver adds a listener to receive status updates on the
	// evolution progress.
	AddEvolutionObserver(observer EvolutionObserver)

	// RemoveEvolutionObserver removes an evolution progress listener.
	RemoveEvolutionObserver(observer EvolutionObserver)

	// SatisfiedTerminationConditions returns a list of all
	// TerminationCondition's that are satisfied by the current state of the
	// evolution engine.
	//
	// Usually this list will contain only one item, but it is possible that
	// mutliple termination conditions will become satisfied at the same time.
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
	// The error value will be ErrIllegalStateI if this method is invoked on an
	// evolution engine before evolution is started or while it is still in
	// progress.
	SatisfiedTerminationConditions() ([]TerminationCondition, error)
}

type ErrIllegalState error
