package factory

import (
	"math/rand"
	"strings"
	"testing"

	"github.com/aurelien-rainone/evolve/framework"
	"github.com/stretchr/testify/assert"
)

const (
	stringLength   = 8
	populationSize = 10
	alphabet       = "abcdefg"
)

// Make sure each candidate is valid (is the right length and contains only
// valid characters).
// @param population The population to be validated.
func validatePopulation(t *testing.T, population []framework.Candidate, alphabet string) {
	for _, candidate := range population {
		assert.IsType(t, string(""), candidate)
		s := candidate.(string)

		assert.Len(t, []rune(s), stringLength)

		// check generated string is only made of alphabet characters
		for _, r := range s {
			assert.True(t, strings.ContainsRune(alphabet, rune(r)),
				"%#U is not contained in '%s'\n", r, alphabet)
		}
	}
}

func TestStringFactory(t *testing.T) {
	rng := rand.New(rand.NewSource(99))

	t.Run("string population with ascii-only aplhabet", func(*testing.T) {
		sf, err := NewStringFactory(alphabet, stringLength)
		assert.NoError(t, err)
		pop := sf.GenerateInitialPopulation(populationSize, rng)
		validatePopulation(t, pop, alphabet)
	})

	t.Run("string population with non ascii-only aplhabet", func(*testing.T) {
		_, err := NewStringFactory("日本語", stringLength)
		assert.Error(t, err)
	})

	t.Run("StringFactory with empty aplhabet", func(*testing.T) {
		_, err := NewStringFactory("", stringLength)
		assert.Error(t, err)
	})

	t.Run("StringFactory with zero string length", func(*testing.T) {
		_, err := NewStringFactory(alphabet, 0)
		assert.Error(t, err)
	})
}

func TestStringGenerator(t *testing.T) {
	rng := rand.New(rand.NewSource(99))

	stringLength := 8
	gen := stringGenerator{alphabet, stringLength}

	// create some random candidates
	for i := 0; i < 10; i++ {
		iface := gen.GenerateRandomCandidate(rng)

		// check candidate type is string
		assert.IsType(t, "", iface)
		s := iface.(string)

		// check string length
		assert.Len(t, s, int(stringLength))

		// check generated string is only made of alphabet characters
		for _, r := range s {
			assert.True(t, strings.ContainsRune(alphabet, rune(r)),
				"%#U is not contained in '%s'\n", r, alphabet)
		}
	}
}
