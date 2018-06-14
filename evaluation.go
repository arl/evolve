package evolve

import (
	"fmt"
	"runtime"

	"github.com/arl/evolve/worker"
)

// EvaluatePopulation takes a population, assigns a fitness score to each member
// and returns the members with their scores attached, sorted in descending
// order of fitness (descending order of fitness score for natural scores,
// ascending order of scores for non-natural scores).
// population is the population to evaluate (each candidate is assigned a
// fitness score).
//
// Returns the evaluated population (a slice of candidates with attached fitness
// scores).
func EvaluatePopulation(pop []interface{}, e Evaluator, concurrent bool) Population {

	// Do fitness evaluations
	evpop := make(Population, len(pop))

	if !concurrent {

		for i, candidate := range pop {
			evpop[i] = &Individual{
				Candidate: candidate,
				Fitness:   e.Fitness(candidate, pop),
			}
		}

	} else {

		// Create a worker pool that will divides the required number of fitness
		// evaluations equally among the available goroutines and coordinate
		// them so that we do not proceed until all of them have finished
		// processing.
		wp := worker.NewPool(runtime.NumCPU())
		workers := make([]worker.Worker, len(pop))

		for i := range pop {
			workers[i] = &fitnessWorker{
				idx:       i,
				pop:       pop,
				evaluator: e,
			}
		}

		results, err := wp.Submit(workers)
		if err != nil {
			panic(fmt.Sprintf("Error while submitting workers to the pool: %v", err))
		}

		for i, result := range results {
			evpop[i] = result.(*Individual)
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

type fitnessWorker struct {
	idx       int           // index of candidate to evaluate
	pop       []interface{} // full population
	evaluator Evaluator
}

func (w *fitnessWorker) Work() (interface{}, error) {
	return &Individual{
		Candidate: w.pop[w.idx],
		Fitness:   w.evaluator.Fitness(w.pop[w.idx], w.pop),
	}, nil
}
