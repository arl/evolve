package bitstring

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

//
// conversion to unsigned integers
//

func TestBitstringUintn(t *testing.T) {
	tests := []struct {
		input    string
		nbits, i uint
		want     uint32
	}{
		// LSB and MSB 8 are both on the same word
		{input: "10",
			nbits: 1, i: 0, want: 0},
		{input: "111",
			nbits: 1, i: 0, want: 1},
		{input: "101",
			nbits: 1, i: 1, want: 0},
		{input: "010",
			nbits: 1, i: 1, want: 1},
		{input: "100",
			nbits: 2, i: 0, want: 0},
		{input: "1101",
			nbits: 2, i: 1, want: 2},
		{input: "10100000000000000000000000000000",
			nbits: 3, i: 29, want: 5},
		{input: "10000000000000000000000000000000",
			nbits: 1, i: 31, want: 1},

		// // LSB and MSB 8 are on 2 separate words
		{input: "1111111111111111111111111111111111111111111111111111111111111111",
			nbits: 3, i: 31, want: 7},
		{input: "1111111111111111111111111111111111111111111111111111111111111111",
			nbits: 3, i: 30, want: 7},
		{input: "0000000000000000000000000000001010000000000000000000000000000000",
			nbits: 3, i: 31, want: 5},
		{input: "0000000000000000000000000000000101000000000000000000000000000000",
			nbits: 3, i: 30, want: 5},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			bs, _ := MakeFromString(tt.input)
			got := bs.Uintn(tt.nbits, tt.i)
			require.Equalf(t, tt.want, got,
				"Uintn(nbits=%v, %v) got %s, want %s", tt.nbits, tt.i,
				sprintubits(got, tt.nbits), sprintubits(tt.want, tt.nbits))
		})
	}
}

func TestBitstringUint32(t *testing.T) {
	tests := []struct {
		input string
		i     uint
		want  uint32
	}{
		// LSB and MSB 8 are both on the same word
		{input: "00000000000000000000000000000001",
			i: 0, want: 1},
		{input: "00000000000000000000000000000010",
			i: 0, want: 2},
		{input: "01000000000000000000000000000010",
			i: 0, want: 1<<30 + 2},
		{input: "1111111111111111111111111111111101000000000000000000000000000010",
			i: 0, want: 1<<30 + 2},
		{input: "0000000000000000000000000000000111111111111111111111111111111111",
			i: 32, want: 1},
		{input: "1000000000000000000000000000000111111111111111111111111111111111",
			i: 32, want: 1<<31 + 1},

		// LSB and MSB 8 are on 2 separate words
		{input: "100000000000000000000000000000000",
			i: 1, want: 1 << 31},
		{input: "1111111111111111111101000000000000000000000000000010111111111111",
			i: 12, want: 1<<30 + 2},
		{input: "0000111111111111111111111111111111100000000000000000000000000000",
			i: 28, want: math.MaxUint32 - 1},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			bs, _ := MakeFromString(tt.input)
			got := bs.Uint32(tt.i)
			require.Equalf(t, tt.want, got,
				"Bitstring(%s).Uint32(%v) got %s, want %s", tt.i,
				sprintubits(uint32(got), 32), sprintubits(uint32(tt.want), 32))
		})
	}
}

func TestBitstringUint16(t *testing.T) {
	tests := []struct {
		input string
		i     uint
		uwant uint16
	}{
		// LSB and MSB 8 are both on the same word
		{input: "00000000000000000000000000000001",
			i: 0, uwant: 1},
		{input: "00000000000000000000000000000010",
			i: 0, uwant: 2},
		{input: "00000000000000000100000000000010",
			i: 0, uwant: 1<<14 + 2},
		{input: "11111111111111110100000000000010",
			i: 0, uwant: 1<<14 + 2},
		{input: "0000000000000000000000000000000111111111111111111111111111111111",
			i: 32, uwant: 1},
		{input: "0000000000000000100000000000000111111111111111111111111111111111",
			i: 32, uwant: 1<<15 + 1},
		{input: "10000000000000000",
			i: 1, uwant: 1 << 15},

		// LSB and MSB 8 are on 2 separate words
		{input: "111111111111111111111110100000000000010111111111111111111111111",
			i: 24, uwant: 1<<14 + 2},
		{input: "000000000000000000000001111111111111110000000000000000000000000",
			i: 24, uwant: math.MaxUint16 - 1},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			bs, _ := MakeFromString(tt.input)
			got := bs.Uint16(tt.i)
			require.Equalf(t, tt.uwant, got,
				"Bitstring(%s).Uint16(%v) got %s, want %s", tt.input, tt.i,
				sprintubits(uint32(got), 16), sprintubits(uint32(tt.uwant), 16))
		})
	}
}

