package operators

import (
	"math/rand"

	"github.com/aurelien-rainone/evolve/base"
)

// NewStringCrossover creates a crossover operator for string candidates.
func NewStringCrossover(options ...abstractCrossoverOption) (*AbstractCrossover, error) {
	return NewAbstractCrossover(stringMater{}, options...)
}

type stringMater struct{}

func (m stringMater) Mate(parent1, parent2 base.Candidate,
	numberOfCrossoverPoints int64,
	rng *rand.Rand) []base.Candidate {

	p1, p2 := parent1.(string), parent2.(string)

	if len(p1) != len(p2) {
		panic("Cannot perform cross-over with different length parents.")
	}
	offspring1, offspring2 := []byte(p1), []byte(p2)

	// Apply as many cross-overs as required.
	for i := int64(0); i < numberOfCrossoverPoints; i++ {
		// Cross-over index is always greater than zero and less than
		// the length of the parent so that we always pick a point that
		// will result in a meaningful cross-over.
		crossoverIndex := (1 + rng.Intn(len(p1)-1))
		for j := 0; j < crossoverIndex; j++ {
			// swap elements j of both offsprings
			offspring1[j], offspring2[j] = offspring2[j], offspring1[j]
		}
	}
	return append([]base.Candidate{}, string(offspring1), string(offspring2))
}
