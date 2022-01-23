package xover

import (
	"math/rand"
	"reflect"
	"sort"
	"testing"

	"github.com/arl/evolve/generator"
)

type byteseq interface {
	~string | ~[]byte
}

func sameStringPop[T byteseq](t *testing.T, a, b []T) {
	t.Helper()

	s1 := make([]string, 0)
	s2 := make([]string, 0)
	for _, a := range a {
		s1 = append(s1, string(a))
	}
	for _, b := range b {
		s2 = append(s2, string(b))
	}
	sort.Strings(s1)
	sort.Strings(s2)
	if !reflect.DeepEqual(s1, s2) {
		t.Errorf("different strings\na = %v\nb = %v", s1, s2)
	}
}

func TestCrossoverApply(t *testing.T) {
	rng := rand.New(rand.NewSource(99))

	t.Run("zero_crossover_points_is_noop", func(t *testing.T) {
		pop := [][]byte{[]byte("abcde"), []byte("fghij"), []byte("klmno"), []byte("pqrst"), []byte("uvwxy")}
		xover := New[[]byte](SliceMater[byte]{})
		xover.Points = generator.Const(0)
		xover.Probability = generator.Const(1.0)

		got := xover.Apply(pop, rng)
		sameStringPop(t, pop, got)
	})

	t.Run("zero_crossover_probability_is_noop", func(t *testing.T) {
		pop := []string{string("abcde"), string("fghij"), string("klmno"), string("pqrst"), string("uvwxy")}
		xover := New[string](StringMater{})
		xover.Points = generator.Const(1)
		xover.Probability = generator.Const(0.0)

		got := xover.Apply(pop, rng)
		sameStringPop(t, pop, got)
	})
}
