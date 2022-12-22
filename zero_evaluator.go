package evolve

// ZeroEvaluator is a fitness evaluator always giving a score of 0.
type ZeroEvaluator[T any] struct{}

// Fitness returns a score of zero, regardless of the candidate being
// evaluated.
func (ZeroEvaluator[T]) Fitness(T, []T) float64 { return 0 }

// IsNatural always returns true. However, it shouldn't be relevant since
// fitness is always 0.
func (ZeroEvaluator[T]) IsNatural() bool { return true }
