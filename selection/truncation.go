package selection

import (
	"constraints"
	"math"
	"math/rand"

	"github.com/arl/evolve"
	"github.com/arl/evolve/generator"
)

// Truncation is a selection strategy that selects the fittest individuals of a
// population. The number of individuals that are selected depends on the
// SelectionRatio.
type Truncation[T any] struct {
	// SelectionRatio informs the ratio of the fittest individuals that are
	// selected. For example, a value of 0.5 would selects from the fittest
	// half. A value of 1 would select from the whole population.
	SelectionRatio generator.Float
}

// Select selects n individuals from a ratio of the fittest individuals. If the
// set of eligible candidates is smaller than n, the fittest are selected more
// than once.
func (ts *Truncation[T]) Select(pop *evolve.Population[T], natural bool, n int, rng *rand.Rand) []T {
	sel := make([]T, 0, n)

	eligible := int(math.Round(ts.SelectionRatio.Next() * float64(pop.Len())))
	if eligible > n {
		eligible = n
	}

	for {
		count := min(eligible, n-len(sel))
		for i := 0; i < count; i++ {
			sel = append(sel, pop.Candidates[i])
		}
		if len(sel) >= n {
			break
		}
	}
	return sel
}

func (ts *Truncation[T]) String() string {
	return "Truncation Selection"
}

func min[T constraints.Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}
