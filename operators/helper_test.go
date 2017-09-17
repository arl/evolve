package operators

import (
	"math/rand"

	"github.com/aurelien-rainone/evolve/framework"
)

// integerAdjuster is trivial test operator that mutates all integers by adding
// a fixed offset.
type integerAdjuster int

func (op integerAdjuster) Apply(selectedCandidates []framework.Candidate, rng *rand.Rand) []framework.Candidate {
	result := make([]framework.Candidate, len(selectedCandidates))
	for i, c := range selectedCandidates {
		result[i] = c.(int) + int(op)
	}
	return result
}
