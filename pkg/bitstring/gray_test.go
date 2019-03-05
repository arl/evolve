package bitstring

import (
	"math"
	"testing"
)

func TestBitstringGrayn(t *testing.T) {
	tests := []struct {
		input string
		nbits uint
		want  word
	}{
		{input: "00000000",
			nbits: 1, want: 0},
		{input: "00000111",
			nbits: 3, want: 5},
		{input: "00000111",
			nbits: 4, want: 5},
		{input: "00000111",
			nbits: 5, want: 5},
		{input: "10000000",
			nbits: 8, want: math.MaxUint8},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			bs, _ := MakeFromString(tt.input)
			got := bs.Grayn(tt.nbits, 0)
			if tt.want != got {
				t.Errorf("Bitstring(%s).Grayn(%d, 0) got %s, want %s", tt.input, tt.nbits,
					sprintubits(word(got), tt.nbits), sprintubits(word(tt.want), tt.nbits))
			}
		})
	}
}

func TestBitstringGray16(t *testing.T) {
	tests := []struct {
		input string
		want  uint16
	}{
		{input: "0000000000000000",
			want: 0},
		{input: "0000000000000111",
			want: 5},
		{input: "1000000000000000",
			want: math.MaxUint16},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			bs, _ := MakeFromString(tt.input)
			got := bs.Gray16(0)
			if tt.want != got {
				t.Errorf("Bitstring(%s).Gray16(0) got %s, want %s", tt.input,
					sprintubits(word(got), 16), sprintubits(word(tt.want), 16))
			}
		})
	}
}

func TestBitstringGray32(t *testing.T) {
	tests := []struct {
		input string
		want  uint32
	}{
		{input: "00000000000000000000000000000000",
			want: 0},
		{input: "00000000000000000000000000000111",
			want: 5},
		{input: "10000000000000000000000000000000",
			want: math.MaxUint32},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			bs, _ := MakeFromString(tt.input)
			got := bs.Gray32(0)
			if tt.want != got {
				t.Errorf("Bitstring(%s).Gray32(0) got %s, want %s", tt.input,
					sprintubits(word(got), 32), sprintubits(word(tt.want), 32))
			}
		})
	}
}

func TestBitstringGray64(t *testing.T) {
	tests := []struct {
		input string
		want  uint64
	}{
		{input: "0000000000000000000000000000000000000000000000000000000000000000",
			want: 0},
		{input: "0000000000000000000000000000000000000000000000000000000000000111",
			want: 5},
		{input: "1000000000000000000000000000000000000000000000000000000000000000",
			want: math.MaxUint64},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			bs, _ := MakeFromString(tt.input)
			got := bs.Gray64(0)
			if tt.want != got {
				t.Errorf("Bitstring(%s).Gray64(0) got %s, want %s", tt.input,
					sprintubits(word(got), 64), sprintubits(word(tt.want), 64))
			}
		})
	}
}

func TestBitstringGray8(t *testing.T) {
	tests := []struct {
		input string
		want  uint8
	}{
		{input: "00000000",
			want: 0},
		{input: "00000111",
			want: 5},
		{input: "10000000",
			want: math.MaxUint8},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			bs, _ := MakeFromString(tt.input)
			got := bs.Gray8(0)
			if tt.want != got {
				t.Errorf("Bitstring(%s).Gray8(0) got %s, want %s", tt.input,
					sprintubits(word(got), 8), sprintubits(word(tt.want), 8))
			}
		})
	}
}
