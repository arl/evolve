package xover

import (
	"math/rand"
	"testing"

	"github.com/arl/evolve/generator"
	"github.com/stretchr/testify/assert"
)

func TestStringMater(t *testing.T) {
	rng := rand.New(rand.NewSource(99))

	xover := New(StringMater{})
	xover.Points = generator.ConstInt(1)
	xover.Probability = generator.ConstFloat64(1)

	pop := make([]interface{}, 4)
	pop[0] = "abcde"
	pop[1] = "fghij"
	pop[2] = "klmno"
	pop[3] = "pqrst"

	for i := 0; i < 20; i++ {
		values := make(map[rune]struct{}, 20) // used as a set of runes
		pop = xover.Apply(pop, rng)

		assert.Lenf(t, pop, 4, "population size changed")

		for _, individual := range pop {
			s := individual.(string)
			assert.Lenf(t, s, 5, "wrong individual length")

			for _, value := range s {
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

	xover := New(StringMater{})
	pop := []interface{}{"abcde", "fghijklm"}

	assert.Panics(t, func() { xover.Apply(pop, rng) })
}

func BenchmarkStringMater(b *testing.B) {
	rng := rand.New(rand.NewSource(99))

	xover := New(StringMater{})
	xover.Probability = generator.ConstFloat64(1.0)
	xover.Points = generator.ConstInt(1)

	pop := make([]interface{}, 4)
	pop[0] = "abcde"
	pop[1] = "fghij"
	pop[2] = "klmno"
	pop[3] = "pqrst"

	b.ReportAllocs()
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		pop = xover.Apply(pop, rng)
	}
}
