package evolve

import "sync"

// EvaluatePopulation evaluates individuals and returns a sorted population.
//
// Every individual is assigned a fitness score with the provided evaluator,
// after that individuals is inserted into a population. The population is then
// sorted, either in descending order of fitness for natural scores, or ascending
// for non natural scores.
//
// Returns the evaluated population (a slice of individuals, each of which associated with its fitness).
func EvaluatePopulation(pop []interface{}, e Evaluator, concurrent bool) Population {
	var evpop Population

	if !concurrent {

		evpop = make(Population, len(pop))
		for i, candidate := range pop {
			evpop[i] = &Individual{
				Candidate: candidate,
				Fitness:   e.Fitness(candidate, pop),
			}
		}

	} else {

		evpop = make(Population, len(pop))

		var w sync.WaitGroup
		w.Add(len(pop))

		for i := range pop {
			go func(i int) {
				ind := &Individual{
					Candidate: pop[i],
					Fitness:   e.Fitness(pop[i], pop),
				}
				evpop[i] = ind
				w.Done()
			}(i)
		}

		w.Wait()

		// TODO: handle goroutine termination
	}

	return evpop
}
