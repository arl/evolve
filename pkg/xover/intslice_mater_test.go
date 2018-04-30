package operators

import (
	"math/rand"
	"testing"

	"github.com/aurelien-rainone/evolve/framework"
	"github.com/aurelien-rainone/evolve/number"
	"github.com/stretchr/testify/assert"
)

func TestIntSliceMater(t *testing.T) {
	rng := rand.New(rand.NewSource(99))

	xover, err := NewCrossover(IntSliceMater{})
	if assert.NoError(t, err) {
		pop := make([]framework.Candidate, 4)
		pop[0] = []int{1, 2, 3, 4, 5}
		pop[1] = []int{6, 7, 8, 9, 10}
		pop[2] = []int{11, 12, 13, 14, 15}
		pop[3] = []int{16, 17, 18, 19, 20}

		for i := 0; i < 20; i++ {
			values := make(map[int]struct{}, 20) // used as a set of runes
			pop = xover.Apply(pop, rng)
			assert.Len(t, pop, 4, "Population size changed after crossover.")
			for _, ind := range pop {
				s := ind.([]int)
				assert.Lenf(t, s, 5, "Invalid candidate length: %v", len(s))
				for _, value := range s {
					values[value] = struct{}{}
				}
			}
			// All of the individual elements should still be present, just jumbled up
			// between individuals.
			assert.Len(t, values, 20, "Information lost during crossover.")
		}
	}
}

// The IntArrayCrossover operator is only defined to work on populations
// containing arrays of equal lengths. Any attempt to apply the operation to
// populations that contain different length arrays should panic. Not panicking
// should be considered a bug since it could lead to hard to trace bugs
// elsewhere.
func TestIntArrayCrossoverWithDifferentLengthParents(t *testing.T) {
	rng := rand.New(rand.NewSource(99))

	xover, err := NewCrossover(IntSliceMater{},
		ConstantCrossoverPoints(1),
		ConstantProbability(number.ProbabilityOne),
	)
	if assert.NoError(t, err) {
		pop := make([]framework.Candidate, 2)
		pop[0] = []int{1, 2, 3, 4, 5}
		pop[1] = []int{2}

		// This should panic since the parents are different lengths.
		// TODO: why panicking and not returning an error?
		assert.Panics(t, func() {
			xover.Apply(pop, rng)
		})
	}
}
