package evolve

import "sync"

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
