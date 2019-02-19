// Package bitstring implements a fixed length bit string type and bit string
// manipulation functions
package bitstring

import (
	"errors"
	"fmt"
	"math/big"
	"math/rand"
)

var (
	// ErrIndexOutOfRange is passed to panic if a bit index is out of range
	ErrIndexOutOfRange = errors.New("bitstring.Bitstring: index out of range")

	// ErrInvalidLength is returned when the provided bitstring length is
	// invalid.
	ErrInvalidLength = errors.New("bitstring.Bitstring: invalid length")
)

const (
	wordlenlog2 = 5
	wordlen     = 1 << wordlenlog2
)

// Bitstring implements a fixed-length bit string.
//
// Internally, bits are packed into an array of uint32. This implementation
// makes more efficient use of space than the alternative approach of using an
// array of booleans.
type Bitstring struct {
	// length in bits of the bit string
	length uint

	// bits are packed in an array of 32-bit ints.
	data []uint32
}

// New creates a bit string of the specified length (in bits) with all bits
// initially set to zero (off).
func New(length uint) (*Bitstring, error) {
	if length > 0 {
		bs := &Bitstring{
			length: length,
			data:   make([]uint32, (length+wordlen-1)/wordlen),
		}
		return bs, nil
	}
	return nil, ErrInvalidLength
}

// Random creates a Bitstring of the length l in which each bit is assigned a
// random value using rng.
//
// Random randomly sets the uint32 values of the underlying slice, so it should
// be faster than creating a bit string and then randomly setting each
// individual bits.
func Random(length uint, rng *rand.Rand) (*Bitstring, error) {
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
	used := uint32(length % wordlen)
	if used < wordlen {
		unused := wordlen - used
		mask := uint32(0xffffffff >> unused)
		bs.data[len(bs.data)-1] &= mask
	}
	return bs, nil
}

// MakeFromString returns the corresponding Bitstring for the given string of 1s
// and 0s in big endian order.
func MakeFromString(s string) (*Bitstring, error) {
	bs, err := New(uint(len(s)))
	if err != nil {
		return nil, err
	}

	for i, c := range s {
		switch c {
		case '0':
			continue
		case '1':
			bs.SetBit(uint(len(s)-i-1), true)
		default:
			return nil, fmt.Errorf("illegal character at position %v: %#U", i, c)
		}
	}
	return bs, nil
}

// Len returns the number of bits of bs.
func (bs *Bitstring) Len() int {
	return int(bs.length)
}

// Data returns the bitstring underlying slice
func (bs *Bitstring) Data() []uint32 {
	return bs.data
}

// Bit returns the bit at index i.
//
// Index 0 index is the index of the bit to look-up (0 is the least-significant bit).
// Returns a boolean indicating whether the bit is set or not.
//
// If index is negative or greater than bs.Len(), Bit will panic with
// ErrIndexOutOfRange.
func (bs *Bitstring) Bit(i uint) bool {
	bs.mustExist(i)

	word := uint32(i / wordlen)
	offset := uint32(i % wordlen)
	return (bs.data[word] & (1 << offset)) != 0
}

// SetBit sets the bit at index i. Index 0 is the LSB.
//
// If index is negative or greater than bs.Len(), SetBit will panic with
// ErrIndexOutOfRange.
func (bs *Bitstring) SetBit(i uint, v bool) {
	bs.mustExist(i)

	word := uint32(i / wordlen)
	offset := uint32(i % wordlen)
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
func (bs *Bitstring) FlipBit(i uint) {
	bs.mustExist(i)

	word := uint32(i / wordlen)
	offset := uint32(i % wordlen)
	bs.data[word] ^= (1 << offset)
}

// Ensures i is a valid index for bs, if the index is negative or greater than
// bs.length mustExist will panic with ErrIndexOutOfRange.
func (bs *Bitstring) mustExist(i uint) {
	if i >= bs.length {
		panic(ErrIndexOutOfRange)
	}
}

// OnesCount returns the number of one bits.
func (bs *Bitstring) OnesCount() uint {
	var count uint
	for _, x := range bs.data {
		for x != 0 {
			x &= (x - 1) // Unsets the least significant on bit.
			count++      // Count how many times we have to unset a bit before x equals zero.
		}
	}
	return count
}

// ZeroesCount returns the number of zero bits.
func (bs *Bitstring) ZeroesCount() uint {
	return bs.length - bs.OnesCount()
}

// BigInt returns the big.Int representation of bs.
func (bs *Bitstring) BigInt() *big.Int {
	bi := new(big.Int)
	if _, ok := bi.SetString(bs.String(), 2); !ok {
		// XXX: by design, this panic should only happen when something very
		// wrong happens. For bi.SetString to fail the string passed should
		// contain runes other than 0's and 1's, or be empty.
		// bs.String() guarantees the string is made of 0's and 1's, plus, of
		// all the ways to create a Bitstring none of them allows the bitstring
		// to be empty though one could still have a zero value, by doing
		// bitstring.Bitstring{}. If it panics in that case that's just fair...
		panic(fmt.Sprintf("couldn't convert bit string \"%s\" to big.Int", bs.String()))
	}
	return bi
}

// SwapRange efficiently swaps the same range of bits between bs and other.
//
// Both Bitstring should not necessarily have the same length but should contain
// the range of bits specified by start and length.
func SwapRange(bs1, bs2 *Bitstring, start, length uint) {
	bs1.mustExist(start)
	bs2.mustExist(start)

	word := start / wordlen
	partial := (int(wordlen) - int(start)) % wordlen
	if partial > 0 {
		bs1.swapBits(bs2, word, 0xffffffff<<uint32(wordlen-partial))
		word++
	}

	remain := int(length) - partial // can be negative
	stop := remain / wordlen
	for i := int(word); i < stop; i++ {
		bs1.data[i], bs2.data[i] = bs2.data[i], bs1.data[i]
	}

	remain %= wordlen
	if remain > 0 {
		bs1.swapBits(bs2, uint(len(bs1.data)-1), 0xffffffff>>uint32(wordlen-remain))
	}
}

// other is the Bitstring to exchange bits with.
// i is the index of the word that will be swapped between the two
// bit strings.
// mask is a mask that specifies which bits in the word will be swapped.
func (bs *Bitstring) swapBits(other *Bitstring, i uint, mask uint32) {
	preserveMask := ^mask
	preservedThis := bs.data[i] & preserveMask
	preservedThat := other.data[i] & preserveMask
	swapThis := bs.data[i] & mask
	swapThat := other.data[i] & mask
	bs.data[i] = preservedThis | swapThat
	other.data[i] = preservedThat | swapThis
}

// String returns a string representation of bs in big endian order.
func (bs *Bitstring) String() string {
	b := make([]byte, bs.length)
	for i := uint(0); i < bs.length; i++ {
		if bs.Bit(i) {
			b[bs.length-1-i] = '1'
		} else {
			b[bs.length-1-i] = '0'
		}
	}
	return string(b)
}

// Copy creates and returns a new Bitstring that is the exact copy of a source
// Bitstring.
func Copy(src *Bitstring) *Bitstring {
	dst := make([]uint32, len(src.data))
	copy(dst, src.data)
	return &Bitstring{
		length: src.length,
		data:   dst,
	}
}

// Equals returns true if other is a Bitstring instance and both bit strings are
// the same length with identical bits set/unset.
func (bs *Bitstring) Equals(other *Bitstring) bool {
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
