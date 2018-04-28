package operators

import (
	"math/rand"

	"github.com/aurelien-rainone/evolve/framework"
	"github.com/aurelien-rainone/evolve/pkg/bitstring"
)

// BitStringMater mates a pair of bit strings to produce a new pair of bit
// strings
type BitstringMater struct{}

// Mate performs crossover on a pair of parents to generate a pair of
// offspring.
//
// parent1 and parent2 are the two individuals that provides the source
// material for generating offspring.
func (BitstringMater) Mate(
	p1, p2 framework.Candidate, nxpts int64,
	rng *rand.Rand) []framework.Candidate {

	p1_, p2_ := p1.(*bitstring.Bitstring), p2.(*bitstring.Bitstring)

	if p1_.Len() != p2_.Len() {
		panic("Cannot mate parents of different lengths")
	}
	off1 := p1_.Copy()
	off2 := p2_.Copy()

	// Apply as many crossovers as required.
	for i := int64(0); i < nxpts; i++ {
		// Cross-over index is always greater than zero and less than the
		// length of the parent so that we always pick a point that will
		// result in a meaningful crossover.
		crossoverIndex := (1 + rng.Intn(p1_.Len()-1))
		off1.SwapRange(off2, 0, crossoverIndex)
	}
	return []framework.Candidate{off1, off2}
}
