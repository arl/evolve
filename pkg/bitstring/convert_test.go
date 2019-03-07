package bitstring

import (
	"math"
	"testing"
)

//
// conversion to unsigned integers
//

func TestBitstringUintn(t *testing.T) {
	tests := []struct {
		input    string
		nbits, i uint
		want     word
	}{
		// LSB and MSB are both on the same word
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

		// // LSB and MSB are on 2 separate words
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
			if tt.want != got {
				t.Errorf("Bitstring(%s).Uintn(%d, %d) got %s, want %s", tt.input, tt.nbits, tt.i,
					sprintubits(got, tt.nbits), sprintubits(tt.want, tt.nbits))
			}
		})
	}
}

func TestBitstringUint64(t *testing.T) {
	tests := []struct {
		input string
		i     uint
		want  uint64
	}{
		// LSB and MSB are both on the same word
		{input: "0000000000000000000000000000000000000000000000000000000000000001",
			i: 0, want: 1},
		{input: "0000000000000000000000000000000000000000000000000000000000000010",
			i: 0, want: 2},
		{input: "0100000000000000000000000000000000000000000000000000000000000010",
			i: 0, want: 1<<62 + 2},
		{input: "11111111111111111111111111111111111111111111111111111111111111110100000000000000000000000000000000000000000000000000000000000010",
			i: 0, want: 1<<62 + 2},
		{input: "00000000000000000000000000000000000000000000000000000000000000011111111111111111111111111111111111111111111111111111111111111111",
			i: 64, want: 1},
		{input: "10000000000000000000000000000000000000000000000000000000000000011111111111111111111111111111111111111111111111111111111111111111",
			i: 64, want: 1<<63 + 1},

		// LSB and MSB are on 2 separate words
		{input: "10000000000000000000000000000000000000000000000000000000000000000",
			i: 1, want: 1 << 63},
		{input: "1111111111111111111111111110100000000000000000000000000000000000000000000000000000000000010111111111111111111111111111111111111111111111111111111111111",
			i: 60, want: 1<<62 + 2},
		{input: "000011111111111111111111111111111111111111111111111111111111111111100000000000000000000000000000000000000000000000000000000000",
			i: 58, want: math.MaxUint64 - 1},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			bs, _ := MakeFromString(tt.input)
			got := bs.Uint64(tt.i)
			if tt.want != got {
				t.Errorf("Bitstring(%s).Uint64(%d) got %s, want %s", tt.input, tt.i,
					sprintubits(word(got), 64), sprintubits(word(tt.want), 64))
			}
		})
	}
}

func TestBitstringUint32(t *testing.T) {
	tests := []struct {
		input string
		i     uint
		want  uint32
	}{
		// LSB and MSB are both on the same word
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

		// LSB and MSB are on 2 separate words
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
			if tt.want != got {
				t.Errorf("Bitstring(%s).Uint32(%d) got %s, want %s", tt.input, tt.i,
					sprintubits(word(got), 32), sprintubits(word(tt.want), 32))
			}
		})
	}
}

func TestBitstringUint16(t *testing.T) {
	tests := []struct {
		input string
		i     uint
		want  uint16
	}{
		// LSB and MSB are both on the same word
		{input: "00000000000000000000000000000001",
			i: 0, want: 1},
		{input: "00000000000000000000000000000010",
			i: 0, want: 2},
		{input: "00000000000000000100000000000010",
			i: 0, want: 1<<14 + 2},
		{input: "11111111111111110100000000000010",
			i: 0, want: 1<<14 + 2},
		{input: "0000000000000000000000000000000111111111111111111111111111111111",
			i: 32, want: 1},
		{input: "0000000000000000100000000000000111111111111111111111111111111111",
			i: 32, want: 1<<15 + 1},
		{input: "10000000000000000",
			i: 1, want: 1 << 15},

		// LSB and MSB are on 2 separate words
		{input: "111111111111111111111110100000000000010111111111111111111111111",
			i: 24, want: 1<<14 + 2},
		{input: "000000000000000000000001111111111111110000000000000000000000000",
			i: 24, want: math.MaxUint16 - 1},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			bs, _ := MakeFromString(tt.input)
			got := bs.Uint16(tt.i)
			if tt.want != got {
				t.Errorf("Bitstring(%s).Uint16(%d) got %s, want %s", tt.input, tt.i,
					sprintubits(word(got), 16), sprintubits(word(tt.want), 16))
			}
		})
	}
}

