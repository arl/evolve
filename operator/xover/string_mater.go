package xover

import (
	"math/rand"
)

// StringMater mates a pair of strings to produce a new pair of bit strings
type StringMater struct{}

// Mate performs crossover on a pair of parents to generate a pair of
// offspring.
func (m StringMater) Mate(
	parent1, parent2 interface{}, nxpts int64,
	rng *rand.Rand) []interface{} {

	p1, p2 := parent1.(string), parent2.(string)
	if len(p1) != len(p2) {
		panic("StringMater only mates string having the same length")
	}

	off1, off2 := []byte(p1), []byte(p2)

	// Apply as many crossovers as required.
	for i := int64(0); i < nxpts; i++ {
		// Cross-over index is always greater than zero and less than
		// the length of the parent so that we always pick a point that
		// will result in a meaningful crossover.
		xidx := (1 + rng.Intn(len(p1)-1))
		for j := 0; j < xidx; j++ {
			// swap elements j of both offsprings
			off1[j], off2[j] = off2[j], off1[j]
		}
	}
	return []interface{}{string(off1), string(off2)}
}
