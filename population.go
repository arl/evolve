package evolve

import (
	"math/rand"
	"sync"
)

// A Population holds a group of candidates alongside their fitness.
type Population[T any] struct {
	Candidates []T
	Fitness    []float64
	Evaluated  []bool

	Evaluator Evaluator[T]
}

// NewPopulation creates a new population of n candidates. Candidates are the
// zero-value of T.
func NewPopulation[T any](n int, evaluator Evaluator[T]) *Population[T] {
	return &Population[T]{
		Candidates: make([]T, n),
		Fitness:    make([]float64, n),
		Evaluated:  make([]bool, n),
		Evaluator:  evaluator,
	}
}

// NewPopulation creates a new population with the following candidates.
func NewPopulationOf[T any](cands []T, evaluator Evaluator[T]) *Population[T] {
	return &Population[T]{
		Candidates: cands,
		Fitness:    make([]float64, len(cands)),
		Evaluated:  make([]bool, len(cands)),
		Evaluator:  evaluator,
	}
}

// NewPopulationWithCapacity creates a new population, pre-allocating internal
// slices to the given length and capacity each.
func NewPopulationWithCapacity[T any](len, cap int, evaluator Evaluator[T]) *Population[T] {
	return &Population[T]{
		Candidates: make([]T, len, cap),
		Fitness:    make([]float64, len, cap),
		Evaluated:  make([]bool, len, cap),
		Evaluator:  evaluator,
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
	p.Evaluated[i], p.Evaluated[j] = p.Evaluated[j], p.Evaluated[i]
}

// Evaluate evaluates all candidates that do not have a fitness yet.
//
// Concurrency controls the number of goroutines to use in the process.
func (p *Population[T]) Evaluate(concurrency int) {
	if concurrency <= 1 {
		// Synchronous evaluation
		for i := 0; i < p.Len(); i++ {
			if !p.Evaluated[i] {
				p.Fitness[i] = p.Evaluator.Fitness(p.Candidates[i])
				p.Evaluated[i] = true
			}
		}
	}

	wg := sync.WaitGroup{}
	wg.Add(p.Len())
	sem := make(chan struct{}, concurrency)
	for i := 0; i < p.Len(); i++ {
		if p.Evaluated[i] {
			continue
		}

		i := i
		sem <- struct{}{}
		go func() {
			p.Fitness[i] = p.Evaluator.Fitness(p.Candidates[i])
			p.Evaluated[i] = true
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
//
// TODO(arl) consider switching to func only (no interface)
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

// GeneratePopulation creates and initializes a new Population of n candidates
// which are randomly generated using the given factory.
func GeneratePopulation[T any](n int, fac Factory[T], e Evaluator[T], rng *rand.Rand) *Population[T] {
	pop := NewPopulationWithCapacity(0, n, e)
	for i := 0; i < n; i++ {
		pop.Candidates = append(pop.Candidates, fac.New(rng))
	}

	// Reslice other slices
	pop.Fitness = pop.Fitness[0:pop.Len()]
	pop.Evaluated = pop.Evaluated[0:pop.Len()]
	return pop
}

// SeedPopulation creates and initializes a new Population of n candidates, some
// of which are seeded (from the provided seeds) and the rest are randomly
// generated using the given factory. If len(seeds) > n then only the first n
// seeds in the slice are used.
//
// Sometimes it is desirable to seed the initial population with some known good
// candidates, providing some hints for the evolution process.
func SeedPopulation[T any](n int, seeds []T, fac Factory[T], e Evaluator[T], rng *rand.Rand) *Population[T] {
	pop := NewPopulationWithCapacity(0, n, e)

	// Seed what we can.
	min := n
	if len(seeds) < n {
		min = len(seeds)
	}
	pop.Candidates = append(pop.Candidates, seeds[:min]...)

	// Generate the rest.
	for pop.Len() < n {
		pop.Candidates = append(pop.Candidates, fac.New(rng))
	}

	// Reslice other slices.
	pop.Fitness = pop.Fitness[0:pop.Len()]
	pop.Evaluated = pop.Evaluated[0:pop.Len()]
	return pop
}
