package evolve

import "sync"

// EvaluatePopulation evaluates all individuals and returns an evaluated
// population, sorted, either in descending order of fitness for natural scores,
// or ascending for non-natural scores. If concurrency is greater than 1, then
// the fitness is evaluated concurrently, using a number of goroutines equal to
// 'concurrency'.
func EvaluatePopulation[T any](pop []T, e Evaluator[T], concurrency int) Population[T] {
	if concurrency < 2 {
		// Synchronous evaluation
		evpop := make(Population[T], len(pop))
		for i, candidate := range pop {
			evpop[i] = &Individual[T]{
				Candidate: candidate,
				Fitness:   e.Fitness(candidate, pop),
			}
		}
		return evpop
	}

	evpop := make(Population[T], len(pop))
	wg := sync.WaitGroup{}
	wg.Add(len(pop))
	sem := make(chan struct{}, concurrency)
	for i := range pop {
		i := i
		sem <- struct{}{}
		go func() {
			evpop[i] = &Individual[T]{
				Candidate: pop[i],
				Fitness:   e.Fitness(pop[i], pop),
			}
			wg.Done()
			<-sem
		}()
	}
	wg.Wait()

	return evpop
}
