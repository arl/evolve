package generator

import (
	"math/rand"
	"testing"
)

const (
	stringLength   = 8
	populationSize = 10
	alphabet       = "abcdefg"
)

func TestNewString(t *testing.T) {
	tests := []struct {
		name     string
		alphabet string
		length   int
		wantErr  error
	}{
		{
			"valid string generator",
			"abcdefgh12324;?:",
			9,
			nil,
		},
		{
			"empty string generator",
			"abcdefgh12324;?:",
			0,
			nil,
		},
		{
			"not ASCII-only alphabet",
			"abcdefgh12324本;?:",
			0,
			ErrNotASCIIAlphabet,
		},

		{
			"empty alphabet",
			"",
			10,
			ErrEmptyAlphabet,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewString(tt.alphabet, tt.length)
			if tt.wantErr != err {
				t.Errorf("NewString(), wantErr = %v, got %v", tt.wantErr, err)
			}
		})
	}
}

func TestStringGenerator(t *testing.T) {
	gen, err := NewString("ABCdefg", 2)
	if err != nil {
		t.Error(err)
	}
	s := gen.GenerateCandidate(rand.New(rand.NewSource(99)))
	if s, ok := s.(string); !ok {
		t.Errorf("GenerateCandidate should generate string candidates, got %T", s)
	}
}
