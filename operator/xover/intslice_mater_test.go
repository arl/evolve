package xover

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/arl/evolve/generator"
)

func TestIntSliceMater(t *testing.T) {
	rng := rand.New(rand.NewSource(99))

	xover := New(IntSliceMater{})
	xover.Points = generator.ConstInt(1)
	xover.Probability = generator.ConstFloat64(1)

	pop := make([]interface{}, 4)
	pop[0] = []int{1, 2, 3, 4, 5}
	pop[1] = []int{6, 7, 8, 9, 10}
	pop[2] = []int{11, 12, 13, 14, 15}
	pop[3] = []int{16, 17, 18, 19, 20}

	for i := 0; i < 20; i++ {
		values := make(map[int]struct{}, 20)
		pop = xover.Apply(pop, rng)

		assert.Lenf(t, pop, 4, "population size changed")

		for _, ind := range pop {
			s := ind.([]int)
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

// IntSliceMater must operate on []int of equal length. It should panic if
// different length slices are used.
func TestIntSliceMaterDifferentLength(t *testing.T) {
	rng := rand.New(rand.NewSource(99))

	xover := New(IntSliceMater{})
	xover.Points = generator.ConstInt(1)
	xover.Probability = generator.ConstFloat64(1)

	pop := make([]interface{}, 2)
	pop[0] = []int{1, 2, 3, 4, 5}
	pop[1] = []int{2}

	assert.Panics(t, func() { xover.Apply(pop, rng) })
}