func TestBitstringUint8(t *testing.T) {
	tests := []struct {
		input string
		i     uint
		want  uint8
	}{
		// LSB and MSB are both on the same word
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

		// LSB and MSB are on separate words
		{input: "11111111111111111111111010000101111111111111111111111111111111",
			i: 31, want: 1<<6 + 2},
		{input: "00000000000000000000000111111100000000000000000000000000000000",
			i: 31, want: math.MaxUint8 - 1},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			bs, _ := MakeFromString(tt.input)
			got := bs.Uint8(tt.i)
			if tt.want != got {
				t.Errorf("Bitstring(%s).Uint8(%d) got %s, want %s", tt.input, tt.i,
					sprintubits(word(got), 8), sprintubits(word(tt.want), 8))
			}
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
		// LSB and MSB are both on the same word
		{input: "11111111111111111111111111111111",
			i: 0, want: -1},
		{input: "01111111111111111111111111111111",
			i: 0, want: math.MaxInt32},
		{input: "10000000000000000000000000000000",
			i: 0, want: math.MinInt32},
		// LSB and MSB are on 2 separate words
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
			if tt.want != got {
				t.Errorf("Bitstring(%s).Int32(%d) got %s, want %s", tt.input, tt.i,
					sprintsbits(sword(got), 32), sprintsbits(sword(tt.want), 32))
			}
		})
	}
}

func TestBitstringInt16(t *testing.T) {
	tests := []struct {
		input string
		i     uint
		want  int16
	}{
		// LSB and MSB are both on the same word
		{input: "1111111111111111",
			i: 0, want: -1},
		{input: "0111111111111111",
			i: 0, want: math.MaxInt16},
		{input: "1000000000000000",
			i: 0, want: math.MinInt16},
		// LSB and MSB are on 2 separate words
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
			if tt.want != got {
				t.Errorf("Bitstring(%s).Int16(%d) got %s, want %s", tt.input, tt.i,
					sprintsbits(sword(got), 16), sprintsbits(sword(tt.want), 16))
			}
		})
	}
}

func TestBitstringInt8(t *testing.T) {
	tests := []struct {
		input string
		i     uint
		want  int8
	}{
		// LSB and MSB are both on the same word
		{input: "11111111",
			i: 0, want: -1},
		{input: "01111111",
			i: 0, want: math.MaxInt8},
		{input: "10000000",
			i: 0, want: math.MinInt8},
		// LSB and MSB are on 2 separate words
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
			if tt.want != got {
				t.Errorf("Bitstring(%s).Int8(%d) got %s, want %s", tt.input, tt.i,
					sprintsbits(sword(got), 8), sprintsbits(sword(tt.want), 8))
			}
		})
	}
}

//
// conversion from unsigned integers
//

