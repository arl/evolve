package factory

import (
	"math/rand"
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/stretchr/testify/assert"
)

func TestStringFactory(t *testing.T) {
	const (
		stringLength   = 8
		populationSize = 10
	)
	rng := rand.New(rand.NewSource(99))

	t.Run("StringFactory with a ascii only aplhabet", func(*testing.T) {
		sf, err := NewStringFactory("abcdefg", stringLength)
		assert.Nil(t, err)
		pop := sf.GenerateInitialPopulation(populationSize, rng)
		assert.Len(t, pop, populationSize)
	})

	t.Run("StringFactory with non-ascii aplhabet", func(*testing.T) {
		sf, err := NewStringFactory("日本語", stringLength)
		assert.Nil(t, err)
		pop := sf.GenerateInitialPopulation(populationSize, rng)
		assert.Len(t, pop, populationSize)
	})

	t.Run("StringFactory with empty aplhabet", func(*testing.T) {
		_, err := NewStringFactory("", stringLength)
		assert.NotNil(t, err)
	})

	t.Run("StringFactory with empty aplhabet", func(*testing.T) {
		_, err := NewStringFactory("abcdefg", 0)
		assert.NotNil(t, err)
	})
}

func TestAsciiStringGenerator(t *testing.T) {
	rng := rand.New(rand.NewSource(99))

	alphabet := "abcdefg"
	stringLength := 8
	gen := asciiStringGenerator{alphabet, stringLength}

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

func TestUnicodeStringGenerator(t *testing.T) {
	rng := rand.New(rand.NewSource(99))

	alphabet := []rune("日本語")
	stringLength := 8
	gen := unicodeStringGenerator{alphabet, stringLength}

	// create some random candidates
	for i := 0; i < 10; i++ {
		iface := gen.GenerateRandomCandidate(rng)

		// check candidate type is string
		assert.IsType(t, "", iface)
		s := iface.(string)

		// check string length
		assert.Equal(t, utf8.RuneCountInString(s), stringLength)

		// check generated string is only made of alphabet characters
		for _, runeValue := range s {
			assert.Truef(t, strings.ContainsRune(string(alphabet), runeValue),
				"%#U is not contained in '%s'\n", runeValue, string(alphabet))
		}
	}
}
