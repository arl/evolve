package generator

import (
	"constraints"
	"math/rand"

	"github.com/arl/evolve/pkg/bitstring"
)

// Binomial generates of binomially-distributed, unsigned integers.
type Binomial[U constraints.Unsigned] struct {
	rng *rand.Rand

	n Unsigned[U]
	p Float

	// Cache the fixed-point representation of p to avoid having to
	// recalculate it for each value generated.  Only calculate it
	// if and when p changes.
	pBits *bitstring.Bitstring
	lastp float64
}

// NewBinomial creates a Binomial generator.
//
// numTrials generates the number of trials, it's the maximum possible value
// returned by the generator, the values it generates must be strictly positive.
// prob generates the probability of success in any one trial, the values it
// generates must be in the [0 1] range.
func NewBinomial[U constraints.Unsigned](numTrials Unsigned[U], prob Float, rng *rand.Rand) *Binomial[U] {
	return &Binomial[U]{n: numTrials, p: prob, rng: rng}
}

// Next returns the next generated binomially-distributed value.
func (g *Binomial[U]) Next() U {
	// Regenerate the fixed point representation of p if it has changed.
	newP := g.p.Next()
	if g.pBits == nil || newP != g.lastp {
		g.lastp = newP
		g.pBits = floatToFixedBits(newP)
	}

	var totalSuccesses U
	trials := g.n.Next()
	pidx := g.pBits.Len() - 1

	for trials > 0 && pidx >= 0 {
		successes := g.evenProbability(trials)
		trials -= successes
		if g.pBits.Bit(uint(pidx)) {
			totalSuccesses += successes
		}
		pidx--
	}

	return totalSuccesses
}

// generates a binomial with even probability (p=0.5). We simply generate n
// random bits and count the 1's.
func (g *Binomial[U]) evenProbability(n U) U {
	bs := bitstring.Random(uint(n), g.rng)
	return U(bs.OnesCount())
}

// floatToFixedBits converts a floating point value [0 1) into a fixed point
// bit string (with MSB having a value of 0.5).
func floatToFixedBits(v float64) *bitstring.Bitstring {
	if v < 0 || v >= 1 {
		panic("value must be between 0 and 1")
	}

	s := make([]byte, 64)
	bitval := 0.5
	d := v
	i := 0
	for d > 0 {
		if d >= bitval {
			d -= bitval
			s[i] = '1'
		} else {
			s[i] = '0'
		}
		bitval /= 2
		i++
	}

	bs, _ := bitstring.MakeFromString(string(s[:i]))
	return bs
}
