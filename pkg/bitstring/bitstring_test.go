package bitstring

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBitstringCreation(t *testing.T) {
	// Check that a bit string is constructed correctly, with
	// the correct length and all bits initially set to zero.
	bs := New(100)
	assert.Equalf(t, bs.Len(), 100, "want Bitstring length 100, got: %v", bs.Len())
	for i := 0; i < bs.Len(); i++ {
		assert.False(t, bs.Bit(uint(i)), "Bit ", i, " should not be set.")
	}
}

func TestBitstringCreateRandomBitstring(t *testing.T) {
	rng := rand.New(rand.NewSource(99))
	// Check that a random bit string of the correct length is constructed.
	bs := Random(100, rng)
	assert.Equalf(t, bs.Len(), 100, "want Bitstring length 100, got: %v", bs.Len())
}

func TestBitstringSetBits(t *testing.T) {
	// Make sure that bits are set correctly.
	bs := New(5)

	bs.SetBit(1)
	bs.SetBit(4)

	// Testing with non-symmetrical string to ensure that there
	// aren't any problem with endianness nor indices.
	assert.False(t, bs.Bit(0), "Bit 0 should not be set.")
	assert.True(t, bs.Bit(1), "Bit 1 should be set.")
	assert.False(t, bs.Bit(2), "Bit 2 should not be set.")
	assert.False(t, bs.Bit(3), "Bit 3 should not be set.")
	assert.True(t, bs.Bit(4), "Bit 4 should be set.")

	bs.ClearBit(4)
	assert.False(t, bs.Bit(4), "Bit 4 should be unset.")
}

func TestBitstringFlipBits(t *testing.T) {
	// Make sure bit-flipping works as expected.
	bs := New(5)

	bs.FlipBit(2)
	assert.True(t, bs.Bit(2), "Flipping unset bit failed.")

	bs.FlipBit(2)
	assert.False(t, bs.Bit(2), "Flipping set bit failed.")
}

func TestBitstringToString(t *testing.T) {
	// Checks that string representations are correctly generated.
	bs := New(10)

	bs.SetBit(3)
	bs.SetBit(7)
	bs.SetBit(8)

	// Testing with leading zero to make sure it isn't omitted.
	want := "0110001000"
	got := bs.String()
	assert.Equalf(t, want, got, "Incorrect string representation, want %s, got: %s", want, got)
}

func TestBitstringParsing(t *testing.T) {
	// Checks that the String-parsing constructor works correctly.
	// Use a 33-bit string to check that word boundaries are dealt with correctly.
	want := "111010101110101100010100101000101"
	bs, err := MakeFromString(want)
	assert.NoError(t, err)
	got := bs.String()
	assert.Equalf(t, want, got, "Failed parsing: want %v, got %v", want, got)
}

// Checks that integer conversion is correct.
func TestBitstringToNumber(t *testing.T) {
	bs := New(10)

	bs.SetBit(0)
	bs.SetBit(9)
	bint := bs.BigInt()
	assert.True(t, bint.IsInt64())
	assert.EqualValuesf(t, 513, bint.Int64(), "Incorrect big.Int conversion, want %v, got: %v", 513, bint.Int64())
}

func TestBitstringCountSetBits(t *testing.T) {
	// Checks that the bit string can correctly count its number of set bits.
	bs := New(64)
	assert.Zerof(t, bs.OnesCount(), "Initial string should have no 1s, got: %v, repr \"%v\"", bs.OnesCount(), bs)

	// The bits to set have been chosen because they deal with boundary cases.
	bs.SetBit(0)
	bs.SetBit(31)
	bs.SetBit(32)
	bs.SetBit(33)
	bs.SetBit(63)
	setBits := bs.OnesCount()
	assert.EqualValuesf(t, 5, setBits, "want set bits = 5, got: %v", setBits)
}

// Checks that the bit string can correctly count its number of unset bits.
func TestBitstringCountUnsetBits(t *testing.T) {
	bs := New(12)
	assert.EqualValuesf(t, 12, bs.ZeroesCount(), "Initial string should have no 1s, got: %v, repr \"%v\"", bs.ZeroesCount(), bs)

	bs.SetBit(0)
	bs.SetBit(5)
	bs.SetBit(6)
	bs.SetBit(9)
	bs.SetBit(10)
	setBits := bs.ZeroesCount()
	assert.EqualValuesf(t, 7, setBits, "want set bits = 7, got: %v", setBits)
}

func TestBitstringCopy(t *testing.T) {
	rng := rand.New(rand.NewSource(99))
	bs := Random(2000, rng)
	bs.SetBit(3)
	bs.SetBit(7)
	bs.SetBit(8)
	cpy := Copy(bs)

	// Check the copy is a bit-for-bit duplicate.
	for i := uint(0); i < uint(bs.Len()); i++ {
		assert.Equalf(t, bs.Bit(i), cpy.Bit(i), "copy differs original at bit %d", i)
	}

	assert.False(t, cpy == bs, "copy and original should be different bitstring instances")
	bs.FlipBit(2)
	assert.False(t, cpy.Bit(2) == bs.Bit(2), "copy and original are not independent from each other")
}

func TestBitstringEquality(t *testing.T) {
	bs := New(10)
	bs.SetBit(2)
	bs.SetBit(5)
	bs.SetBit(8)

	assert.True(t, bs.Equals(bs), "Bitstring should always equals itself.")
	assert.False(t, bs.Equals(nil), "Valid Bitstring should never equals nil.")
	assert.False(t, bs.Equals(&Bitstring{}), "Bitstring should not equals another instance")

	cpy := Copy(bs)
	assert.Truef(t, cpy.Equals(bs), "Freshly copied Bitstring should equals original, bs=%s copy=%s", bs, cpy)

	// Changing one of the objects should result in them no longer being
	// considered equal.
	cpy.FlipBit(0)
	assert.Falsef(t, cpy.Equals(bs), "bs=%s copy=%s", bs, cpy)

	// Bit strings of different lengths but with the same bits set should not
	// be considered equal.
	bs2 := New(9)
	bs2.SetBit(2)
	bs2.SetBit(5)
	bs2.SetBit(8)
	assert.False(t, bs2.Equals(bs))
}

