package xover

import (
	"math/rand"

	"github.com/aurelien-rainone/evolve/framework"
)

// StringMater mates a pair of strings to produce a new pair of bit strings
type StringMater struct{}

func (m StringMater) Mate(
	p1, p2 framework.Candidate, nxpts int64,
	rng *rand.Rand) []framework.Candidate {

	p1_, p2_ := p1.(string), p2.(string)
	if len(p1_) != len(p2_) {
		panic("Cannot perform crossover with different length parents.")
	}

	off1, off2 := []byte(p1_), []byte(p2_)

	// Apply as many crossovers as required.
	for i := int64(0); i < nxpts; i++ {
		// Cross-over index is always greater than zero and less than
		// the length of the parent so that we always pick a point that
		// will result in a meaningful crossover.
		xidx := (1 + rng.Intn(len(p1_)-1))
		for j := 0; j < xidx; j++ {
			// swap elements j of both offsprings
			off1[j], off2[j] = off2[j], off1[j]
		}
	}
	return append([]framework.Candidate{}, string(off1), string(off2))
}
