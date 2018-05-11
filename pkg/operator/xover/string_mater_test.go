package xover

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringMater(t *testing.T) {
	rng := rand.New(rand.NewSource(99))

	xover := New(StringMater{})
	pop := make([]interface{}, 4)
	pop[0] = "abcde"
	pop[1] = "fghij"
	pop[2] = "klmno"
	pop[3] = "pqrst"

	for i := 0; i < 20; i++ {
		values := make(map[rune]struct{}, 20) // used as a set of runes
		pop = xover.Apply(pop, rng)
		if len(pop) != 4 {
			t.Error("population size changed, want 4, got", len(pop))
		}

		for _, individual := range pop {
			s := individual.(string)
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

// StringMater is only defined to work on populations
// containing strings of equal lengths. Any attempt to apply the operation to
// populations that contain different length strings should panic.
func TestStringMaterWithDifferentLengthParents(t *testing.T) {
	rng := rand.New(rand.NewSource(99))

	xover := New(StringMater{})
	pop := []interface{}{"abcde", "fghijklm"}

	assert.Panics(t, func() { xover.Apply(pop, rng) })
}
