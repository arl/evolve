package operator

import (
	"math/rand"
	"reflect"
	"sort"
	"testing"

	"github.com/arl/evolve/generator"
	"github.com/arl/evolve/operator/xover"
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
		xover := NewCrossover[[]byte](xover.SliceMater[byte]{})
		xover.Points = generator.Const(0)
		xover.Probability = generator.Const(1.0)

		got := xover.Apply(pop, rng)
		sameStringPop(t, pop, got)
	})

	t.Run("zero_crossover_probability_is_noop", func(t *testing.T) {
		pop := []string{string("abcde"), string("fghij"), string("klmno"), string("pqrst"), string("uvwxy")}
		xover := NewCrossover[string](xover.StringMater{})
		xover.Points = generator.Const(1)
		xover.Probability = generator.Const(0.0)

		got := xover.Apply(pop, rng)
		sameStringPop(t, pop, got)
	})
}

var sink any

func BenchmarkCrossoverApply(b *testing.B) {
	b.ReportAllocs()
	rng := rand.New(rand.NewSource(99))

	pop := [][]byte{[]byte("abcde"), []byte("fghij"), []byte("klmno"), []byte("pqrst"), []byte("uvwxy")}

	b.ResetTimer()
	var res [][]byte
	for n := 0; n < b.N; n++ {
		xover := NewCrossover[[]byte](xover.SliceMater[byte]{})
		xover.Points = generator.Const(1)
		xover.Probability = generator.Const(1.0)
		res = xover.Apply(pop, rng)
	}

	sink = res
}
