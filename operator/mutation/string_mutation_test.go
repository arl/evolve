package mutation

import (
	"math/rand"
	"testing"

	"github.com/arl/evolve"
	"github.com/arl/evolve/generator"
	"github.com/arl/evolve/operator"

	"github.com/stretchr/testify/require"
)

func TestStringMutation(t *testing.T) {
	rng := rand.New(rand.NewSource(99))
	alphabet := "abcd"

	sm := &String{
		Alphabet:    alphabet,
		Probability: generator.Const(0.5),
	}

	mut := operator.NewMutation[string](sm)

	individual1 := "abcd"
	individual2 := "abab"
	individual3 := "cccc"

	items := []string{individual1, individual2, individual3}
	pop := evolve.NewPopulationOf(items, nil)

	// Perform several iterations.
	for i := 0; i < 20; i++ {
		mut.Apply(pop, rng)
		require.Lenf(t, pop, 3, "Population size changed after mutation: %v", pop.Len())

		// Check that each individual is still valid
		for _, ind := range pop.Candidates {
			require.Lenf(t, ind, 4, "Individual size changed after mutation: %d", len(ind))
			for _, c := range []byte(ind) {
				require.Containsf(t, []byte(alphabet), c, "Mutation introduced invalid character: %v", c)
			}
		}
	}
}
