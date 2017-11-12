package evolve

import (
	"sort"
	"time"

	"github.com/aurelien-rainone/evolve/framework"
)

// Utility methods used by different evolution implementations. This class
// exists to avoid duplication of this logic among multiple evolution
// implementations.

// ShouldContinue returns a list of satisfied termination conditions if the
// evolution has reached some pre-specified state, an empty list if the
// evolution should stop because of a thread interruption, or null if the
// evolution should continue.
//
// Given data about the current population and a set of termination conditions,
// determines whether or not the evolution should continue.
// data is the current state of the population.
// conditions represents one or more termination conditions. The evolution
// should not continue if any of these is satisfied.
func ShouldContinue(
	data *framework.PopulationData,
	conditions ...framework.TerminationCondition) []framework.TerminationCondition {

	// If the thread has been interrupted, we should abort and return whatever
	// result we currently have.
	// TODO: ? what is this????
	//if (Thread.currentThread().isInterrupted()) {
	//return Collections.emptyList();
	//}
	//// Otherwise check the termination conditions for the evolution.
	satisfiedConditions := make([]framework.TerminationCondition, 0)
	for _, condition := range conditions {
		if condition.ShouldTerminate(data) {
			satisfiedConditions = append(satisfiedConditions, condition)
		}
	}
	if len(satisfiedConditions) == 0 {
		return nil

	}
	return satisfiedConditions
}

// SortEvaluatedPopulation sorts an evaluated population in descending order of
// fitness (descending order of fitness score for natural scores, ascending
// order of scores for non-natural scores).
func SortEvaluatedPopulation(evaluatedPopulation framework.EvaluatedPopulation, naturalFitness bool) {
	// Sort candidates in descending order according to fitness.
	if naturalFitness {
		// Descending values for natural fitness.
		sort.Sort(sort.Reverse(evaluatedPopulation))
	} else {
		// Ascending values for non-natural fitness.
		sort.Sort(evaluatedPopulation)
	}
}

// ComputePopulationData computes statistics about the current generation of
// evolved individuals, including the fittest candidate.
//
// evaluatedPopulation is the population of candidate solutions with their
// associated fitness scores.
// naturalFitness should be true if higher fitness scores mean fitter
// individuals, false otherwise.
// eliteCount is the number of candidates preserved via elitism.
// iterationNumber is the zero-based index of the current generation/epoch.
// startTime is the time at which the evolution began.
func ComputePopulationData(
	evaluatedPopulation framework.EvaluatedPopulation,
	naturalFitness bool,
	eliteCount int,
	iterationNumber int,
	startTime time.Time) *framework.PopulationData {

	stats := framework.NewDataSet(framework.WithInitialCapacity(len(evaluatedPopulation)))
	for _, candidate := range evaluatedPopulation {
		stats.AddValue(candidate.Fitness())
	}
	return framework.NewPopulationData(evaluatedPopulation[0].Candidate(),
		evaluatedPopulation[0].Fitness(),
		stats.ArithmeticMean(),
		stats.StandardDeviation(),
		naturalFitness,
		stats.Len(),
		eliteCount,
		iterationNumber,
		time.Since(startTime))
}
