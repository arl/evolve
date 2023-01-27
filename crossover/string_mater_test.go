package crossover

import (
	"math/rand"
	"testing"

	"github.com/arl/evolve"
	"github.com/arl/evolve/generator"
)

func TestStringMater(t *testing.T) {
	rng := rand.New(rand.NewSource(99))

	xover := evolve.Crossover[string]{
		Probability: generator.Const(1.0),
		Mater: &StringMater{
			Points: generator.Const(2),
		},
	}

	items := []string{"abcde", "fghij", "klmno", "pqrst"}

	pop := evolve.NewPopulationOf(items, nil)

	for i := 0; i < 20; i++ {
		genes := make(map[rune]struct{}, 20) // used as a set of runes
		xover.Apply(pop, rng)

		if pop.Len() != 4 {
			t.Errorf("pop.Len() = %v, want 4", pop.Len())
		}

		for _, ind := range pop.Candidates {
			if len(ind) != 5 {
				t.Errorf("len(ind) = %v, want 5", pop.Len())
			}

			for _, value := range ind {
				genes[value] = struct{}{}
			}
		}

		// All of the genes should still be present, just mixed up up between
		// individuals.
		if len(genes) != 20 {
			t.Errorf("got %d different genes, want 20", len(genes))
		}
	}
}

func TestStringMaterZeroPoints(t *testing.T) {
	rng := rand.New(rand.NewSource(99))

	m := StringMater{Points: generator.Const(0)}

	p1 := "abcde"
	p2 := "fghij"

	off1, off2 := m.Mate(p1, p2, rng)
	if off1 != p1 {
		t.Errorf("got offspring1 = %v, want %v", off1, p1)
	}
	if off2 != p2 {
		t.Errorf("got offspring2 = %v, want %v", off2, p2)
	}
}

func BenchmarkStringMater(b *testing.B) {
	rng := rand.New(rand.NewSource(99))

	xover := evolve.Crossover[string]{
		Probability: generator.Const(1.0),
		Mater: &StringMater{
			Points: generator.Const(1),
		},
	}

	pop := evolve.NewPopulationOf([]string{"abcde", "fghij", "klmno", "pqrst"}, nil)

	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		xover.Apply(pop, rng)
	}
}
