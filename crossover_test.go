package evolve

import (
	"math/rand"
	"reflect"
	"sort"
	"testing"

	"github.com/arl/evolve/crossover"
	"github.com/arl/evolve/generator"
)

func sameStringPop(t *testing.T, a, b []string) {
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

	org := []string{
		"abcde",
		"fghij",
		"klmno",
		"pqrst",
		"uvwxy",
	}

	t.Run("zero_crossover_points_is_noop", func(t *testing.T) {
		xover := NewCrossover[string](crossover.StringMater{})
		xover.Points = generator.Const(0)
		xover.Probability = generator.Const(1.0)

		items := make([]string, len(org))
		copy(items, org)

		pop := NewPopulationOf(items, nil)

		xover.Apply(pop, rng)
		sameStringPop(t, pop.Candidates, org)
	})

	t.Run("zero_crossover_probability_is_noop", func(t *testing.T) {
		xover := NewCrossover[string](crossover.StringMater{})
		xover.Points = generator.Const(1)
		xover.Probability = generator.Const(0.0)

		items := make([]string, len(org))
		copy(items, org)

		pop := NewPopulationOf(items, nil)

		xover.Apply(pop, rng)
		sameStringPop(t, pop.Candidates, org)
	})
}

var sink any

func BenchmarkCrossoverApply(b *testing.B) {
	b.ReportAllocs()
	rng := rand.New(rand.NewSource(99))

	items := [][]byte{
		[]byte("abcde"),
		[]byte("fghij"),
		[]byte("klmno"),
		[]byte("pqrst"),
		[]byte("uvwxy"),
	}
	pop := NewPopulationOf(items, nil)

	b.ResetTimer()
	var res [][]byte
	for n := 0; n < b.N; n++ {
		xover := NewCrossover[[]byte](crossover.SliceMater[byte]{})
		xover.Points = generator.Const(1)
		xover.Probability = generator.Const(1.0)

		xover.Apply(pop, rng)
	}

	sink = res
}
