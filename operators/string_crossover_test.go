package operators

import (
	"math/rand"
	"testing"

	"github.com/aurelien-rainone/evolve/framework"
	"github.com/aurelien-rainone/evolve/number"
	"github.com/stretchr/testify/assert"
)

func TestStringCrossover(t *testing.T) {
	rng := rand.New(rand.NewSource(99))

	crossover, err := NewStringCrossover()
	if assert.NoError(t, err) {
		population := make([]framework.Candidate, 4)
		population[0] = "abcde"
		population[1] = "fghij"
		population[2] = "klmno"
		population[3] = "pqrst"

		for i := 0; i < 20; i++ {
			values := make(map[rune]struct{}, 20) // used as a set of runes
			population = crossover.Apply(population, rng)
			assert.Len(t, population, 4, "Population size changed after cross-over.")
			for _, individual := range population {
				s := individual.(string)
				assert.Lenf(t, s, 5, "Invalid candidate length: %v", len(s))
				for _, value := range s {
					values[value] = struct{}{}
				}
			}
			// All of the individual elements should still be present, just jumbled up
			// between individuals.
			assert.Len(t, values, 20, "Information lost during cross-over.")
		}
	}
}

// The StringCrossover operator is only defined to work on populations
// containing strings of equal lengths. Any attempt to apply the operation to
// populations that contain different length strings should panic. Not panicking
// should be considered a bug since it could lead to hard to trace bugs
// elsewhere.
func TestStringCrossoverWithDifferentLengthParents(t *testing.T) {
	rng := rand.New(rand.NewSource(99))

	crossover, err := NewStringCrossover(
		WithConstantCrossoverPoints(1),
		WithConstantCrossoverProbability(number.ProbabilityOne),
	)
	if assert.NoError(t, err) {
		population := make([]framework.Candidate, 2)
		population[0] = "abcde"
		population[1] = "fghijklm"

		// This should panic since the parents are different lengths.
		// TODO: why panicking and not returning an error?
		assert.Panics(t, func() {
			crossover.Apply(population, rng)
		})
	}
}

// Number of cross-over points must be greater than zero otherwise the operator
// is a no-op.
func TestStringCrossoverZeroPoints(t *testing.T) {
	op, err := NewStringCrossover(WithConstantCrossoverPoints(0))
	assert.Error(t, err)
	assert.Nilf(t, op, "want string crossover to be nil if invalid, got %v", op)
}