func TestBitstringFromString(t *testing.T) {
	tests := []struct {
		name  string
		str   string
		valid bool
	}{
		{"invalid ascii chars", "0101012", false},
		{"non ascii chars", "10日本", false},
		{"only 0s", "0000", true},
		{"only 1s", "1111111", true},
		{"mixed 0s and 1s", "1000111011", true},
		{"empty string", "", true},
		{"with spaces", "11 ", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MakeFromString(tt.str)
			if tt.valid {
				assert.NotNil(t, got)
				assert.NoError(t, err)
			} else {
				assert.Nil(t, got)
				assert.Error(t, err)
			}
		})
	}
}

func TestBitstringSetBit(t *testing.T) {
	bs := New(1)
	t.Run("panics on index too high", func(t *testing.T) {
		// The index of an individual bit must be within the range 0 to
		// length-1.
		assert.Panics(t, func() { bs.ClearBit(1) })
	})
}

func TestBitstringSwapRange(t *testing.T) {
	tests := []struct {
		x, y          string
		start, length uint
		wantx, wanty  string
	}{
		{
			x:     "1",
			y:     "0",
			start: 0, length: 1,
			wantx: "0",
			wanty: "1",
		},
		{
			x:     "1111111111111111111111111111111111111111111111111111111111111111",
			y:     "0000000000000000000000000000000000000000000000000000000000000000",
			start: 0, length: 32,
			wantx: "1111111111111111111111111111111100000000000000000000000000000000",
			wanty: "0000000000000000000000000000000011111111111111111111111111111111",
		},
		{
			x:     "1111111111111111111111111111111111111111",
			y:     "0000000000000000000000000000000000000000",
			start: 2, length: 30,
			wantx: "1111111100000000000000000000000000000011",
			wanty: "0000000011111111111111111111111111111100",
		},
		{
			x:     "1111111111",
			y:     "0000000000",
			start: 0, length: 3,
			wantx: "1111111000",
			wanty: "0000000111",
		},
		{
			x:     "111",
			y:     "000",
			start: 1, length: 2,
			wantx: "001",
			wanty: "110",
		},
		{
			x:     "111",
			y:     "000",
			start: 0, length: 3,
			wantx: "000",
			wanty: "111",
		},
		{
			x:     "11111111111111111111111111111111",
			y:     "00000000000000000000000000000000",
			start: 0, length: 32,
			wantx: "00000000000000000000000000000000",
			wanty: "11111111111111111111111111111111",
		},
		{
			x:     "111111111111111111111111111111111",
			y:     "000000000000000000000000000000000",
			start: 0, length: 33,
			wantx: "000000000000000000000000000000000",
			wanty: "111111111111111111111111111111111",
		},
		{
			x:     "111111111111111111111111111111111111111111111111111111111111111",
			y:     "000000000000000000000000000000000000000000000000000000000000000",
			start: 0, length: 63,
			wantx: "000000000000000000000000000000000000000000000000000000000000000",
			wanty: "111111111111111111111111111111111111111111111111111111111111111",
		},
		{
			x:     "1111111111111111111111111111111111111111111111111111111111111111",
			y:     "0000000000000000000000000000000000000000000000000000000000000000",
			start: 0, length: 64,
			wantx: "0000000000000000000000000000000000000000000000000000000000000000",
			wanty: "1111111111111111111111111111111111111111111111111111111111111111",
		},
		{
			x:     "11111111111111111111111111111111111111111111111111111111111111111",
			y:     "00000000000000000000000000000000000000000000000000000000000000000",
			start: 0, length: 65,
			wantx: "00000000000000000000000000000000000000000000000000000000000000000",
			wanty: "11111111111111111111111111111111111111111111111111111111111111111",
		},
		{
			x:     "1111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111",
			y:     "0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
			start: 94, length: 1,
			wantx: "1101111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111",
			wanty: "0010000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
		},
		{
			x:     "111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111",
			y:     "000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
			start: 1, length: 256,
			wantx: "100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001",
			wanty: "011111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111110",
		},
		{
			x:     "1111111111111111111111111111111111111111111111111111111111111111111",
			y:     "0000000000000000000000000000000000000000000000000000000000000000000",
			start: 64, length: 2,
			wantx: "1001111111111111111111111111111111111111111111111111111111111111111",
			wanty: "0110000000000000000000000000000000000000000000000000000000000000000",
		},
		{
			x:     "1111111111111111111111111111111111111111111111111111111111111111111",
			y:     "0000000000000000000000000000000000000000000000000000000000000000000",
			start: 65, length: 1,
			wantx: "1011111111111111111111111111111111111111111111111111111111111111111",
			wanty: "0100000000000000000000000000000000000000000000000000000000000000000",
		},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			x, err1 := MakeFromString(tt.x)
			assert.NoError(t, err1)
			y, err2 := MakeFromString(tt.y)
			assert.NoError(t, err2)
			SwapRange(x, y, tt.start, tt.length)

			assert.Equalf(t, tt.wantx, x.String(),
				"want %s, got %s", tt.wantx, x.String())
			assert.Equalf(t, tt.wanty, y.String(),
				"want %s, got %s", tt.wanty, y.String())
		})
	}
}
