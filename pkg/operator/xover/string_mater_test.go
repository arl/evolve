package xover

import (
	"math/rand"
	"testing"

	"github.com/aurelien-rainone/evolve/framework"
	"github.com/stretchr/testify/assert"
)

func TestStringMater(t *testing.T) {
	rng := rand.New(rand.NewSource(99))

	xover := NewCrossover(StringMater{})
	population := make([]framework.Candidate, 4)
	population[0] = "abcde"
	population[1] = "fghij"
	population[2] = "klmno"
	population[3] = "pqrst"

	for i := 0; i < 20; i++ {
		values := make(map[rune]struct{}, 20) // used as a set of runes
		population = xover.Apply(population, rng)
		assert.Len(t, population, 4, "Population size changed after crossover.")
		for _, individual := range population {
			s := individual.(string)
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

// The StringCrossover operator is only defined to work on populations
// containing strings of equal lengths. Any attempt to apply the operation to
// populations that contain different length strings should panic. Not panicking
// should be considered a bug since it could lead to hard to trace bugs
// elsewhere.
func TestStringMaterWithDifferentLengthParents(t *testing.T) {
	rng := rand.New(rand.NewSource(99))

	xover := NewCrossover(StringMater{})
	population := []framework.Candidate{"abcde", "fghijklm"}

	// This should panic since the parents are different lengths.
	// TODO: why panicking and not returning an error?
	assert.Panics(t, func() {
		xover.Apply(population, rng)
	})
}

//func TestStringCrossoverNoop(t *testing.T) {
//rng := rand.New(rand.NewSource(99))

//t.Run("constant_crossover_points_cant_be_zero", func(t *testing.T) {
//// If created with a specified (constant) number of crossover points,
//// this number must be greater than 0 or the operator is a no-op.
//xover := NewCrossover(StringMater{})
//xover.SetPoints(0)
//assert.Error(t, err)
//assert.Nilf(t, xover, "want string crossover to be nil if invalid, got %v", xover)
//})

//t.Run("zero_crossover_points_is_noop", func(t *testing.T) {
//// If created with a variable number of crossover points,
//// verifies that when this number happens to be 0, the operator is a
//// no-op.
//xover, err := NewCrossover(StringMater{}, VariableCrossoverPoints(zeroGenerator{}))
//if assert.NoError(t, err) {
//population := []framework.Candidate{"abcde", "fghij"}
//crossed := xover.Apply([]framework.Candidate{population[0], population[1]}, rng)
//assert.Equal(t, population, crossed)
//}
//})

//t.Run("zero_crossover_probability_is_noop", func(t *testing.T) {
//// If created wit a variable number of crossover probability,
//// verifies that when this number happens to be 0, the operator is a
//// no-op.
//xover, err := NewCrossover(StringMater{}, ConstantProbability(number.ProbabilityZero))
//if assert.NoError(t, err) {
//population := []framework.Candidate{"abcde", "fghij"}
//crossed := xover.Apply([]framework.Candidate{population[0], population[1]}, rng)
//assert.Equal(t, population, crossed)
//}
/*})*/
//}

type zeroGenerator struct{}

func (g zeroGenerator) NextValue() int64 {
	return 0
}