package xover

import (
	"math/rand"
)

// ByteSliceMater mates two []byte and produces a new pair of []byte
type ByteSliceMater struct{}

// Mate performs crossover on a pair of parents to generate a pair of
// offspring.
func (m ByteSliceMater) Mate(
	parent1, parent2 interface{}, nxpts int64,
	rng *rand.Rand) []interface{} {

	p1, p2 := parent1.([]byte), parent2.([]byte)

	// TODO: we chose to panick here, making those assumptions:
	// - always returning an error and checking its value in the caller would
	// eat up some cpu cycle on the critical path.
	// - in this specific case panics 'should' only happen for clear and
	// irrecoverable logic errors, for example when trying to mate 2 slices of
	// different lengths...
	// But even if that's true, the cost should really be negligible compared to
	// the individual fitness evaluation, that is what takes time in non-trivial
	// GAs.
	// So that let us, once again, with having to chose between error checking
	// and/or panicking. IMO it should be refactored in favor of error checking!
	if len(p1) != len(p2) {
		panic("Cannot perform crossover with different length parents.")
	}
	off1 := append([]byte{}, p1...)
	off2 := append([]byte{}, p2...)

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
	return []interface{}{off1, off2}
}
