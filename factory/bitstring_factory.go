package factory

import (
	"math/rand"

	"github.com/aurelien-rainone/evolve/bitstring"
	"github.com/aurelien-rainone/evolve/framework"
)

// BitStringFactory is a general purpose candidate factory for generating bit
// strings for genetic algorithms.
type BitStringFactory struct {
	AbstractCandidateFactory
}

// NewBitStringFactory creates a factory that generates bit strings of the
// specified length
func NewBitStringFactory(length int) *BitStringFactory {
	return &BitStringFactory{
		AbstractCandidateFactory{
			&bitStringGenerator{
				length: length,
			},
		},
	}
}

type bitStringGenerator struct {
	length int
}

// GenerateRandomCandidate generates a random bit string, with a uniform
// distribution of ones and zeroes.
func (g *bitStringGenerator) GenerateRandomCandidate(rng *rand.Rand) framework.Candidate {
	bs, _ := bitstring.NewRandom(g.length, rng)
	return bs
}