func TestBitstringSetUint8(t *testing.T) {
	tests := []struct {
		bs   string // starting bitstring
		x    uint8  // value to set
		i    uint   // bit index where to set it
		want string
	}{
		// LSB and MSB are both on the same word
		{
			bs: "1111111111111111",
			x:  2, i: 0,
			want: "1111111100000010"},
		{
			bs: "1111111111111111",
			x:  2, i: 8,
			want: "0000001011111111"},
		{
			bs: "11111111111111111111111111111111",
			x:  2, i: 16,
			want: "11111111000000101111111111111111"},
		{
			bs: "11111111111111111111111111111111",
			x:  2, i: 24,
			want: "00000010111111111111111111111111"},
		{
			bs: "11111111111111111111111111111111",
			x:  2, i: 22,
			want: "11000000101111111111111111111111"},
		// LSB and MSB are on separate words
		{
			bs: "111111111111111111111111111111111111111111111111111111111111111111111",
			x:  2, i: 61,
			want: "000000101111111111111111111111111111111111111111111111111111111111111"},
		{
			bs: "11111111111111111111111111111111111111111111111111111111111111111111111",
			x:  15, i: 63,
			want: "00001111111111111111111111111111111111111111111111111111111111111111111"},
		{
			bs: "0000000000000000000000000000000000000000000000000000000000000000000",
			x:  math.MaxUint8, i: 59,
			want: "1111111100000000000000000000000000000000000000000000000000000000000"},
		{
			bs: "0011101010101010101010101010101010101010101010101010101010101010101010101010101010",
			x:  0xaa, i: 63,
			want: "0011101010110101010010101010101010101010101010101010101010101010101010101010101010"},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			bs, _ := MakeFromString(tt.bs)
			bs.SetUint8(tt.i, tt.x)
			want, _ := MakeFromString(tt.want)
			if !want.Equals(bs) {
				t.Errorf("Bitstring(%s).SetUint8(%d, %d) got %s, want %s",
					tt.bs, tt.i, tt.x, bs, want)
			}
		})
	}
}

func TestBitstringSetUint16(t *testing.T) {
	tests := []struct {
		bs   string // starting bitstring
		x    uint16 // value to set
		i    uint   // bit index where to set it
		want string
	}{
		// LSB and MSB are both on the same word
		{
			bs: "1111111111111111",
			x:  2, i: 0,
			want: "0000000000000010"},
		{
			bs: "111111111111111111111111",
			x:  2, i: 8,
			want: "000000000000001011111111"},
		{
			bs: "11111111111111111111111111111111",
			x:  2, i: 16,
			want: "00000000000000101111111111111111"},
		{
			bs: "1111111111111111111111111111111111111111",
			x:  2, i: 24,
			want: "0000000000000010111111111111111111111111"},
		{
			bs: "1111111111111111111111111111111111111111",
			x:  2, i: 22,
			want: "1100000000000000101111111111111111111111"},
		// LSB and MSB are on separate words
		{
			bs: "11111111111111111111111111111111111111111111111111111111111111111111111111111",
			x:  2, i: 61,
			want: "00000000000000101111111111111111111111111111111111111111111111111111111111111"},
		{
			bs: "1111111111111111111111111111111111111111111111111111111111111111111111111111111",
			x:  15, i: 63,
			want: "0000000000001111111111111111111111111111111111111111111111111111111111111111111"},
		{
			bs: "000000000000000000000000000000000000000000000000000000000000000000000000000",
			x:  math.MaxUint16, i: 59,
			want: "111111111111111100000000000000000000000000000000000000000000000000000000000"},
		{
			bs: "001110101010101010101010101010101010101010101010101010101010101010101010101010101010101010",
			x:  0xaaaa, i: 63,
			want: "001110101011010101010101010010101010101010101010101010101010101010101010101010101010101010"},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			bs, _ := MakeFromString(tt.bs)
			bs.SetUint16(tt.i, tt.x)
			want, _ := MakeFromString(tt.want)
			if !want.Equals(bs) {
				t.Errorf("Bitstring(%s).SetUint16(%d, %d) got %s, want %s",
					tt.bs, tt.i, tt.x, bs, want)
			}
		})
	}
}

