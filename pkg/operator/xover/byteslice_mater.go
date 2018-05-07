package xover

import (
	"math/rand"

	"github.com/aurelien-rainone/evolve/pkg/api"
)

// ByteSliceMater mates two []byte and produces a new pair of []byte
type ByteSliceMater struct{}

// Mate performs crossover on a pair of parents to generate a pair of
// offspring.
func (m ByteSliceMater) Mate(
	parent1, parent2 api.Candidate, nxpts int64,
	rng *rand.Rand) []api.Candidate {

	p1, p2 := parent1.([]byte), parent2.([]byte)

	if len(p1) != len(p2) {
		panic("Cannot perform crossover with different length parents.")
	}
	off1 := make([]byte, len(p1))
	copy(off1, p1)
	off2 := make([]byte, len(p1))
	copy(off2, p2)

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
	return append([]api.Candidate{}, off1, off2)
}
