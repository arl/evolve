package crossover

import (
	"math/rand"

	"github.com/arl/bitstring"
	"github.com/arl/evolve/generator"
)

// BitstringMater is a crossover helper that mates pairs of parent bitstrings
// and produces pairs of offsprings.
type BitstringMater struct {
	// Points generator decided the number of cut points to apply, for each mating.
	Points generator.Generator[int]
}

// Mate performs crossover on a pair of parent strings and generate a pair of offsprings.
// Mate is undefined if p1 and p2 do not have the same length.
func (m *BitstringMater) Mate(p1, p2 *bitstring.Bitstring, rng *rand.Rand) (*bitstring.Bitstring, *bitstring.Bitstring) {
	// Decide the number of cut points.
	npts := int(m.Points.Next())

	off1 := bitstring.Clone(p1)
	off2 := bitstring.Clone(p2)

	// Apply as many crossovers as required.
	for i := 0; i < npts; i++ {
		// Cross-over index is always greater than zero and less than the
		// length of the parent so that we always pick a point that will
		// result in a meaningful crossover.
		bitstring.SwapRange(off1, off2, 0, 1+rng.Intn(p1.Len()-1))
	}

	return off1, off2
}
