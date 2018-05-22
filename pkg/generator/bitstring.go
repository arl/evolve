package generator

import (
	"math/rand"

	"github.com/aurelien-rainone/evolve/pkg/bitstring"
)

// Bitstring generates random Bitstring of a specified length.
type Bitstring int

// GenerateCandidate generates a random bit string. The a distribution of ones
// and zeroes depends on rng.
func (i Bitstring) GenerateCandidate(rng *rand.Rand) interface{} {
	bs, _ := bitstring.Random(int(i), rng)
	return bs
}
