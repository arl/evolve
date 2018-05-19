package api

import (
	"sort"
	"time"
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
	data *PopulationData,
	conds ...TerminationCondition) []TerminationCondition {

	// If the thread has been interrupted, we should abort and return whatever
	// result we currently have.
	// TODO: ? what is this????
	//if (Thread.currentThread().isInterrupted()) {
	//return Collections.emptyList();
	//}
	//// Otherwise check the termination conditions for the evolution.
	satisfied := make([]TerminationCondition, 0)
	for _, cond := range conds {
		if cond.ShouldTerminate(data) {
			satisfied = append(satisfied, cond)
		}
	}
	if len(satisfied) == 0 {
		return nil

	}
	return satisfied
}

// SortEvaluatedPopulation sorts an evaluated population in descending order of
// fitness (descending order of fitness score for natural scores, ascending
// order of scores for non-natural scores).
func SortEvaluatedPopulation(evpop Population, natural bool) {
	// Sort candidates in descending order according to fitness.
	if natural {
		// Descending values for natural fitness.
		sort.Sort(sort.Reverse(evpop))
	} else {
		// Ascending values for non-natural fitness.
		sort.Sort(evpop)
	}
}

// ComputePopulationData computes statistics about the current generation of
// evolved individuals, including the fittest candidate.
//
// evpop is the population of candidate solutions with their associated fitness
// scores.
// natural should be true if higher fitness scores mean fitter individuals,
// false otherwise.
// nelites is the number of candidates preserved via elitism.
// genidx is the zero-based index of the current generation/epoch.
// start is the time at which the evolution began.
func ComputePopulationData(
	pop Population,
	natural bool,
	nelites int,
	genidx int,
	start time.Time) *PopulationData {

	stats := NewDataSet(WithInitialCapacity(len(pop)))
	for _, cand := range pop {
		stats.AddValue(cand.Fitness)
	}

	return &PopulationData{
		BestCand:    pop[0].Candidate,
		BestFitness: pop[0].Fitness,
		Mean:        stats.ArithmeticMean(),
		StdDev:      stats.StandardDeviation(),
		Natural:     natural,
		Size:        stats.Len(),
		NumElites:   nelites,
		GenNumber:   genidx,
		Elapsed:     time.Since(start),
	}
}
