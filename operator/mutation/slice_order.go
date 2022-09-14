package mutation

import (
	"math/rand"

	"github.com/arl/evolve/generator"
)

// SliceOrder is mutation operator acting on slices of items of type T. When
// applied, a number of Count items are mutated. The Amount indicates how many
// times each item is swapped with another, random one, in the slice.
type SliceOrder[T any] struct {
	Count, Amount generator.Generator[int]
	Probability   generator.Float
}

func (op *SliceOrder[T]) Apply(sel [][]T, rng *rand.Rand) [][]T {
	mutpop := make([][]T, len(sel))
	for i := range sel {
		// Copy current candidate.
		cand := sel[i]
		cpy := make([]T, len(cand))
		copy(cpy, cand)

		// Find out the probability for this candidate.
		prob := op.Probability.Next()
		if rng.Float64() < prob {
			// Perform mutation. Now we determine the mutation count.
			count := op.Count.Next()

			for i := 0; i < count; i++ {
				istart := rng.Intn(len(cpy))

				// Determine the amount of mutations for current item.
				amount := op.Amount.Next()
				iend := (istart + amount) % len(cpy)
				if iend < 0 {
					iend += len(cpy)
				}

				// swap the 2 items
				cpy[istart], cpy[iend] = cpy[iend], cpy[istart]
			}
		}

		mutpop[i] = cpy
	}
	return mutpop
}
