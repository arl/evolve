package bitstring

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBitStringCreation(t *testing.T) {
	// Check that a bit string is constructed correctly, with
	// the correct length and all bits initially set to zero.
	bitString, err := New(100)
	assert.NoError(t, err)
	assert.Equalf(t, bitString.Len(), 100, "want BitString length 100, got: %v", bitString.Len())
	for i := 0; i < bitString.Len(); i++ {
		assert.False(t, bitString.Bit(i), "Bit ", i, " should not be set.")
	}
}

func TestBitStringCreateRandomBitString(t *testing.T) {
	rng := rand.New(rand.NewSource(99))
	// Check that a random bit string of the correct length is constructed.
	bitString, err := NewRandom(100, rng)
	assert.NoError(t, err)
	assert.Equalf(t, bitString.Len(), 100, "want BitString length 100, got: %v", bitString.Len())
}

func TestBitStringSetBits(t *testing.T) {
	// Make sure that bits are set correctly.
	bitString, err := New(5)
	assert.NoError(t, err)

	bitString.SetBit(1, true)
	bitString.SetBit(4, true)

	// Testing with non-symmetrical string to ensure that there are no endian
	// problems.
	assert.False(t, bitString.Bit(0), "Bit 0 should not be set.")
	assert.True(t, bitString.Bit(1), "Bit 1 should be set.")
	assert.False(t, bitString.Bit(2), "Bit 2 should not be set.")
	assert.False(t, bitString.Bit(3), "Bit 3 should not be set.")
	assert.True(t, bitString.Bit(4), "Bit 4 should be set.")

	// Test unsetting a bit.
	bitString.SetBit(4, false)
	assert.False(t, bitString.Bit(4), "Bit 4 should be unset.")
}

func TestBitStringFlipBits(t *testing.T) {
	// Make sure bit-flipping works as expected.
	bitString, err := New(5)
	assert.NoError(t, err)

	bitString.FlipBit(2)
	assert.True(t, bitString.Bit(2), "Flipping unset bit failed.")

	bitString.FlipBit(2)
	assert.False(t, bitString.Bit(2), "Flipping set bit failed.")
}

func TestBitStringToString(t *testing.T) {
	// Checks that string representations are correctly generated.
	bitString, err := New(10)
	assert.NoError(t, err)

	bitString.SetBit(3, true)
	bitString.SetBit(7, true)
	bitString.SetBit(8, true)

	// Testing with leading zero to make sure it isn't omitted.
	exp := "0110001000"
	got := bitString.String()
	assert.Equalf(t, exp, got, "Incorrect string representation, want %s, got: %s", exp, got)
}

func TestBitStringParsing(t *testing.T) {
	// Checks that the String-parsing constructor works correctly.
	// Use a 33-bit string to check that word boundaries are dealt with correctly.
	fromString := "111010101110101100010100101000101"
	bitString, err := NewFromString(fromString)
	assert.NoError(t, err)
	toString := bitString.String()
	assert.Equal(t, toString, fromString, "Failed parsing: String representations do not match.")
}

// Checks that integer conversion is correct.
func TestBitStringToNumber(t *testing.T) {
	bitString, err := New(10)
	assert.NoError(t, err)

	bitString.SetBit(0, true)
	bitString.SetBit(9, true)
	bint := bitString.ToBigInt()
	assert.True(t, bint.IsInt64())
	assert.EqualValuesf(t, 513, bint.Int64(), "Incorrect big.Int conversion, want %v, got: %v", 513, bint.Int64())
}

func TestBitStringCountSetBits(t *testing.T) {
	// Checks that the bit string can correctly count its number of set bits.
	bitString, err := New(64)
	assert.NoError(t, err)
	assert.Zerof(t, bitString.CountSetBits(), "Initial string should have no 1s, got: %v, repr \"%v\"", bitString.CountSetBits(), bitString)

	// The bits to set have been chosen because they deal with boundary cases.
	bitString.SetBit(0, true)
	bitString.SetBit(31, true)
	bitString.SetBit(32, true)
	bitString.SetBit(33, true)
	bitString.SetBit(63, true)
	setBits := bitString.CountSetBits()
	assert.Equalf(t, 5, setBits, "want set bits = 5, got: %v", setBits)
}

// Checks that the bit string can correctly count its number of unset bits.
func TestBitStringCountUnsetBits(t *testing.T) {
	bitString, err := New(12)
	assert.NoError(t, err)
	assert.Equalf(t, 12, bitString.CountUnsetBits(), "Initial string should have no 1s, got: %v, repr \"%v\"", bitString.CountUnsetBits(), bitString)

	bitString.SetBit(0, true)
	bitString.SetBit(5, true)
	bitString.SetBit(6, true)
	bitString.SetBit(9, true)
	bitString.SetBit(10, true)
	setBits := bitString.CountUnsetBits()
	assert.Equalf(t, 7, setBits, "want set bits = 7, got: %v", setBits)
}

func TestBitStringClone(t *testing.T) {
	bitString, err := New(10)
	assert.NoError(t, err)
	bitString.SetBit(3, true)
	bitString.SetBit(7, true)
	bitString.SetBit(8, true)
	clone := bitString.Clone()

	// Check the clone is a bit-for-bit duplicate.
	for i := 0; i < bitString.Len(); i++ {
		assert.Equalf(t, bitString.Bit(i), clone.Bit(i), "Cloned bit string does not match in position %v", i)
	}

	// Check that clone is distinct from original (i.e. it does not change
	// if the original is modified).
	assert.False(t, clone == bitString, "want clone and original different objects, got the same")
	bitString.FlipBit(2)
	assert.False(t, clone.Bit(2), "Clone is not independent from original.")
}

