package mutation

import (
	"math/rand"
	"strings"
	"testing"

	"github.com/arl/evolve"
	"github.com/arl/evolve/generator"
	"github.com/arl/evolve/operator"
)

func TestStringMutation(t *testing.T) {
	rng := rand.New(rand.NewSource(99))

	const alphabet = "abcd"
	mut := operator.NewMutation[string](
		&String{
			Alphabet:    alphabet,
			Probability: generator.Const(0.5),
		},
	)

	items := []string{"abcd", "abab", "cccc"}
	pop := evolve.NewPopulationOf(items, nil)

	// Mutate the population multiple times, check the population size doesn't
	// change and that mutants only contains characters of the alphabet.
	for i := 0; i < 20; i++ {
		mut.Apply(pop, rng)
		if pop.Len() != 3 {
			t.Errorf("pop.Len() = %d, want 3", pop.Len())
		}

		// Check that each individual is still valid
		for _, ind := range pop.Candidates {
			if len(ind) != 4 {
				t.Errorf("len(ind) = %d, want 4", len(ind))
			}
			for _, c := range ind {
				if !strings.Contains(alphabet, string(c)) {
					t.Fatalf("invalid char introduced by mutation %v", c)
				}
			}
		}
	}
}
