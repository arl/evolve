package xover

import (
	"math/rand"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/arl/evolve/generator"
)

func sameStringPop(t *testing.T, a, b []interface{}) {
	t.Helper()

	s1 := make([]string, 0)
	s2 := make([]string, 0)
	for _, a := range a {
		s1 = append(s1, a.(string))
	}
	for _, b := range b {
		s2 = append(s2, b.(string))
	}
	sort.Strings(s1)
	sort.Strings(s2)
	assert.EqualValues(t, s1, s2)
}

func TestCrossover_Apply(t *testing.T) {
	rng := rand.New(rand.NewSource(99))
	pop := []interface{}{"abcde", "fghij", "klmno", "pqrst", "uvwxy"}

	t.Run("zero_crossover_points_is_noop", func(t *testing.T) {
		xover := New(StringMater{})
		xover.Points = generator.ConstInt(0)
		xover.Probability = generator.ConstFloat64(1)

		got := xover.Apply(pop, rng)
		sameStringPop(t, pop, got)
	})

	t.Run("zero_crossover_probability_is_noop", func(t *testing.T) {
		xover := New(StringMater{})
		xover.Points = generator.ConstInt(1)
		xover.Probability = generator.ConstFloat64(0.0)

		got := xover.Apply(pop, rng)
		sameStringPop(t, pop, got)
	})
}
