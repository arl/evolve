package bitstring

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBitstringCreation(t *testing.T) {
	// Check that a bit string is constructed correctly, with
	// the correct length and all bits initially set to zero.
	bs, err := New(100)
	assert.NoError(t, err)
	assert.Equalf(t, bs.Len(), 100, "want Bitstring length 100, got: %v", bs.Len())
	for i := 0; i < bs.Len(); i++ {
		assert.False(t, bs.Bit(i), "Bit ", i, " should not be set.")
	}
}

func TestBitstringCreateRandomBitstring(t *testing.T) {
	rng := rand.New(rand.NewSource(99))
	// Check that a random bit string of the correct length is constructed.
	bs, err := Random(100, rng)
	assert.NoError(t, err)
	assert.Equalf(t, bs.Len(), 100, "want Bitstring length 100, got: %v", bs.Len())
}

func TestBitstringSetBits(t *testing.T) {
	// Make sure that bits are set correctly.
	bs, err := New(5)
	assert.NoError(t, err)

	bs.SetBit(1, true)
	bs.SetBit(4, true)

	// Testing with non-symmetrical string to ensure that there are no endian
	// or index problems.
	assert.False(t, bs.Bit(0), "Bit 0 should not be set.")
	assert.True(t, bs.Bit(1), "Bit 1 should be set.")
	assert.False(t, bs.Bit(2), "Bit 2 should not be set.")
	assert.False(t, bs.Bit(3), "Bit 3 should not be set.")
	assert.True(t, bs.Bit(4), "Bit 4 should be set.")

	// Test unsetting a bit.
	bs.SetBit(4, false)
	assert.False(t, bs.Bit(4), "Bit 4 should be unset.")
}

func TestBitstringFlipBits(t *testing.T) {
	// Make sure bit-flipping works as expected.
	bs, err := New(5)
	assert.NoError(t, err)

	bs.FlipBit(2)
	assert.True(t, bs.Bit(2), "Flipping unset bit failed.")

	bs.FlipBit(2)
	assert.False(t, bs.Bit(2), "Flipping set bit failed.")
}

func TestBitstringToString(t *testing.T) {
	// Checks that string representations are correctly generated.
	bs, err := New(10)
	assert.NoError(t, err)

	bs.SetBit(3, true)
	bs.SetBit(7, true)
	bs.SetBit(8, true)

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
	bs, err := New(10)
	assert.NoError(t, err)

	bs.SetBit(0, true)
	bs.SetBit(9, true)
	bint := bs.BigInt()
	assert.True(t, bint.IsInt64())
	assert.EqualValuesf(t, 513, bint.Int64(), "Incorrect big.Int conversion, want %v, got: %v", 513, bint.Int64())
}

func TestBitstringCountSetBits(t *testing.T) {
	// Checks that the bit string can correctly count its number of set bits.
	bs, err := New(64)
	assert.NoError(t, err)
	assert.Zerof(t, bs.OnesCount(), "Initial string should have no 1s, got: %v, repr \"%v\"", bs.OnesCount(), bs)

	// The bits to set have been chosen because they deal with boundary cases.
	bs.SetBit(0, true)
	bs.SetBit(31, true)
	bs.SetBit(32, true)
	bs.SetBit(33, true)
	bs.SetBit(63, true)
	setBits := bs.OnesCount()
	assert.Equalf(t, 5, setBits, "want set bits = 5, got: %v", setBits)
}

// Checks that the bit string can correctly count its number of unset bits.
func TestBitstringCountUnsetBits(t *testing.T) {
	bs, err := New(12)
	assert.NoError(t, err)
	assert.Equalf(t, 12, bs.ZeroesCount(), "Initial string should have no 1s, got: %v, repr \"%v\"", bs.ZeroesCount(), bs)

	bs.SetBit(0, true)
	bs.SetBit(5, true)
	bs.SetBit(6, true)
	bs.SetBit(9, true)
	bs.SetBit(10, true)
	setBits := bs.ZeroesCount()
	assert.Equalf(t, 7, setBits, "want set bits = 7, got: %v", setBits)
}

func TestBitstringClone(t *testing.T) {
	bs, err := New(10)
	assert.NoError(t, err)
	bs.SetBit(3, true)
	bs.SetBit(7, true)
	bs.SetBit(8, true)
	clone := Copy(bs)

	// Check the clone is a bit-for-bit duplicate.
	for i := 0; i < bs.Len(); i++ {
		assert.Equalf(t, bs.Bit(i), clone.Bit(i), "Cloned bit string does not match in position %v", i)
	}

	// Check that clone is distinct from original (i.e. it does not change
	// if the original is modified).
	assert.False(t, clone == bs, "want clone and original different objects, got the same")
	bs.FlipBit(2)
	assert.False(t, clone.Bit(2), "Clone is not independent from original.")
}

