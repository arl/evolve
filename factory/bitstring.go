package factory

import (
	"math/rand"

	"github.com/arl/bitstring"
)

// Bitstring creates random bit strings of a specified length.
type Bitstring uint

// New creates a random bit string in which the distribution of ones and zeroes
// depends on rng.
func (i Bitstring) New(rng *rand.Rand) *bitstring.Bitstring {
	return bitstring.Random(int(i), rng)
}
