package mutation

import (
	"math/rand"

	"github.com/arl/evolve/generator"
)

type ListOrder struct {
	MutationCount  generator.Int
	MutationAmount generator.Int
}

func (op *ListOrder) Apply(sel []interface{}, rng *rand.Rand) []interface{} {
	mutpop := make([]interface{}, len(sel))
	for i := range sel {
		// copy current candidate
		cand := sel[i].([]int)
		newCand := make([]int, len(cand))
		copy(newCand, cand)

		// determine the mutation count
		nmut := int(op.MutationCount.Next())
		// if op.varnmut {
		// 	nmut = op.nmutmin + rng.Intn(op.nmutmax-op.nmutmin)
		// } else {
		// 	nmut = op.nmut
		// }

		for imut := 0; imut < nmut; imut++ {
			istart := rng.Intn(len(newCand))

			// determine the amount of mutations for current item
			mutAmount := int(op.MutationAmount.Next())
			// if op.varMutAmount {
			// 	mutAmount = op.mutAmountMin + rng.Intn(op.mutAmountMax-op.mutAmountMin)
			// } else {
			// 	mutAmount = op.mutAmount
			// }
			iend := (istart + mutAmount) % len(newCand)
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
