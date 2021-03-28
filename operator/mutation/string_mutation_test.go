package mutation

import (
	"math/rand"
	"testing"

	"github.com/arl/evolve/generator"

	"github.com/stretchr/testify/require"
)

func TestStringMutation(t *testing.T) {
	rng := rand.New(rand.NewSource(99))
	alphabet := "abcd"

	sm := &String{
		Alphabet:    alphabet,
		Probability: generator.ConstFloat64(0.5),
	}

	mut := New(sm)

	individual1 := "abcd"
	individual2 := "abab"
	individual3 := "cccc"

	population := []interface{}{individual1, individual2, individual3}

	// Perform several iterations.
	for i := 0; i < 20; i++ {
		population = mut.Apply(population, rng)
		require.Lenf(t, population, 3, "Population size changed after mutation: %v", len(population))

		// Check that each individual is still valid
		for _, individual := range population {
			sind := individual.(string)
			require.Lenf(t, sind, 4, "Individual size changed after mutation: %d", len(sind))
			for _, c := range []byte(sind) {
				require.Containsf(t, []byte(alphabet), c, "Mutation introduced invalid character: %v", c)
			}
		}
	}
}
