package evolve

// A Population holds a group of candidates alongside their fitness.
type Population[T any] struct {
	Candidates       []T
	Fitness          []float64
	FitnessEvaluated []bool

	Evaluator Evaluator[T]
}

// NewPopulation creates a new population, pre-allocating the slices of
// candidates and fitness to n items each.
func NewPopulation[T any](n int, evaluator Evaluator[T]) *Population[T] {
	return &Population[T]{
		Candidates:       make([]T, n),
		Fitness:          make([]float64, n),
		FitnessEvaluated: make([]bool, n),
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
