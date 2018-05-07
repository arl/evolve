package factory

import (
	"math/rand"

	"github.com/aurelien-rainone/evolve/pkg/api"
	"github.com/aurelien-rainone/evolve/pkg/bitstring"
)

// BitstringFactory is a general purpose candidate factory for generating bit
// strings for genetic algorithms.
type BitstringFactory struct {
	AbstractCandidateFactory
}

// NewBitstringFactory creates a factory that generates bit strings of the
// specified length
func NewBitstringFactory(length int) *BitstringFactory {
	return &BitstringFactory{
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
func (g *bitStringGenerator) GenerateRandomCandidate(rng *rand.Rand) api.Candidate {
	bs, _ := bitstring.Random(g.length, rng)
	return bs
}
