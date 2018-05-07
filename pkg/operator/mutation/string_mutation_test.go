package mutation

import (
	"math/rand"
	"testing"

	"github.com/aurelien-rainone/evolve/pkg/api"
	"github.com/stretchr/testify/assert"
)

func TestStringMutationTest(t *testing.T) {
	rng := rand.New(rand.NewSource(99))
	alphabet := []byte{'a', 'b', 'c', 'd'}

	mutation := NewStringMutation(string(alphabet))
	err := mutation.SetProb(0.5)
	errcheck(t, err)
	individual1 := "abcd"
	individual2 := "abab"
	individual3 := "cccc"

	population := []api.Candidate{individual1, individual2, individual3}

	// Perform several iterations.
	for i := 0; i < 20; i++ {
		population = mutation.Apply(population, rng)
		assert.Lenf(t, population, 3, "Population size changed after mutation: %v", len(population))
		// Check that each individual is still valid
		for _, individual := range population {
			sind := individual.(string)
			assert.Lenf(t, sind, 4, "Individual size changed after mutation: %d", len(sind))
			for _, c := range sind {
				assert.Containsf(t, alphabet, byte(c), "Mutation introduced invalid character: %c", c)
			}
		}
	}
}
