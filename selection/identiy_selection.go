package selection

import (
	"math/rand"

	"github.com/aurelien-rainone/evolve/framework"
)

// Identity is a selection strategy that returns identical candidates
type Identity struct{}

// Select selects the specified number of candidates from the population.
func (sel Identity) Select(
	population framework.EvaluatedPopulation,
	naturalFitnessScores bool,
	selectionSize int,
	rng *rand.Rand) []framework.Candidate {

	selection := make([]framework.Candidate, selectionSize)
	for i := 0; i < selectionSize; i++ {
		selection[i] = population[i].Candidate()
	}
	return selection
}

func (sel Identity) String() string {
	return "Identity Selection"
}
