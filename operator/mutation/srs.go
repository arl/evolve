package mutation

import (
	"math/rand"

	"github.com/arl/evolve/generator"
)

// SRS (Swap and Reverse Sections) is a mutation operator creating a permutation
// of the original candidate. To do so, the candidate is divided into 3 sections
// by randomly choosing 2 cut points. The middle section is copied as is to the
// mutant, while the right and left are swapped. Before copying, the elements of
// the left section are reversed.
//
// Example:
//
//	[1, 2, 3, 4, 5, 6, 7].
//
// 2 cut points are chosen randomly, for example xp1 = 2 and xp2 = 4:
//
//	[   1    2    3    4    5    6    7]
//	           |          |
//	     left  |  middle  |    right
//	          xp1        xp2
//
// Then if the random probability allows for the candidate to be mutated, the
// resulting candidate will be:
//
//	[5, 6, 7, 3, 4, 2, 1],
//
// SRS is undefined for candidates with less than 3 elements.
type SRS[T any] struct {
	//  Probability is the probability for each candidate to be mutated.
	Probability generator.Float
}

func (op *SRS[T]) Apply(sel [][]T, rng *rand.Rand) [][]T {
	mutpop := make([][]T, len(sel))
	for i := range sel {
		// Copy current candidate.
		cand := sel[i]
		cpy := make([]T, len(cand))
		copy(cpy, cand)

		// Find out the probability of mutation for this candidate.
		prob := op.Probability.Next()
		var xp1, xp2 int
		if rng.Float64() < prob {
			//  Find 2 cut points between 1 and len(can)-1
			xp1 = 1 + rng.Intn(len(cand)-1)
			for {
				xp2 = 1 + rng.Intn(len(cand)-1)
				if xp2 != xp1 {
					if xp2 < xp1 {
						xp1, xp2 = xp2, xp1
					}
					break
				}
			}
			srs(cand, cpy, xp1, xp2)
		}
		mutpop[i] = cpy
	}

	return mutpop
}

// srs performs the 'swap and reverse sections' mutation.
func srs[T any](org, mut []T, xp1, xp2 int) {
	copy(mut, org[xp2:])
	copy(mut[len(org)-xp2:], org[xp1:xp2])
	copy(mut[len(org)-xp1:], org[:xp1])
	reverse(mut[len(org)-xp1:])
}

func reverse[T any](s []T) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}
