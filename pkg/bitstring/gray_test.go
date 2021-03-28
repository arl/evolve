package bitstring

import (
	"math"
	"testing"
)

func TestBitstringGray8(t *testing.T) {
	tests := []struct {
		input string
		want  uint8
	}{
		{input: "00000000", want: 0},
		{input: "00000111", want: 5},
		{input: "10000000", want: math.MaxUint8},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			bs, _ := MakeFromString(tt.input)
			got := bs.Gray8(0)
			if tt.want != got {
				t.Errorf("Bitstring(%s).Gray8(0) got %s, want %s", tt.input,
					sprintubits(uint(got), 8), sprintubits(uint(tt.want), 8))
			}
		})
	}
}

func TestBitstringGray16(t *testing.T) {
	tests := []struct {
		input string
		want  uint16
	}{
		{input: "0000000000000000", want: 0},
		{input: "0000000000000111", want: 5},
		{input: "1000000000000000", want: math.MaxUint16},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			bs, _ := MakeFromString(tt.input)
			got := bs.Gray16(0)
			if tt.want != got {
				t.Errorf("Bitstring(%s).Gray16(0) got %s, want %s", tt.input,
					sprintubits(uint(got), 16), sprintubits(uint(tt.want), 16))
			}
		})
	}
}

func TestBitstringGray32(t *testing.T) {
	tests := []struct {
		input string
		want  uint32
	}{
		{input: "00000000000000000000000000000000", want: 0},
		{input: "00000000000000000000000000000111", want: 5},
		{input: "10000000000000000000000000000000", want: math.MaxUint32},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			bs, _ := MakeFromString(tt.input)
			got := bs.Gray32(0)
			if tt.want != got {
				t.Errorf("Bitstring(%s).Gray32(0) got %s, want %s", tt.input,
					sprintubits(uint(got), 32), sprintubits(uint(tt.want), 32))
			}
		})
	}
}

func TestBitstringGray64(t *testing.T) {
	tests := []struct {
		input string
		want  uint64
	}{
		{input: "0000000000000000000000000000000000000000000000000000000000000000", want: 0},
		{input: "0000000000000000000000000000000000000000000000000000000000000111", want: 5},
		{input: "1000000000000000000000000000000000000000000000000000000000000000", want: math.MaxUint64},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			bs, _ := MakeFromString(tt.input)
			got := bs.Gray64(0)
			if tt.want != got {
				t.Errorf("Bitstring(%s).Gray64(0) got %s, want %s", tt.input,
					sprintubits(uint(got), 64), sprintubits(uint(tt.want), 64))
			}
		})
	}
}

func TestBitstringGrayn(t *testing.T) {
	tests := []struct {
		input string
		nbits uint
		want  uint
	}{
		{input: "00000000", nbits: 1, want: 0},
		{input: "00000111", nbits: 3, want: 5},
		{input: "00000111", nbits: 4, want: 5},
		{input: "00000111", nbits: 5, want: 5},
		{input: "10000000", nbits: 8, want: math.MaxUint8},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			bs, _ := MakeFromString(tt.input)
			got := bs.Grayn(tt.nbits, 0)
			if tt.want != got {
				t.Errorf("Bitstring(%s).Grayn(%d, 0) got %s, want %s", tt.input, tt.nbits,
					sprintubits(uint(got), tt.nbits), sprintubits(uint(tt.want), tt.nbits))
			}
		})
	}
}
