// Package bitstring implmements a fixed length bit string type and common
// operations on bit strings.
package bitstring

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"math/big"
	"math/rand"
)

const wordLength = 32

// BitString implements a fixed-length bit-string.
//
// Internally, bits are packed into an array of ints. This implementation makes
// more efficient use of space than the alternative approach of using an array
// of booleans.
type BitString struct {
	// length in bits of the bit string
	length int

	// bits are packed in an array of 32-bit ints.
	data []uint32
}

// New creates a bit string of the specified length (in bits) with all bits
// initially set to zero (off).
func New(length int) (*BitString, error) {
	if length < 0 {
		return nil, errors.New("BitString length must be non-negative")
	}
	return &BitString{
		length: length,
		data:   make([]uint32, (length+wordLength-1)/wordLength),
	}, nil
}

// NewRandom creates a BitString of the specified length with each bit set
// randomly (the distribution of bits is uniform so long as the output
// from the provided pseudo-random number generator is also uniform).
//
// Using NewRandom is more efficient than creating a bit string and then
// randomly setting each bit individually.
func NewRandom(length int, rng *rand.Rand) (*BitString, error) {
	bt, err := New(length)
	if err != nil {
		return nil, err
	}

	// Instead of setting each bit with a random boolean, we create a random
	// byte buffer and fill the data field with it.
	buf := make([]byte, len(bt.data)*4)
	rng.Read(buf) // no need for error checking on Rand.Read
	for i := 0; i < len(bt.data); i++ {
		randomInt, _ := binary.Varint(buf[i : i+4])
		bt.data[i] = uint32(randomInt)
	}

	// If the last word is not fully utilised, zero any out-of-bounds bits.
	// This is necessary because the CountSetBits() methods will count
	// out-of-bounds bits.
	bitsUsed := uint32(length % wordLength)
	if bitsUsed < wordLength {
		unusedBits := wordLength - bitsUsed
		mask := uint32(0xffffffff >> unusedBits)
		bt.data[len(bt.data)-1] &= mask
	}
	return bt, nil
}

// NewFromString creates a new BitString from a character string of 1s and 0s
// in big-endian order.
func NewFromString(value string) (*BitString, error) {
	bt, err := New(len(value))
	if err != nil {
		return nil, err
	}

	for i, c := range value {
		switch c {
		case '0':
			continue
		case '1':
			bt.SetBit(len(value)-i-1, true)
		default:
			return nil, fmt.Errorf("illegal character at position %v: %#U", i, c)
		}
	}
	return bt, nil
}

// Len returns the length of the BitString in bits.
func (bt *BitString) Len() int {
	return bt.length
}

// Bit returns the bit at the specified index.
//
//  - index is the index of the bit to look-up (0 is the least-significant
//  bit).
//
// Returns a boolean indicating whether the bit is set or not.
//
// Will panic if the specified index is not a bit position in this bit string.
func (bt *BitString) Bit(index int) bool {
	bt.assertValidIndex(index)
	word := uint32(index / wordLength)
	offset := uint32(index % wordLength)
	return (bt.data[word] & (1 << offset)) != 0
}

// SetBit sets the bit at the specified index.
//
//  - index is the index of the bit to set (0 is the least-significant bit).
//  - set is a boolean indicating whether the bit should be set or not.
//
// Will panic if the specified index is not a bit position in this bit string.
func (bt *BitString) SetBit(index int, set bool) {
	bt.assertValidIndex(index)
	word := uint32(index / wordLength)
	offset := uint32(index % wordLength)
	if set {
		bt.data[word] |= (1 << offset)
	} else {
		// Unset the bit.
		bt.data[word] &= ^(1 << offset)
	}
}

// FlipBit inverts the value of the bit at the specified index.
//
// - param index is the bit to flip (0 is the least-significant bit).
//
// Will panic if the specified index is not a bit position in this bit string.
func (bt *BitString) FlipBit(index int) {
	bt.assertValidIndex(index)
	word := uint32(index / wordLength)
	offset := uint32(index % wordLength)
	bt.data[word] ^= (1 << offset)
}

