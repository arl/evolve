package operators

import (
	"math/rand"

	"github.com/aurelien-rainone/evolve/bitstring"
	"github.com/aurelien-rainone/evolve/framework"
)

// NewBitStringCrossover creates a crossover operator with a configurable number
// of points (fixed or random) for bit strings candidates.
func NewBitStringCrossover(options ...Option) (*AbstractCrossover, error) {
	return NewAbstractCrossover(bitStringMater{}, options...)
}

type bitStringMater struct{}

// Mate performs crossover on a pair of parents to generate a pair of
// offspring.
//
// parent1 and parent2 are the two individuals that provides the source
// material for generating offspring.
func (m bitStringMater) Mate(parent1, parent2 framework.Candidate,
	numberOfCrossoverPoints int64,
	rng *rand.Rand) []framework.Candidate {

	p1, p2 := parent1.(*bitstring.BitString), parent2.(*bitstring.BitString)

	if p1.Len() != p2.Len() {
		panic("Cannot perform crossover with different length parents")
	}
	offspring1 := p1.Copy()
	offspring2 := p2.Copy()

	// Apply as many crossovers as required.
	for i := int64(0); i < numberOfCrossoverPoints; i++ {
		// Cross-over index is always greater than zero and less than the
		// length of the parent so that we always pick a point that will
		// result in a meaningful crossover.
		crossoverIndex := (1 + rng.Intn(p1.Len()-1))
		offspring1.SwapRange(offspring2, 0, crossoverIndex)
	}
	return []framework.Candidate{offspring1, offspring2}
}
