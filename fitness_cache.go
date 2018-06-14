package evolve

import "sync"

// FitnessCache provides caching for any Evaluator implementation. The
// results of fitness evaluations are stored in a cache so that if the same
// candidate is evaluated twice, the expense of the fitness calculation can be
// avoided the second time.
//
// Caching of fitness values can be a useful optimisation in situations where the
// fitness evaluation is expensive and there is a possibility that some candidates
// will survive from generation to generation unmodified.  Programs that use elitism
// are one example of candidates surviving unmodified.  Another scenario is when the
// configured evolutionary operator does not always modify every candidate in the
// population for every generation.
//
// Caching of fitness scores is only valid when fitness evaluations are isolated
// and repeatable. An isolated fitness evaluation is one where the result
// depends only upon the candidate being evaluated.  This is not the case when
// candidates are evaluated against the other members of the population.  So
// unless the fitness evaluator ignores the second parameter to the
// Evaluator.Fitness method, caching must not be used.
type FitnessCache struct {

	// Wrapped is the fitness evaluator for which we want to provide caching
	Wrapped Evaluator
	cache   sync.Map
}

// Fitness calculates a fitness score for the given candidate.
//
// This implementation performs a cache look-up every time it is invoked.  If
// the fitness evaluator has already calculated the fitness score for the
// specified candidate that score is returned without delegating to the wrapped
// evaluator.
func (c *FitnessCache) Fitness(cand interface{}, pop []interface{}) float64 {
	var fitness float64
	val, ok := c.cache.Load(cand)
	if ok {
		fitness = val.(float64)
	} else {
		fitness = c.Wrapped.Fitness(cand, pop)
		c.cache.Store(cand, fitness)
	}
	return fitness
}

// IsNatural specifies whether this evaluator generates 'natural' fitness
// scores or not.
func (c *FitnessCache) IsNatural() bool { return c.Wrapped.IsNatural() }
