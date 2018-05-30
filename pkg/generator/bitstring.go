package generator

import (
	"math/rand"

	"github.com/arl/evolve/pkg/bitstring"
)

// Bitstring generates random Bitstring of a specified length.
type Bitstring int

// Generate generates a random bit string in which the distribution of ones and
// zeroes depends on rng.
func (i Bitstring) Generate(rng *rand.Rand) interface{} {
	bs, err := bitstring.Random(int(i), rng)
	if err != nil {
		panic(err)
	}
	return bs
}
