package factory

import (
	"math/rand"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	stringLength   = 8
	populationSize = 10
	alphabet       = "abcdefg"
)

// Make sure each candidate is valid (has the right length and contains only
// valid characters).
func validatePopulation(t *testing.T, pop []interface{}, alphabet string) {
	for _, candidate := range pop {
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
		sf, err := NewString(alphabet, stringLength)
		assert.NoError(t, err)
		pop := sf.GenPopulation(populationSize, rng)
		validatePopulation(t, pop, alphabet)
	})

	t.Run("string population with non ascii-only aplhabet", func(*testing.T) {
		_, err := NewString("日本語", stringLength)
		assert.Error(t, err)
	})

	t.Run("StringFactory with empty aplhabet", func(*testing.T) {
		_, err := NewString("", stringLength)
		assert.Error(t, err)
	})

	t.Run("StringFactory with zero string length", func(*testing.T) {
		_, err := NewString(alphabet, 0)
		assert.Error(t, err)
	})
}

func TestStringGenerator(t *testing.T) {
	rng := rand.New(rand.NewSource(99))

	stringLength := 8
	gen := stringGenerator{alphabet, stringLength}

	// create some random candidates
	for i := 0; i < 10; i++ {
		iface := gen.GenerateCandidate(rng)

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
