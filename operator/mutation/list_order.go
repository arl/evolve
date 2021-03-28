package mutation

import (
	"math/rand"

	"github.com/arl/evolve/generator"
)

type ListOrder struct {
	Count          generator.Int
	MutationAmount generator.Int
}

func (op *ListOrder) Apply(sel []interface{}, rng *rand.Rand) []interface{} {
	mutpop := make([]interface{}, len(sel))
	for i := range sel {
		// Copy current candidate.
		cand := sel[i].([]int)
		newCand := make([]int, len(cand))
		copy(newCand, cand)

		// Determine the mutation count.
		count := int(op.Count.Next())

		for i := 0; i < count; i++ {
			istart := rng.Intn(len(newCand))

			// Determine the amount of mutations for current item.
			amount := int(op.MutationAmount.Next())
			iend := (istart + amount) % len(newCand)
			if iend < 0 {
				iend += len(newCand)
			}

			// swap the 2 items
			newCand[istart], newCand[iend] = newCand[iend], newCand[istart]

		}
		mutpop[i] = newCand
	}
	return mutpop
}

/*
func (op *listOrder) Mutate(c interface{}, rng *rand.Rand) interface{} {
	s := c.(string)
	buffer := make([]byte, len(s))
	copy(buffer, []byte(s))

	// get/decide a probability for this run
	prob := op.prob
	if op.varprob {
		prob = op.probmin + (op.probmax-op.probmin)*rng.Float64()
	}

	for i := range buffer {
		if rng.Float64() < prob {
			buffer[i] = op.alphabet[rng.Intn(len(op.alphabet))]
		}
	}
	return string(buffer)
}

*/
