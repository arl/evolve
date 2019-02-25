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

// Bitstring implements a fixed-length bit string.
//
// Internally, bits are packed into an array of machine word integers. This
// implementation makes more efficient use of space than the alternative
// approach of using an array of booleans.
type Bitstring struct {
	// length in bits of the bit string
	length uint

	// bits are packed in an array of 32-bit ints.
	data []word
}

// New creates a bit string of the specified length (in bits) with all bits
// initially set to zero (off).
func New(length uint) *Bitstring {
	bs := &Bitstring{
		length: length,
		data:   make([]word, (length+wordlen-1)/wordlen),
	}
	return bs
}

// Random creates a Bitstring of the length l in which each bit is assigned a
// random value using rng.
//
// Random randomly sets the uint32 values of the underlying slice, so it should
// be faster than creating a bit string and then randomly setting each
// individual bits.
func Random(length uint, rng *rand.Rand) *Bitstring {
	bs := New(length)

	a := bs.data[:len(bs.data)] // remove bounds-checking

	// fill words with random values
	switch wordlen {
	case 32:
		for i := range a {
			a[i] = word(rng.Uint32())
		}
	case 64:
		for i := range a {
			a[i] = word(rng.Uint64())
		}
	}

	// If the last word is not fully utilised, zero any out-of-bounds bits.
	// This is necessary because OnesCount and ZeroesCount count the
	// out-of-bounds bits.
	nused := bitoffset(length)
	if nused != 0 {
		mask := genlomask(uint(nused))
		a[len(a)-1] &= mask
	}
	return bs
}

// MakeFromString returns the corresponding Bitstring for the given string of 1s
// and 0s in big endian order.
func MakeFromString(s string) (*Bitstring, error) {
	bs := New(uint(len(s)))

	for i, c := range s {
		switch c {
		case '0':
			continue
		case '1':
			bs.SetBit(uint(len(s) - i - 1))
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
func (bs *Bitstring) Data() []word {
	return bs.data
}

// Bit returns a boolean indicating wether the bit at index i is set or not.
//
// If i is greater than bs.Len(), Bit will panic.
func (bs *Bitstring) Bit(i uint) bool {
	bs.mustExist(i)
	w := wordoffset(i)
	off := bitoffset(i)
	mask := bitmask(uint(off))
	return (bs.data[w] & mask) != 0
}

// SetBit sets the bit at index i.
//
// If i is greater than bs.Len(), SetBit will panic.
func (bs *Bitstring) SetBit(i uint) {
	bs.mustExist(i)

	w := wordoffset(i)
	off := bitoffset(i)
	bs.data[w] |= bitmask(uint(off))
}

// ClearBit clears the bit at index i.
//
// If i is greater than bs.Len(), ClearBit will panic.
func (bs *Bitstring) ClearBit(i uint) {
	bs.mustExist(i)

	w := wordoffset(i)
	off := bitoffset(i)
	bs.data[w] &= ^bitmask(uint(off))
}

// FlipBit flips (i.e toggles) the bit at index i.
//
// param index is the bit to flip (0 is the least-significant bit).
//
// If i is greater than bs.Len(), FlipBit will panic.
func (bs *Bitstring) FlipBit(i uint) {
	bs.mustExist(i)

	w := wordoffset(i)
	off := bitoffset(i)
	bs.data[w] ^= (1 << off)
}

// mustExist panics if i is not a valid bit index for bs, that is if i is
// greater than bs.length.
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

func minuint(x, y uint) uint {
	if x < y {
		return x
	}
	return y
}

// SwapRange swaps the same range of bits between bs and other.
//
// Both Bitstring may have different length the bit with index start+lenght must
// be valid on both or SwapRange will panic.
func SwapRange(bs1, bs2 *Bitstring, start, length uint) {
	bs1.mustExist(start + length - 1)
	bs2.mustExist(start + length - 1)

	// swap the required bits of the first word
	i := word(start / wordlen)
	start = uint(bitoffset(start))
	end := minuint(start+length, wordlen)
	remain := length - (end - start)
	swapBits(bs1, bs2, i, genmask(start, end))
	i++

	// swap whole words but the last one
	for remain > wordlen {
		bs1.data[i], bs2.data[i] = bs2.data[i], bs1.data[i]
		remain -= wordlen
		i++
	}

	// swap the remaining bits of the last word
	if remain != 0 {
		swapBits(bs1, bs2, i, genlomask(remain))
	}
}

// swapBits swaps range of bits from one word to another.
// w is the index of the word containing the bits to swap, and m is a mask that specifies
// whilch bits of that word will be swapped.
func swapBits(x, y *Bitstring, w, mask word) {
	keep := ^mask
	xkeep, ykeep := x.data[w]&keep, y.data[w]&keep
	xswap, yswap := x.data[w]&mask, y.data[w]&mask
	x.data[w] = xkeep | yswap
	y.data[w] = ykeep | xswap
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

// Copy creates and returns a new Bitstring that is a copy of src.
func Copy(src *Bitstring) *Bitstring {
	dst := make([]word, len(src.data))
	copy(dst, src.data)
	return &Bitstring{
		length: src.length,
		data:   dst,
	}
}

// Equals returns true if bs and other have the same lenght and each bit are
// identical, or if bs and other both point to the same Bitstring instance
// (i.e pointer equality).
func (bs *Bitstring) Equals(other *Bitstring) bool {
	switch {
	case bs == other:
		return true
	case bs != nil && other == nil,
		bs == nil && other != nil:
		return false
	case bs.length == other.length:
		od := other.data[:len(bs.data)] // remove BCE
		for i, v := range bs.data {
			if v != od[i] {
				return false
			}
		}
		return true
	}
	return false
}
