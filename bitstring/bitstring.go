// Package bitstring implements a fixed length bit string type and bit string
// manipulation functions
package bitstring

import (
	"bytes"
	"errors"
	"fmt"
	"math/big"
	"math/rand"
)

var (
	// ErrIndexOutOfRange is passed to panic if a bit index is out of the valid
	// range for a a given BitString.
	ErrIndexOutOfRange = errors.New("bitstring.BitString: index out of range")

	ErrInvalidLength = errors.New("bitstring.BitString: invalid length")
)

const wordLength = 32

// BitString implements a fixed-length bit string.
//
// Internally, bits are packed into an array of uint32. This implementation
// makes more efficient use of space than the alternative approach of using an
// array of booleans.
type BitString struct {
	// length in bits of the bit string
	length int

	// bits are packed in an array of 32-bit ints.
	data []uint32
}

// New creates a bit string of the specified length (in bits) with all bits
// initially set to zero (off).
func New(length int) (*BitString, error) {
	var (
		bs  *BitString
		err = ErrInvalidLength
	)
	if length > 0 {
		slicelen := (length + wordLength - 1) / wordLength
		bs = &BitString{
			length: length,
			data:   make([]uint32, slicelen),
		}
		err = nil
	}
	return bs, err
}

// Random creates a BitString of the length l in which each bit is assigned a
// random value using rng.
//
// Random randomly sets the uint32 values of the underlying slice, so it should
// be faster than creating a bit string and then randomly setting each
// individual bits.
func Random(length int, rng *rand.Rand) (*BitString, error) {
	bs, err := New(length)
	if err != nil {
		return nil, err
	}

	// fill the slice with random values
	for i := 0; i < len(bs.data); i++ {
		bs.data[i] = rng.Uint32()
	}

	// If the last word is not fully utilised, zero any out-of-bounds bits.
	// This is necessary because OnesCount and ZeroesCount count the
	// out-of-bounds bits.
	bitsUsed := uint32(length % wordLength)
	if bitsUsed < wordLength {
		unusedBits := wordLength - bitsUsed
		mask := uint32(0xffffffff >> unusedBits)
		bs.data[len(bs.data)-1] &= mask
	}
	return bs, nil
}

// MakeFromString returns the corresponding BitString for the given string of 1s
// and 0s in big endian order.
func MakeFromString(from string) (*BitString, error) {
	bs, err := New(len(from))
	if err != nil {
		return nil, err
	}

	for i, c := range from {
		switch c {
		case '0':
			continue
		case '1':
			bs.SetBit(len(from)-i-1, true)
		default:
			return nil, fmt.Errorf("illegal character at position %v: %#U", i, c)
		}
	}
	return bs, nil
}

// Len returns the number of bits of bs.
func (bs *BitString) Len() int {
	return bs.length
}

// Data returns the bitstring underlying slice
func (bs *BitString) Data() []uint32 {
	return bs.data
}

// Bit returns the bit at index i.
//
// Index 0 index is the index of the bit to look-up (0 is the least-significant bit).
// Returns a boolean indicating whether the bit is set or not.
//
// If index is negative or greater than bs.Len(), Bit will panic with
// ErrIndexOutOfRange.
func (bs *BitString) Bit(i int) bool {
	bs.mustExist(i)

	word := uint32(i / wordLength)
	offset := uint32(i % wordLength)
	return (bs.data[word] & (1 << offset)) != 0
}

// SetBit sets the bit at index i. Index 0 is the LSB.
//
// If index is negative or greater than bs.Len(), SetBit will panic with
// ErrIndexOutOfRange.
func (bs *BitString) SetBit(i int, v bool) {
	bs.mustExist(i)

	word := uint32(i / wordLength)
	offset := uint32(i % wordLength)
	if v {
		bs.data[word] |= (1 << offset)
	} else {
		// Unset the bit.
		bs.data[word] &= ^(1 << offset)
	}
}