// Helper method to check whether a bit index is valid or not.
// Will panic if the index is not valid.
func (bt *BitString) assertValidIndex(index int) {
	if index >= bt.length || index < 0 {
		panic(fmt.Sprintf("invalid index: %v (length: %v)", index, bt.length))
	}
}

// CountSetBits returns the number of bits that are 1s rather than 0s.
func (bt *BitString) CountSetBits() int {
	var count int
	for _, x := range bt.data {
		for x != 0 {
			x &= (x - 1) // Unsets the least significant on bit.
			count++      // Count how many times we have to unset a bit before x equals zero.
		}
	}
	return count
}

// CountUnsetBits returns the number of bits that are 0s rather than 1s.
func (bt *BitString) CountUnsetBits() int {
	return bt.length - bt.CountSetBits()
}

// ToBigInt interprets this bit string as being a binary numeric value and returns
// the integer that it represents.
//
// Returns a big.Int that contains the numeric value represented by this bit
// string.
func (bt *BitString) ToBigInt() *big.Int {
	bi := new(big.Int)
	if _, ok := bi.SetString(bt.String(), 2); !ok {
		panic(fmt.Sprintf("couldn't convert bit string \"%s\" to big.Int", bt.String()))
	}
	return bi
}

// SwapSubstring is an efficient method for exchanging data between two bit
// strings. Both bit strings must be long enough that they contain the full
// length of the specified substring.
//
//  - other is the bitstring with which this bitstring should swap bits.
//  - start is the start position for the substrings to be exchanged. All bit
//  indices are big-endian, which means position 0 is the rightmost bit.
//  - length is the number of contiguous bits to swap.
func (bt *BitString) SwapSubstring(other *BitString, start, length int) {
	bt.assertValidIndex(start)
	other.assertValidIndex(start)

	word := start / wordLength
	partialWordSize := (wordLength - start) % wordLength
	if partialWordSize > 0 {
		bt.swapBits(other, word, 0xffffffff<<uint32(wordLength-partialWordSize))
		word++
	}

	remainingBits := length - partialWordSize
	stop := remainingBits / wordLength
	for i := word; i < stop; i++ {
		bt.data[i], other.data[i] = other.data[i], bt.data[i]
	}

	remainingBits %= wordLength
	if remainingBits > 0 {
		bt.swapBits(other, len(bt.data)-1, 0xffffffff>>uint32(wordLength-remainingBits))
	}
}

//  - other is the BitString to exchange bits with.
//  - word is the word index of the word that will be swapped between the two
//  bit strings.
//  - swapMask is a mask that specifies which bits in the word will be swapped.
func (bt *BitString) swapBits(other *BitString, word int, swapMask uint32) {
	preserveMask := ^swapMask
	preservedThis := bt.data[word] & preserveMask
	preservedThat := other.data[word] & preserveMask
	swapThis := bt.data[word] & swapMask
	swapThat := other.data[word] & swapMask
	bt.data[word] = preservedThis | swapThat
	other.data[word] = preservedThat | swapThis
}

// String creates a textual representation of this bit string in big-endian
// order (index 0 is the right-most bit).
//
// Returns this bit string rendered as a string of 1s and 0s.
func (bt *BitString) String() string {
	buf := bytes.NewBuffer(make([]byte, 0, bt.length))
	for i := bt.length - 1; i >= 0; i-- {
		if bt.Bit(i) {
			buf.WriteString("1")
		} else {
			buf.WriteString("0")
		}
	}
	return buf.String()
}

// Clone returns an identical copy of this bit string.
func (bt *BitString) Clone() *BitString {
	clone, _ := New(bt.length)
	copy(clone.data, bt.data)
	return clone
}

// Equals returns true if other is a BitString instance and both bit strings are
// the same length with identical bits set/unset.
func (bt *BitString) Equals(other *BitString) bool {
	switch {
	case bt == other:
		return true
	case bt != nil && other == nil:
		break
	case bt == nil && other != nil:
		break
	case bt.length == other.length:
		for i, v := range bt.data {
			if v != other.data[i] {
				return false
			}
		}
		return true
	}
	return false
}
