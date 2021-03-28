package generator

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"
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
	require.NoError(t, err)

	s := gen.Generate(rand.New(rand.NewSource(99)))
	if s, ok := s.(string); !ok {
		t.Errorf("GenerateCandidate should generate string candidates, got %T", s)
	}
}

var sink interface{}

func BenchmarkGenerateString(b *testing.B) {
	rng := rand.New(rand.NewSource(99))

	runs := []int{10, 100, 1000}
	for _, slen := range runs {
		b.Run(fmt.Sprintf("%d", slen), func(b *testing.B) {
			b.ReportAllocs()
			gen, _ := NewString("A", slen)
			for i := 0; i < b.N; i++ {
				sink = gen.Generate(rng)
			}
		})
	}
}