func TestBitstringEquality(t *testing.T) {
	bs, err := New(10)
	assert.NoError(t, err)
	bs.SetBit(2, true)
	bs.SetBit(5, true)
	bs.SetBit(8, true)

	assert.True(t, bs.Equals(bs), "Bitstring should always equal itself.")
	assert.False(t, bs.Equals(nil), "Valid Bitstring should never equal nil.")
	assert.False(t, bs.Equals(&Bitstring{}), "Bitstring should not equal another instance")

	clone := Copy(bs)
	assert.True(t, clone.Equals(bs), "Freshly cloned Bitstring should equal original")

	// Changing one of the objects should result in them no longer being
	// considered equal.
	clone.FlipBit(0)
	assert.False(t, clone.Equals(bs), "want different strings to cancel equality, \"%v\" and \"%s\"", clone, bs)

	// Bit strings of different lengths but with the same bits set should not
	// be considered equal.
	var bs2 *Bitstring
	bs2, err = New(9)
	assert.NoError(t, err)
	bs2.SetBit(2, true)
	bs2.SetBit(5, true)
	bs2.SetBit(8, true)
	assert.False(t, bs2.Equals(bs), "want equal numbers but of different lengths to be considered not equal")
}

func TestBitstringInvalidLength(t *testing.T) {
	// The length of a bit string must be non-negative. If an attempt is made to
	// create a bit string with a negative length, an error and a nil Bitstring
	// pointer should be returned.
	bs, err := New(-1)
	assert.Nil(t, bs)
	assert.Error(t, err)
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
		{"empty string", "", false},
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
	bs, err := New(1)
	assert.NoError(t, err)

	t.Run("panics on negative index", func(t *testing.T) {
		// The index of an individual bit must be non-negative.
		assert.Panics(t, func() { bs.SetBit(-1, false) })
	})

	t.Run("panics on index too high", func(t *testing.T) {
		// The index of an individual bit must be within the range 0 to
		// length-1.
		assert.Panics(t, func() { bs.SetBit(1, false) })
	})
}

func TestBitstringSwapRange(t *testing.T) {
	tests := []struct {
		name              string
		ones, zeros       string
		lo, hi            int // Bit indices are little-endian, so position 0 is the rightmost bit.
		expOnes, expZeros string
	}{
		{
			"word-aligned",
			"1111111111111111111111111111111111111111111111111111111111111111",
			"0000000000000000000000000000000000000000000000000000000000000000",
			0, 32,
			"1111111111111111111111111111111100000000000000000000000000000000",
			"0000000000000000000000000000000011111111111111111111111111111111",
		},
		{
			"non-aligned start",
			"1111111111111111111111111111111111111111",
			"0000000000000000000000000000000000000000",
			2, 30,
			"1111111100000000000000000000000000000011",
			"0000000011111111111111111111111111111100",
		},
		{
			"non-aligned end",
			"1111111111",
			"0000000000",
			0, 3,
			"1111111000",
			"0000000111",
		},
		{
			"smaller than word length",
			"111",
			"000",
			1, 2,
			"001",
			"110",
		},
		{
			"smaller than word length, full swap",
			"111",
			"000",
			0, 3,
			"000",
			"111",
		},
		{
			"word length full swap",
			"11111111111111111111111111111111",
			"00000000000000000000000000000000",
			0, 32,
			"00000000000000000000000000000000",
			"11111111111111111111111111111111",
		},
		{
			"greater than word length full swap",
			"111111111111111111111111111111111",
			"000000000000000000000000000000000",
			0, 33,
			"000000000000000000000000000000000",
			"111111111111111111111111111111111",
		},
		{
			"smaller than 2 times word length full swap",
			"111111111111111111111111111111111111111111111111111111111111111",
			"000000000000000000000000000000000000000000000000000000000000000",
			0, 63,
			"000000000000000000000000000000000000000000000000000000000000000",
			"111111111111111111111111111111111111111111111111111111111111111",
		},
		{
			"2 times word length full swap",
			"1111111111111111111111111111111111111111111111111111111111111111",
			"0000000000000000000000000000000000000000000000000000000000000000",
			0, 64,
			"0000000000000000000000000000000000000000000000000000000000000000",
			"1111111111111111111111111111111111111111111111111111111111111111",
		},
		{
			"greater than 2 times word length full swap",
			"11111111111111111111111111111111111111111111111111111111111111111",
			"00000000000000000000000000000000000000000000000000000000000000000",
			0, 65,
			"00000000000000000000000000000000000000000000000000000000000000000",
			"11111111111111111111111111111111111111111111111111111111111111111",
		},
		{
			"greater than 3 times word length",
			"1111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111",
			"0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
			94, 1,
			"0111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111",
			"1000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ones, err1 := MakeFromString(tt.ones)
			assert.NoError(t, err1)
			zeros, err2 := MakeFromString(tt.zeros)
			assert.NoError(t, err2)
			SwapRange(ones, zeros, tt.lo, tt.hi)

			assert.Equalf(t, tt.expOnes, ones.String(),
				"want %s, got %s", tt.expOnes, ones.String())
			assert.Equalf(t, tt.expZeros, zeros.String(),
				"want %s, got %s", tt.expZeros, zeros.String())
		})
	}
}
