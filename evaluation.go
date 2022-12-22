package evolve

import (
	"sync"
)

// EvaluatePopulation evaluates all individuals and returns an evaluated
// population, sorted, either in descending order of fitness for natural scores,
// or ascending for non-natural scores. If concurrency is greater than 1, then
// the fitness is evaluated concurrently, using a number of goroutines equal to
// 'concurrency'.
func EvaluatePopulation[T any](pop []T, e Evaluator[T], concurrency int) *Population[T] {
	evpop := &Population[T]{
		Candidates: make([]T, len(pop)),
		Fitness:    make([]float64, len(pop)),
	}

	if concurrency < 2 {
		// Synchronous evaluation
		for i, cand := range pop {
			evpop.Candidates[i] = cand
			evpop.Fitness[i] = e.Fitness(cand, pop)
		}
		return evpop
	}

	wg := sync.WaitGroup{}
	wg.Add(len(pop))
	sem := make(chan struct{}, concurrency)
	for i := range pop {
		i := i
		sem <- struct{}{}
		go func() {
			evpop.Candidates[i] = pop[i]
			evpop.Fitness[i] = e.Fitness(pop[i], pop)
			wg.Done()
			<-sem
		}()
	}
	wg.Wait()

	return evpop
}
