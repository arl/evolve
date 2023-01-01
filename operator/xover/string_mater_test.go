package xover

import (
	"math/rand"
	"testing"

	"github.com/arl/evolve/generator"
	"github.com/arl/evolve/operator"
	"github.com/stretchr/testify/assert"
)

func TestStringMater(t *testing.T) {
	rng := rand.New(rand.NewSource(99))

	xover := operator.NewCrossover[string](StringMater{})
	xover.Points = generator.Const(1)
	xover.Probability = generator.Const(1.0)

	pop := make([]string, 4)
	pop[0] = "abcde"
	pop[1] = "fghij"
	pop[2] = "klmno"
	pop[3] = "pqrst"

	for i := 0; i < 20; i++ {
		values := make(map[rune]struct{}, 20) // used as a set of runes
		pop = xover.Apply(pop, rng)

		assert.Lenf(t, pop, 4, "population size changed")

		for _, ind := range pop {
			assert.Lenf(t, ind, 5, "wrong individual length")

			for _, value := range ind {
				values[value] = struct{}{}
			}
		}

		// All of the elements should still be present, just jumbled up between
		// individuals.
		assert.Lenf(t, values, 20, "wrong number of individuals")
	}
}

// StringMater is only defined to work on populations
// containing strings of equal lengths. Any attempt to apply the operation to
// populations that contain different length strings should panic.
func TestStringMaterWithDifferentLengthParents(t *testing.T) {
	rng := rand.New(rand.NewSource(99))

	xover := operator.NewCrossover[string](StringMater{})
	pop := []string{"abcde", "fghijklm"}

	assert.Panics(t, func() { xover.Apply(pop, rng) })
}

func BenchmarkStringMater(b *testing.B) {
	rng := rand.New(rand.NewSource(99))

	xover := operator.NewCrossover[string](StringMater{})
	xover.Probability = generator.Const(1.0)
	xover.Points = generator.Const(1)

	pop := []string{"abcde", "fghij", "klmno", "pqrst"}

	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		pop = xover.Apply(pop, rng)
	}
}