func TestBitStringEquality(t *testing.T) {
	bitString, err := New(10)
	assert.NoError(t, err)
	bitString.SetBit(2, true)
	bitString.SetBit(5, true)
	bitString.SetBit(8, true)

	assert.True(t, bitString.Equals(bitString), "BitString should always equal itself.")
	assert.False(t, bitString.Equals(nil), "Valid BitString should never equal nil.")
	assert.False(t, bitString.Equals(&BitString{}), "BitString should not equal another instance")

	clone := bitString.Clone()
	assert.True(t, clone.Equals(bitString), "Freshly cloned BitString should equal original")

	// Changing one of the objects should result in them no longer being
	// considered equal.
	clone.FlipBit(0)
	assert.False(t, clone.Equals(bitString), "want different strings to cancel equality, \"%v\" and \"%s\"", clone, bitString)

	// Bit strings of different lengths but with the same bits set should not
	// be considered equal.
	var shortBitString *BitString
	shortBitString, err = New(9)
	assert.NoError(t, err)
	shortBitString.SetBit(2, true)
	shortBitString.SetBit(5, true)
	shortBitString.SetBit(8, true)
	assert.False(t, shortBitString.Equals(bitString), "want equal numbers but of different lengths to be considered not equal")
}

func TestBitStringInvalidLength(t *testing.T) {
	// The length of a bit string must be non-negative. If an attempt is made to
	// create a bit string with a negative length, an error and a nil BitString
	// pointer should be returned.
	bitString, err := New(-1)
	assert.Nil(t, bitString)
	assert.Error(t, err)
}

func TestBitStringFromString(t *testing.T) {
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
			got, err := NewFromString(tt.str)
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

func TestBitStringSetBit(t *testing.T) {
	bitString, err := New(1)
	assert.NoError(t, err)

	t.Run("panics on negative index", func(t *testing.T) {
		// The index of an individual bit must be non-negative.
		assert.Panics(t, func() {
			bitString.SetBit(-1, false)
		})
	})

	t.Run("panics on index too high", func(t *testing.T) {
		// The index of an individual bit must be within the range 0 to
		// length-1.
		assert.Panics(t, func() {
			bitString.SetBit(1, false)
		})
	})
}

func TestBitStringSwapSubstring(t *testing.T) {
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
			ones, err1 := NewFromString(tt.ones)
			assert.NoError(t, err1)
			zeros, err2 := NewFromString(tt.zeros)
			assert.NoError(t, err2)
			ones.SwapSubstring(zeros, tt.lo, tt.hi)

			assert.Equalf(t, tt.expOnes, ones.String(),
				"want %s, got %s", tt.expOnes, ones.String())
			assert.Equalf(t, tt.expZeros, zeros.String(),
				"want %s, got %s", tt.expZeros, zeros.String())
		})
	}
}

func ExampleNew() {
	// create a 8 bits BitString
	bitstring, _ := New(8)
	// upon creation all bits are unset
	fmt.Println(bitstring)
	// Output: 00000000
}

func ExampleNewFromString() {
	// create a BitString from string
	bitstring, _ := NewFromString("101001")
	fmt.Println(bitstring)
	// Output: 101001
}

func ExampleBitString_Len() {
	// create a 8 bits BitString
	bitstring, _ := New(8)
	fmt.Println(bitstring.Len())
	// Output: 8
}

func ExampleBitString_Bit() {
	// create a 8 bits BitString
	bitstring, _ := New(8)
	fmt.Println(bitstring.Bit(7))
	// Output: false
}

func ExampleBitString_SetBit() {
	// create a 8 bits BitString
	bitstring, _ := New(8)
	bitstring.SetBit(2, true)
	fmt.Println(bitstring)
	// Output: 00000100
}

func ExampleBitString_FlipBit() {
	// create a 8 bits BitString
	bitstring, _ := New(8)
	bitstring.FlipBit(2)
	fmt.Println(bitstring)
	// Output: 00000100
}

func ExampleBitString_CountUnsetBits() {
	// create a 8 bits BitString
	bitstring, _ := New(8)
	// upon creation all bits are unset
	fmt.Println(bitstring.CountUnsetBits())
	// Output: 8
}

func ExampleBitString_CountSetBits() {
	// create a 8 bits BitString
	bitstring, _ := New(8)
	// upon creation all bits are unset
	fmt.Println(bitstring.CountSetBits())
	// Output: 0
}

func ExampleBitString_ToBigInt() {
	// create a 8 bits BitString
	bitstring, _ := NewFromString("100")
	bi := bitstring.ToBigInt()
	fmt.Println(bi.Int64())
	// Output: 4
}

func ExampleBitString_SwapSubstring() {
	bs1, _ := NewFromString("111")
	bs2, _ := NewFromString("000")
	// starting from bit 2 of bs1, swap 1 bit with bs2
	bs1.SwapSubstring(bs2, 2, 1)
	fmt.Println(bs1)
	// Output: 011
}