func TestBitstringUint8(t *testing.T) {
	tests := []struct {
		input string
		i     uint
		want  uint8
	}{
		// LSB and MSB 8 are both on the same word
		{input: "00000000000000000000000000000001",
			i: 0, want: 1},
		{input: "00000000000000000000000000000010",
			i: 0, want: 2},
		{input: "000000000000000001000010",
			i: 0, want: 1<<6 + 2},
		{input: "111111111111111101000010",
			i: 0, want: 1<<6 + 2},
		{input: "0000000000000000000000000000000111111111111111111111111111111111",
			i: 32, want: 1},
		{input: "00000000000000001000000111111111111111111111111111111111",
			i: 32, want: 1<<7 + 1},
		{input: "100000000",
			i: 1, want: 1 << 7},

		// LSB and MSB 8 are on separate words
		{input: "11111111111111111111111010000101111111111111111111111111111111",
			i: 31, want: 1<<6 + 2},
		{input: "00000000000000000000000111111100000000000000000000000000000000",
			i: 31, want: math.MaxUint8 - 1},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			bs, _ := MakeFromString(tt.input)
			got := bs.Uint8(tt.i)
			require.Equalf(t, tt.want, got,
				"Uint8(%v) got %s, want %s", tt.i,
				sprintubits(uint32(got), 8), sprintubits(uint32(tt.want), 8))
		})
	}
}

//
// conversion to signed integers
//
func TestBitstringInt32(t *testing.T) {
	tests := []struct {
		input string
		i     uint
		want  int32
	}{
		// LSB and MSB 8 are both on the same word
		{input: "11111111111111111111111111111111",
			i: 0, want: -1},
		{input: "01111111111111111111111111111111",
			i: 0, want: math.MaxInt32},
		{input: "10000000000000000000000000000000",
			i: 0, want: math.MinInt32},
		// LSB and MSB 8 are on 2 separate words
		{input: "111111111111111111111111111111110000000000000000000000000000000",
			i: 31, want: -1},
		{input: "011111111111111111111111111111110000000000000000000000000000000",
			i: 31, want: math.MaxInt32},
		{input: "100000000000000000000000000000001111111111111111111111111111111",
			i: 31, want: math.MinInt32},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			bs, _ := MakeFromString(tt.input)
			got := bs.Int32(tt.i)
			require.Equalf(t, tt.want, got,
				"Bitstring(%s).Int32(%v) got %s, want %s", tt.input, tt.i,
				sprintsbits(int32(got), 32), sprintsbits(int32(tt.want), 32))
		})
	}
}

func TestBitstringInt16(t *testing.T) {
	tests := []struct {
		input string
		i     uint
		want  int16
	}{
		// LSB and MSB 8 are both on the same word
		{input: "1111111111111111",
			i: 0, want: -1},
		{input: "0111111111111111",
			i: 0, want: math.MaxInt16},
		{input: "1000000000000000",
			i: 0, want: math.MinInt16},
		// LSB and MSB 8 are on 2 separate words
		{input: "11111111111111110000000000000000000000000000000",
			i: 31, want: -1},
		{input: "01111111111111110000000000000000000000000000000",
			i: 31, want: math.MaxInt16},
		{input: "10000000000000001111111111111111111111111111111",
			i: 31, want: math.MinInt16},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			bs, _ := MakeFromString(tt.input)
			got := bs.Int16(tt.i)
			require.Equalf(t, tt.want, got,
				"Bitstring(%s).Int16(%v) got %s, want %s", tt.input, tt.i,
				sprintsbits(int32(got), 16), sprintsbits(int32(tt.want), 16))
		})
	}
}

// continuer ici a reorganiser le code

func TestBitstringInt8(t *testing.T) {
	tests := []struct {
		input string
		i     uint
		want  int8
	}{
		// LSB and MSB 8 are both on the same word
		{input: "11111111",
			i: 0, want: -1},
		{input: "01111111",
			i: 0, want: math.MaxInt8},
		{input: "10000000",
			i: 0, want: math.MinInt8},
		// LSB and MSB 8 are on 2 separate words
		{input: "111111110000000000000000000000000000000",
			i: 31, want: -1},
		{input: "011111110000000000000000000000000000000",
			i: 31, want: math.MaxInt8},
		{input: "100000001111111111111111111111111111111",
			i: 31, want: math.MinInt8},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			bs, _ := MakeFromString(tt.input)
			got := bs.Int8(tt.i)
			require.Equalf(t, tt.want, got,
				"Bitstring(%s).Int8(%v) got %s, want %s", tt.input, tt.i,
				sprintsbits(int32(got), 8), sprintsbits(int32(tt.want), 8))
		})
	}
}
