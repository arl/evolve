package mutation

import (
	"math/rand"

	"github.com/arl/evolve"
	"github.com/arl/evolve/generator"
)

// Permutation is a mutation operator acting on population of []T. It permutes
// elements position.
type Permutation[T any] struct {
	// Count indicates how many slice items are mutated.
	Count generator.Generator[int]

	// Amount indicates how many times each mutated item gets swapped with
	// another randomly selected one.
	Amount generator.Generator[int]

	// Probability specifies the probability for a candidate to be mutated.
	Probability generator.Float
}

func (op *Permutation[T]) Apply(pop *evolve.Population[[]T], rng *rand.Rand) {
	for i := 0; i < pop.Len(); i++ {
		cand := pop.Candidates[i]

		if rng.Float64() < op.Probability.Next() {
			count := op.Count.Next()

			for i := 0; i < count; i++ {
				istart := rng.Intn(len(cand))
				iend := (istart + op.Amount.Next()) % len(cand)
				if iend < 0 {
					iend += len(cand)
				}

				cand[istart], cand[iend] = cand[iend], cand[istart]
			}
		}
	}
}
