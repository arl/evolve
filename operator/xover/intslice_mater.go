package xover

import (
	"math/rand"
)

// IntSliceMater mates a pair of int slices to produce a new pair of int slices
type IntSliceMater struct{}

// Mate performs crossover on a pair of parents to generate a pair of
// offspring.
func (m IntSliceMater) Mate(parent1, parent2 interface{}, npts int64,
	rng *rand.Rand) []interface{} {

	p1, p2 := parent1.([]int), parent2.([]int)

	if len(p1) != len(p2) {
		panic("Cannot perform crossover with different length parents.")
	}
	off1 := make([]int, len(p1))
	copy(off1, p1)
	off2 := make([]int, len(p1))
	copy(off2, p2)

	// Apply as many crossovers as required.
	for i := int64(0); i < npts; i++ {
		// Crossover index is always greater than zero and less than the length
		// of the parent so that we always pick a point that will result in a
		// meaningful crossover.
		xidx := (1 + rng.Intn(len(p1)-1))
		for j := 0; j < xidx; j++ {
			// swap elements j of both offsprings
			off1[j], off2[j] = off2[j], off1[j]
		}
	}
	return []interface{}{off1, off2}
}
