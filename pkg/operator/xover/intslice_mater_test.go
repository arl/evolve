package xover

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIntSliceMater(t *testing.T) {
	rng := rand.New(rand.NewSource(99))

	xover := New(IntSliceMater{})
	pop := make([]interface{}, 4)
	pop[0] = []int{1, 2, 3, 4, 5}
	pop[1] = []int{6, 7, 8, 9, 10}
	pop[2] = []int{11, 12, 13, 14, 15}
	pop[3] = []int{16, 17, 18, 19, 20}

	for i := 0; i < 20; i++ {
		values := make(map[int]struct{}, 20)
		pop = xover.Apply(pop, rng)
		if len(pop) != 4 {
			t.Error("population size changed, want 4, got", len(pop))
		}

		for _, ind := range pop {
			s := ind.([]int)
			if len(s) != 5 {
				t.Error("wrong candidate length, want 5, got", len(s))
			}
			for _, value := range s {
				values[value] = struct{}{}
			}
		}
		// All of the individual elements should still be present, just jumbled up
		// between individuals.
		if len(values) != 20 {
			t.Error("wrong number of candidates, want 20, got", len(values))
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

	xover := New(IntSliceMater{})
	pop := make([]interface{}, 2)
	pop[0] = []int{1, 2, 3, 4, 5}
	pop[1] = []int{2}

	assert.Panics(t, func() { xover.Apply(pop, rng) })
}
