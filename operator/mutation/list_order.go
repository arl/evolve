package mutation

import (
	"errors"
	"math"
	"math/rand"
)

// ErrInvalidMutationAmount is the error returned when trying to set an invalid
// mutation amount
var ErrInvalidMutationAmount = errors.New("mutation amount must be in [0,MaxInt32]")

type ListOrder struct {
	// *Mutation

	// handle mutation count: the number of mutations to apply.
	nmut             int
	varnmut          bool
	nmutmin, nmutmax int

	// handle mutation amount: the number of positions by which an element gets
	// moved when applying the mutation
	// TODO: remove the n from these variable names
	mutAmount                  int
	varMutAmount               bool
	mutAmountMin, mutAmountMax int
}

func NewListOrder() *ListOrder {
	op := &ListOrder{}
	op.SetMutations(1)
	op.SetMutationAmount(1)

	return op
}

// SetMutations sets the constant number of mutations to apply to each
// individual in the population.
// If nmut is not in the [0,MaxInt32] range SetMutations will return
// ErrInvalidMutationCount.
func (op *ListOrder) SetMutations(nmut int) error {
	if nmut < 0 || nmut > math.MaxInt32 {
		return ErrInvalidMutationCount
	}
	op.nmut = nmut
	op.varnmut = false
	return nil
}

// SetMutationsRange sets the range of possible number of mutations (i.e the
// number of bits that will be flipped if the bitstring candidate is selected
// for mutation).
//
// The specific number of mutations will be randomly chosen with the pseudo
// random number generator argument of Apply, by linearly converting from
// [0,MaxInt32) to [min,max).
//
// If min and max are not bounded by [0,MaxInt32] SetMutationsRange will return
// ErrInvalidMutationCount.
func (op *ListOrder) SetMutationsRange(min, max int) error {
	if min > max || min < 0 || max > math.MaxInt32 {
		return ErrInvalidMutationCount
	}
	op.nmutmin = min
	op.nmutmax = max
	op.varnmut = true
	return nil
}

// TODO(arl) abstract away the number of mutation and mutation amount.
// a same pattern always repeats. A crossover or mutation depends on a specific number (int or float depending on the case)
// or numbers. These numbers can either be constant or variable (depending on a certain distribution).
// All these behaviours could be easiluy encapsulated in a struct on which:
// - one can define if the number is constant
// - one can define if the number varies and in this case, set the range and the distribution

// SetMutationAmount sets the constant number of positions by which to displace an element when mutation

// * @param mutationAmount A random variable that provides a number
// * of positions by which to displace an element when mutating.

func (op *ListOrder) SetMutationAmount(amount int) error {
	if amount < 0 || amount > math.MaxInt32 {
		return ErrInvalidMutationAmount
	}
	op.mutAmount = amount
	op.varMutAmount = false
	return nil
}

// SetMutationAmountRange ...
func (op *ListOrder) SetMutationAmountRange(min, max int) error {
	if min > max || min < 0 || max > math.MaxInt32 {
		return ErrInvalidMutationAmount
	}
	op.mutAmountMin = min
	op.mutAmountMax = max
	op.varMutAmount = true
	return nil
}

func (op *ListOrder) Apply(sel []interface{}, rng *rand.Rand) []interface{} {
	mutpop := make([]interface{}, len(sel))
	for i := range sel {
		// copy current candidate
		cand := sel[i].([]int)
		newCand := make([]int, len(cand))
		copy(newCand, cand)

		// determine the mutation count
		var nmut int
		if op.varnmut {
			nmut = op.nmutmin + rng.Intn(op.nmutmax-op.nmutmin)
		} else {
			nmut = op.nmut
		}

		for imut := 0; imut < nmut; imut++ {
			istart := rng.Intn(len(newCand))

			// determine the amount of mutations for current item
			var mutAmount int
			if op.varMutAmount {
				mutAmount = op.mutAmountMin + rng.Intn(op.mutAmountMax-op.mutAmountMin)
			} else {
				mutAmount = op.mutAmount
			}
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
