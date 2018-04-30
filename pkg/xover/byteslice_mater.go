package operators

import (
	"math/rand"

	"github.com/aurelien-rainone/evolve/framework"
)

// NewByteSliceCrossover creates a crossover operating on slices of bytes.
func NewByteSliceCrossover(options ...Option) (*Crossover, error) {
	return NewCrossover(byteSliceMater{}, options...)
}

type byteSliceMater struct{}

func (m byteSliceMater) Mate(
	p1, p2 framework.Candidate, nxpts int64,
	rng *rand.Rand) []framework.Candidate {

	p1_, p2_ := p1.([]byte), p2.([]byte)

	if len(p1_) != len(p2_) {
		panic("Cannot perform crossover with different length parents.")
	}
	off1 := make([]byte, len(p1_))
	copy(off1, p1_)
	off2 := make([]byte, len(p1_))
	copy(off2, p2_)

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
	return append([]framework.Candidate{}, off1, off2)
}
