package factory

import (
	"bytes"
	"errors"
	"fmt"
	"math/rand"
	"unicode/utf8"

	"github.com/aurelien-rainone/evolve/framework"
)

// StringFactory is a general-purpose candidate factory for EAs that use a
// fixed-length string encoding.
// Generates random strings of a fixed length from a given alphabet.
type StringFactory struct {
	AbstractCandidateFactory
}

// NewStringFactory creates a StringFactory.
//
// - alphabet is the set of characters that can legally occur within a string
// generated by this factory.
// - stringLength The fixed length of all strings generated by this factory.
func NewStringFactory(alphabet string, stringLength int) (*StringFactory, error) {
	// safety checks
	if len(alphabet) == 0 {
		return nil, errors.New("can't create StringFactory with an empty alphabet")
	}
	if stringLength == 0 {
		return nil, errors.New("can't create StringFactory with a string length equals to 0")
	}

	// ascii only alphabet
	if utf8.RuneCountInString(alphabet) != len(alphabet) {
		return nil, fmt.Errorf("non ascii alphabet is not supported %v", alphabet)
	}
	// unicode alphabet
	sf := &StringFactory{
		AbstractCandidateFactory{
			&stringGenerator{
				alphabet:     alphabet,
				stringLength: stringLength,
			},
		},
	}
	return sf, nil
}

type stringGenerator struct {
	alphabet     string
	stringLength int
}

// GenerateRandomCandidate generates a random string of a pre-configured length.
//
// Each character is randomly selected from the pre-configured alphabet. The
// same character may appear multiple times and some characters may not appear
// at all.
func (g *stringGenerator) GenerateRandomCandidate(rng *rand.Rand) framework.Candidate {
	var buffer bytes.Buffer
	for i := 0; i < g.stringLength; i++ {
		idx := rand.Int31n(int32(len(g.alphabet)))
		buffer.WriteByte(g.alphabet[idx])
	}
	return buffer.String()
}
