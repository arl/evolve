package mutation

import (
	"math/rand"

	"github.com/arl/evolve"
	"github.com/arl/evolve/generator"
)

// CIM (Center Inverse Mutation) is a mutation operator creating a permutation
// of the original candidate. To do so, a randomly chosen cut point divides the
// candidate into 2 sections. In each section, the order of elements are
// reversed, before reassembling the section into a new candidate.
//
// Let's take for example the following candidate, with a random cut points of
// 4:
//
//	[1, 2, 3, 4, 5, 6]
//	            |
//	        cut point
//
// If the randomly selected number for that candidate, with respect to
// Probability, allows it to be mutated, it will become:
//
//	[4, 3, 2, 1, 6, 5]
//	            |
//	        cut point
type CIM[T any] struct {
	//  Probability is the probability for each candidate to be mutated.
	Probability generator.Float
}

func (op *CIM[T]) Apply(pop *evolve.Population[[]T], rng *rand.Rand) {
	for i := 0; i < pop.Len(); i++ {
		if rng.Float64() >= op.Probability.Next() {
			continue
		}

		cut := rng.Intn(len(pop.Candidates[i]) - 1)
		cim(pop.Candidates[i], cut)
	}
}

// cim performs the 'center inverse mutation'.
func cim[T any](mut []T, cut int) {
	reverse(mut[:cut])
	reverse(mut[cut:])
}