func TestBitstringSetUint32(t *testing.T) {
	tests := []struct {
		bs   string // starting bitstring
		x    uint32 // value to set
		i    uint   // bit index where to set it
		want string
	}{
		// LSB and MSB are both on the same word
		{
			bs: "11111111111111111111111111111111",
			x:  2, i: 0,
			want: "00000000000000000000000000000010"},
		{
			bs: "1111111111111111111111111111111111111111",
			x:  2, i: 8,
			want: "0000000000000000000000000000001011111111"},
		{
			bs: "111111111111111111111111111111111111111111111111",
			x:  2, i: 16,
			want: "000000000000000000000000000000101111111111111111"},
		{
			bs: "11111111111111111111111111111111111111111111111111111111",
			x:  2, i: 24,
			want: "00000000000000000000000000000010111111111111111111111111"},
		{
			bs: "11111111111111111111111111111111111111111111111111111111",
			x:  2, i: 22,
			want: "11000000000000000000000000000000101111111111111111111111"},
		// LSB and MSB are on separate words
		{
			bs: "111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111",
			x:  2, i: 61,
			want: "000000000000000000000000000000101111111111111111111111111111111111111111111111111111111111111"},
		{
			bs: "11111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111",
			x:  15, i: 63,
			want: "00000000000000000000000000001111111111111111111111111111111111111111111111111111111111111111111"},
		{
			bs: "0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
			x:  math.MaxUint32, i: 59,
			want: "1111111111111111111111111111111100000000000000000000000000000000000000000000000000000000000"},
		{
			bs: "0011101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010",
			x:  0xaaaaaaaa, i: 63,
			want: "0011101010110101010101010101010101010101010010101010101010101010101010101010101010101010101010101010101010"},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			bs, _ := MakeFromString(tt.bs)
			bs.SetUint32(tt.i, tt.x)
			want, _ := MakeFromString(tt.want)
			if !want.Equals(bs) {
				t.Errorf("Bitstring(%s).SetUint32(%d, %d) got %s, want %s",
					tt.bs, tt.i, tt.x, bs, want)
			}
		})
	}
}

func TestBitstringSetUint64(t *testing.T) {
	tests := []struct {
		bs   string // starting bitstring
		x    uint64 // value to set
		i    uint   // bit index where to set it
		want string
	}{
		// LSB and MSB are both on the same word
		{
			bs: "1111111111111111111111111111111111111111111111111111111111111111",
			x:  2, i: 0,
			want: "0000000000000000000000000000000000000000000000000000000000000010"},
		// LSB and MSB are on separate words
		{
			bs: "111111111111111111111111111111111111111111111111111111111111111111111111",
			x:  2, i: 8,
			want: "000000000000000000000000000000000000000000000000000000000000001011111111"},
		{
			bs: "11111111111111111111111111111111111111111111111111111111111111111111111111111111",
			x:  2, i: 16,
			want: "00000000000000000000000000000000000000000000000000000000000000101111111111111111"},
		{
			bs: "1111111111111111111111111111111111111111111111111111111111111111111111111111111111111111",
			x:  2, i: 24,
			want: "0000000000000000000000000000000000000000000000000000000000000010111111111111111111111111"},
		{
			bs: "1111111111111111111111111111111111111111111111111111111111111111111111111111111111111111",
			x:  2, i: 22,
			want: "1100000000000000000000000000000000000000000000000000000000000000101111111111111111111111"},
		{
			bs: "11111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111",
			x:  2, i: 61,
			want: "00000000000000000000000000000000000000000000000000000000000000101111111111111111111111111111111111111111111111111111111111111"},
		{
			bs: "1111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111",
			x:  15, i: 63,
			want: "0000000000000000000000000000000000000000000000000000000000001111111111111111111111111111111111111111111111111111111111111111111"},
		{
			bs: "000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
			x:  math.MaxUint64, i: 59,
			want: "111111111111111111111111111111111111111111111111111111111111111100000000000000000000000000000000000000000000000000000000000"},
		{
			bs: "001110101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010",
			x:  0xaaaaaaaaaaaaaaaa, i: 63,
			want: "001110101011010101010101010101010101010101010101010101010101010101010101010010101010101010101010101010101010101010101010101010101010101010"},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			bs, _ := MakeFromString(tt.bs)
			bs.SetUint64(tt.i, tt.x)
			want, _ := MakeFromString(tt.want)
			if !want.Equals(bs) {
				t.Errorf("Bitstring(%s).SetUint64(%d, %d) got %s, want %s",
					tt.bs, tt.i, tt.x, bs, want)
			}
		})
	}
}
