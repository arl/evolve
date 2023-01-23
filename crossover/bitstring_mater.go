package crossover

import (
	"math/rand"

	"github.com/arl/bitstring"
)

// BitstringMater mates a pair of bit-strings to produce a new pair of
// bit-strings
type BitstringMater struct{}

// Mate performs crossover on a pair of parents to generate a pair of offspring.
//
// p1 and p2 are the two individuals that provides the source material for
// generating offspring.
func (BitstringMater) Mate(p1, p2 *bitstring.Bitstring, nxpts int, rng *rand.Rand) (off1, off2 *bitstring.Bitstring) {
	if p1.Len() != p2.Len() {
		panic("Cannot mate Bitstring of different lengths")
	}
	off1 = bitstring.Clone(p1)
	off2 = bitstring.Clone(p2)

	// Apply as many crossovers as required.
	for i := 0; i < nxpts; i++ {
		// Cross-over index is always greater than zero and less than the
		// length of the parent so that we always pick a point that will
		// result in a meaningful crossover.
		bitstring.SwapRange(off1, off2, 0, 1+rng.Intn(p1.Len()-1))
	}
	return
}
