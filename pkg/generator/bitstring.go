package generator

import (
	"math/rand"

	"github.com/aurelien-rainone/evolve/pkg/bitstring"
)

// Bitstring generates random Bitstring of a specified length.
type Bitstring int

// Generate generates a random bit string. The a distribution of ones
// and zeroes depends on rng.
func (i Bitstring) Generate(rng *rand.Rand) interface{} {
	bs, _ := bitstring.Random(int(i), rng)
	return bs
}
