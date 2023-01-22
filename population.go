package evolve

import (
	"math/rand"
	"sync"
)

// A Population holds a group of candidates alongside their fitness.
type Population[T any] struct {
	Candidates       []T
	Fitness          []float64
	FitnessEvaluated []bool

	Evaluator Evaluator[T]
}

// NewPopulation creates a new population, pre-allocating internal slices to the
// given length each.
func NewPopulation[T any](len int, evaluator Evaluator[T]) *Population[T] {
	return &Population[T]{
		Candidates:       make([]T, len),
		Fitness:          make([]float64, len),
		FitnessEvaluated: make([]bool, len),
		Evaluator:        evaluator,
	}
}

// NewPopulation creates a new population using the provided elements as candidates.
func NewPopulationOf[T any](items []T, evaluator Evaluator[T]) *Population[T] {
	return &Population[T]{
		Candidates:       items,
		Fitness:          make([]float64, len(items)),
		FitnessEvaluated: make([]bool, len(items)),
		Evaluator:        evaluator,
	}
}

// NewPopulationWithCapacity creates a new population, pre-allocating internal
// slices to the given length and capacity each.
func NewPopulationWithCapacity[T any](len, cap int, evaluator Evaluator[T]) *Population[T] {
	return &Population[T]{
		Candidates:       make([]T, len, cap),
		Fitness:          make([]float64, len, cap),
		FitnessEvaluated: make([]bool, len, cap),
		Evaluator:        evaluator,
	}
}

// Len is the number of elements in the collection.
func (p *Population[T]) Len() int { return len(p.Candidates) }

// Less reports whether the element with
// index a should sort before the element with index b.
func (p *Population[T]) Less(i, j int) bool { return p.Fitness[i] < p.Fitness[j] }

// Swap swaps the elements with indexes i and j.
func (p *Population[T]) Swap(i, j int) {
	p.Fitness[i], p.Fitness[j] = p.Fitness[j], p.Fitness[i]
	p.Candidates[i], p.Candidates[j] = p.Candidates[j], p.Candidates[i]
}

// Evaluate evaluates all candidates that do not have a fitness yet.
//
// Concurrency controls the number of goroutines to use in the process.
func (p *Population[T]) Evaluate(concurrency int) {
	if concurrency <= 1 {
		// Synchronous evaluation
		for i := 0; i < p.Len(); i++ {
			if !p.FitnessEvaluated[i] {
				p.Fitness[i] = p.Evaluator.Fitness(p.Candidates[i])
			}
		}
	}

	wg := sync.WaitGroup{}
	wg.Add(p.Len())
	sem := make(chan struct{}, concurrency)
	for i := 0; i < p.Len(); i++ {
		if p.FitnessEvaluated[i] {
			continue
		}

		i := i
		sem <- struct{}{}
		go func() {
			p.Fitness[i] = p.Evaluator.Fitness(p.Candidates[i])
			wg.Done()
			<-sem
		}()
	}
	wg.Wait()
}

// A Factory generates random candidates.
//
// It is used by evolution engine to increase genetic diversity and/or add new
// candidates to a population.
type Factory[T any] interface {
	// New returns a new random candidate, using the provided pseudo-random
	// number generator.
	New(*rand.Rand) T
}

// The FactoryFunc type is an adapter to allow the use of ordinary
// functions as candidate generators. If f is a function with the appropriate
// signature, FactoryFunc(f) is a Factory that calls f.
type FactoryFunc[T any] func(*rand.Rand) T

// New calls f(rng) and returns its return value.
func (f FactoryFunc[T]) New(rng *rand.Rand) T { return f(rng) }

// GeneratePopulation creates a population with n candidates randomly generated
// using the provided factory.
//
// Note: the return Population doesn't have any Evaluator set.
func GeneratePopulation[T any](fac Factory[T], n int, rng *rand.Rand) *Population[T] {
	pop := NewPopulationWithCapacity[T](0, n, nil)
	for i := 0; i < n; i++ {
		pop.Candidates = append(pop.Candidates, fac.New(rng))
	}

	// Reslice other slices
	pop.Fitness = pop.Fitness[0:pop.Len()]
	pop.FitnessEvaluated = pop.FitnessEvaluated[0:pop.Len()]
	return pop
}

// SeedPopulation returns a slice of n candidates, where a part of them are
// seeded while the rest is generated randomly using the provided factory.
// Sometimes it is desirable to seed the initial population with some known good
// candidates, providing some hints for the evolution process.
//
// Note: the return Population doesn't have any Evaluator set.
func SeedPopulation[T any](fac Factory[T], n int, seeds []T, rng *rand.Rand) *Population[T] {
	pop := NewPopulationWithCapacity[T](0, n, nil)
	pop.Candidates = append(pop.Candidates, seeds...)
	for pop.Len() < n {
		pop.Candidates = append(pop.Candidates, fac.New(rng))
	}

	// Reslice other slices
	pop.Fitness = pop.Fitness[0:pop.Len()]
	pop.FitnessEvaluated = pop.FitnessEvaluated[0:pop.Len()]
	return pop
}
