package mutation

import (
	"math/rand"

	"github.com/arl/evolve"
	"github.com/arl/evolve/generator"
)

// SRS (Swap and Reverse Sections) is a mutation operator creating a permutation
// of the original candidate. To do so, the candidate is divided into 3 sections
// by randomly choosing 2 cut points. The middle section remains as-is while the
// left and right are swapped. Before the swap, the elements of the left
// sections are reversed.
//
// Let's take for example the following candidate, with 2 random cut points of 2
// and 4:
//
//	[1, 2, 3, 4, 5, 6, 7]
//	      |     |
//	    cut-1  cut-2
//
// If the randomly selected number for that candidate, with respect to
// Probability, allows it to be mutated, it will become:
//
//	[5, 6, 7, 3, 4, 2, 1]
//	         |     |
//	       cut-1  cut-2
type SRS[T any] struct {
	//  Probability is the probability for each candidate to be mutated.
	Probability generator.Float
}

func (op *SRS[T]) Apply(pop *evolve.Population[[]T], rng *rand.Rand) {
	// lazily created scratch buffer for the reverse+swap.
	var scratch []T

	for i := 0; i < pop.Len(); i++ {
		if rng.Float64() >= op.Probability.Next() {
			continue
		}

		cand := pop.Candidates[i]
		//  Find 2 cut points between 0 and len(can)
		xp1 := rng.Intn(len(cand))     // 0 -> len-1
		xp2 := 1 + rng.Intn(len(cand)) // 1 -> len

		if xp2 < xp1 {
			xp1, xp2 = xp2, xp1
		}

		if len(scratch) < len(cand) {
			scratch = append(scratch, make([]T, len(cand))...)
		}
		copy(scratch, cand)
		srs(scratch, cand, xp1, xp2)
	}
}

// srs performs the 'swap and reverse sections' mutation.
func srs[T any](cpy, mut []T, xp1, xp2 int) {
	copy(mut, cpy[xp2:])
	copy(mut[len(cpy)-xp2:], cpy[xp1:xp2])
	copy(mut[len(cpy)-xp1:], cpy[:xp1])
	reverse(mut[len(cpy)-xp1:])
}

func reverse[T any](s []T) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}
