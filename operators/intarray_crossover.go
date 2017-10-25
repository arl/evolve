package operators

import (
	"math/rand"

	"github.com/aurelien-rainone/evolve/framework"
)

// NewIntArrayCrossover creates a crossover operator for array of primitive
// ints.
func NewIntArrayCrossover(options ...Option) (*AbstractCrossover, error) {
	return NewAbstractCrossover(intArrayMater{}, options...)
}

type intArrayMater struct{}

func (m intArrayMater) Mate(parent1, parent2 framework.Candidate,
	numberOfCrossoverPoints int64,
	rng *rand.Rand) []framework.Candidate {

	p1, p2 := parent1.([]int), parent2.([]int)

	if len(p1) != len(p2) {
		panic("Cannot perform crossover with different length parents.")
	}
	offspring1 := make([]int, len(p1))
	copy(offspring1, p1)
	offspring2 := make([]int, len(p1))
	copy(offspring2, p2)

	// Apply as many crossovers as required.
	for i := int64(0); i < numberOfCrossoverPoints; i++ {
		// Cross-over index is always greater than zero and less than
		// the length of the parent so that we always pick a point that
		// will result in a meaningful crossover.
		crossoverIndex := (1 + rng.Intn(len(p1)-1))
		for j := 0; j < crossoverIndex; j++ {
			// swap elements j of both offsprings
			offspring1[j], offspring2[j] = offspring2[j], offspring1[j]
		}
	}
	return append([]framework.Candidate{}, offspring1, offspring2)
}
