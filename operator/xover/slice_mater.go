package xover

import "math/rand"

// SliceMater mates a pair of slices to produce a new pair of slices.
type SliceMater[T any] struct{}

// Mate performs crossover on a pair of parents to generate a pair of offspring.
func (m SliceMater[T]) Mate(p1, p2 []T, nxpts int, rng *rand.Rand) (off1, off2 []T) {
	if len(p1) != len(p2) {
		panic("SliceMater only mates slices of the same length")
	}

	off1 = make([]T, len(p1))
	off2 = make([]T, len(p2))
	copy(off1, p1)
	copy(off2, p2)

	// Apply as many crossovers as required.
	for i := 0; i < nxpts; i++ {
		// Cross-over index is always greater than zero and less than
		// the length of the parent so that we always pick a point that
		// will result in a meaningful crossover.
		xidx := (1 + rng.Intn(len(p1)-1))
		for j := 0; j < xidx; j++ {
			// swap elements j of both offsprings
			off1[j], off2[j] = off2[j], off1[j]
		}
	}

	return
}