// FlipBit flips the bit at index i.
//
// param index is the bit to flip (0 is the least-significant bit).
//
// If index is negative or greater than bs.Len(), FlipBit will panic with
// ErrIndexOutOfRange.
func (bs *BitString) FlipBit(i int) {
	bs.mustExist(i)

	word := uint32(i / wordLength)
	offset := uint32(i % wordLength)
	bs.data[word] ^= (1 << offset)
}

// Ensures i is a valid index for bs, if the index is negative or greater than
// bs.length mustExist will panic with ErrIndexOutOfRange.
func (bs *BitString) mustExist(i int) {
	if i >= bs.length || i < 0 {
		panic(ErrIndexOutOfRange)
	}
}

// OnesCount returns the number of one bits.
func (bs *BitString) OnesCount() int {
	var count int
	for _, x := range bs.data {
		for x != 0 {
			x &= (x - 1) // Unsets the least significant on bit.
			count++      // Count how many times we have to unset a bit before x equals zero.
		}
	}
	return count
}

// ZeroesCount returns the number of zero bits.
func (bs *BitString) ZeroesCount() int {
	return bs.length - bs.OnesCount()
}

// BigInt returns the big.Int representation of bs.
func (bs *BitString) BigInt() *big.Int {
	bi := new(big.Int)
	if _, ok := bi.SetString(bs.String(), 2); !ok {
		// XXX: by design, this panic should only happen when something very
		// wrong happens. For bi.SetString to fail the string passed should
		// contain other runes other than 0's and 1's, or be empty.
		// bs.String() guarantees the string is made of 0's and 1's and all the
		// ways to create BitString prevent construction them with 0-length.
		panic(fmt.Sprintf("couldn't convert bit string \"%s\" to big.Int", bs.String()))
	}
	return bi
}

// SwapRange efficiently swaps the same range of bits between bs and other.
//
// Both BitString should not necessarily have the same length but should contain
// the range of bits specified by start and length.
func (bs *BitString) SwapRange(other *BitString, start, length int) {
	bs.mustExist(start)
	other.mustExist(start)

	word := start / wordLength
	partialWordSize := (wordLength - start) % wordLength
	if partialWordSize > 0 {
		bs.swapBits(other, word, 0xffffffff<<uint32(wordLength-partialWordSize))
		word++
	}

	remainingBits := length - partialWordSize
	stop := remainingBits / wordLength
	for i := word; i < stop; i++ {
		bs.data[i], other.data[i] = other.data[i], bs.data[i]
	}

	remainingBits %= wordLength
	if remainingBits > 0 {
		bs.swapBits(other, len(bs.data)-1, 0xffffffff>>uint32(wordLength-remainingBits))
	}
}

// other is the BitString to exchange bits with.
// word is the word index of the word that will be swapped between the two
// bit strings.
// swapMask is a mask that specifies which bits in the word will be swapped.
func (bs *BitString) swapBits(other *BitString, word int, swapMask uint32) {
	preserveMask := ^swapMask
	preservedThis := bs.data[word] & preserveMask
	preservedThat := other.data[word] & preserveMask
	swapThis := bs.data[word] & swapMask
	swapThat := other.data[word] & swapMask
	bs.data[word] = preservedThis | swapThat
	other.data[word] = preservedThat | swapThis
}

// String returns a string representation of bs in big endian order.
func (bs *BitString) String() string {
	buf := bytes.NewBuffer(make([]byte, 0, bs.length))
	for i := bs.length - 1; i >= 0; i-- {
		if bs.Bit(i) {
			buf.WriteString("1")
		} else {
			buf.WriteString("0")
		}
	}
	return buf.String()
}

// Copy returns an identical copy of bs.
//
// The new BitString is based off of a new backing array.
func (bs *BitString) Copy() *BitString {
	clone, _ := New(bs.length)
	copy(clone.data, bs.data)
	return clone
}

// Equals returns true if other is a BitString instance and both bit strings are
// the same length with identical bits set/unset.
func (bs *BitString) Equals(other *BitString) bool {
	switch {
	case bs == other:
		return true
	case bs != nil && other == nil:
		break
	case bs == nil && other != nil:
		break
	case bs.length == other.length:
		for i, v := range bs.data {
			if v != other.data[i] {
				return false
			}
		}
		return true
	}
	return false
}
