package mutation

import (
	"math/rand"

	"github.com/arl/evolve/generator"
)

// CIM (Center Inverse Mutation) is a mutation operator creating a permutation
// of the original candidate. To do so, the candidate is divided into 2 sections
// by randomly choosing a cut point. All elements in each section are copied and
// then inversely placed in the same section of the mutant Example:
//
//	[1, 2, 3, 4, 5, 6].
//
// The cut point is chosen randomly, for example xp = 4:
//
//	[   1    2    3    4    5    6 ]
//	                     |
//	                    xp
//
// Then if the random probability allows for the candidate to be mutated, the
// resulting candidate will be:
//
//	[4, 3, 2, 1,  6, 5]
//
// CIM is undefined for candidates with less than 2 elements.
type CIM[T any] struct {
	//  Probability is the probability for each candidate to be mutated.
	Probability generator.Float
}

func (op *CIM[T]) Apply(sel [][]T, rng *rand.Rand) [][]T {
	mutpop := make([][]T, len(sel))
	for i := range sel {
		// Copy current candidate.
		cand := sel[i]
		cpy := make([]T, len(cand))
		copy(cpy, cand)

		// Find out the probability of mutation for this candidate.
		prob := op.Probability.Next()
		var xp int
		if rng.Float64() < prob {
			//  Find a cut point between 1 and len(cand)-1
			xp = 1 + rng.Intn(len(cand)-1)
			cim(cpy, xp)
		}
		mutpop[i] = cpy
	}

	return mutpop
}

// cim performs the 'center inverse mutation'.
func cim[T any](mut []T, xp int) {
	reverse(mut[:xp])
	reverse(mut[xp:])
}
