package factory

import (
	"math/rand"
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/stretchr/testify/assert"
)

func TestStringFactory(t *testing.T) {
	rng := rand.New(rand.NewSource(99))
	alphabet := "abcdefg"
	stringLength := 8
	populationSize := 10

	cf := NewStringFactory(alphabet, stringLength)
	pop := cf.GenerateInitialPopulation(populationSize, rng)
	assert.Len(t, pop, populationSize)
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
