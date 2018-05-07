package factory

import (
	"math/rand"

	"github.com/aurelien-rainone/evolve/pkg/api"
	"github.com/aurelien-rainone/evolve/pkg/bitstring"
)

// Bitstring is a general purpose candidate factory for generating bit
// strings for genetic algorithms.
type Bitstring struct{ BaseFactory }

// NewBitstring returns a factory that generates bitstrings of the
// specified length
func NewBitstring(length int) *Bitstring {
	return &Bitstring{BaseFactory{bitstringGenerator(length)}}
}

type bitstringGenerator int

// GenerateCandidate generates a random bit string, with a uniform
// distribution of ones and zeroes.
func (i bitstringGenerator) GenerateCandidate(rng *rand.Rand) api.Candidate {
	bs, _ := bitstring.Random(int(i), rng)
	return bs
}
