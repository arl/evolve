package xover

import (
	"math/rand"
	"testing"

	"github.com/arl/evolve"
	"github.com/arl/evolve/generator"
	"github.com/arl/evolve/operator"
)

func TestStringMater(t *testing.T) {
	rng := rand.New(rand.NewSource(99))

	xover := operator.NewCrossover[string](StringMater{})
	xover.Points = generator.Const(1)
	xover.Probability = generator.Const(1.0)

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

func TestStringMaterWithDifferentLengthParents(t *testing.T) {
	// StringMater is only defined for population of strings of equal lengths
	rng := rand.New(rand.NewSource(99))
	xover := operator.NewCrossover[string](StringMater{})
	pop := evolve.NewPopulationOf([]string{"abcde", "fghijklm"}, nil)

	if !didPanic(func() { xover.Apply(pop, rng) }) {
		t.Fatalf("Should have panicked")
	}
}

// didPanic returns true if the function passed to it panics
func didPanic(f func()) (panicked bool) {
	panicked = true
	defer func() {
		recover()
	}()
	f()
	panicked = false
	return
}

func BenchmarkStringMater(b *testing.B) {
	rng := rand.New(rand.NewSource(99))

	xover := operator.NewCrossover[string](StringMater{})
	xover.Probability = generator.Const(1.0)
	xover.Points = generator.Const(1)

	pop := evolve.NewPopulationOf([]string{"abcde", "fghij", "klmno", "pqrst"}, nil)

	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		xover.Apply(pop, rng)
	}
}
