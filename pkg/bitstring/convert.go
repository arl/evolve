package bitstring

import (
	"fmt"
)

// Uintn returns the n-bit unsigned integer value represented by the n bits
// starting at the bit index i. It panics if there are not enough bits or if n
// is greater than 32.
// TODO: reverse order of nbits and i params
func (bs *Bitstring) Uintn(nbits, i uint) word {
	if nbits > wordlen || nbits < 1 {
		panic(fmt.Sprintf("Uintn supports unsigned integers from 1 to %d bits long", wordlen))
	}
	bs.mustExist(i + nbits - 1)

	j := wordoffset(i)
	k := wordoffset(i + nbits - 1)
	looff := bitoffset(i)
	loword := bs.data[j]
	if j == k {
		// fast path: same word
		return (loword >> looff) & genlomask(nbits)
	}
	hioff := bitoffset(i + nbits)
	hiword := bs.data[k] & genlomask(uint(hioff))
	loword = genhimask(uint(looff)) & loword >> looff
	return loword | hiword<<(wordlen-looff)
}

// Uint16 returns the uint16 value represented by the 16 bits starting at the
// given bit. It panics if there are not enough bits.
func (bs *Bitstring) Uint16(i uint) uint16 {
	bs.mustExist(i + 15)

	off := bitoffset(i)
	loword := bs.data[wordoffset(i)] >> off
	hiword := bs.data[wordoffset(i+15)] & ((1 << off) - 1)
	return uint16(loword | hiword<<(wordlen-off))
}

// Uint8 returns the uint8 value represented by the 8 bits starting at the
// given bit. It panics if there are not enough bits.
func (bs *Bitstring) Uint8(i uint) uint8 {
	bs.mustExist(i + 7)

	off := bitoffset(i)
	loword := bs.data[wordoffset(i)] >> off
	hiword := bs.data[wordoffset(i+7)] & ((1 << off) - 1)
	return uint8(loword | hiword<<(wordlen-off))
}

// Int32 returns the int32 value represented by the 32 bits starting at the
// given bit. It panics if there are not enough bits.
func (bs *Bitstring) Int32(i uint) int32 {
	return int32(bs.Uint32(i))
}

// Int16 returns the int16 value represented by the 16 bits starting at the
// given bit. It panics if there are not enough bits.
func (bs *Bitstring) Int16(i uint) int16 {
	return int16(bs.Uint16(i))
}

// Intn returns the n-bit signed integer value represented by the n bits
// starting at the i. It panics if there are not enough bits or if n is greater
// than 32.
func (bs *Bitstring) Intn(nbits, i uint) int32 {
	return int32(bs.Uintn(nbits, i))
}

// Int8 returns the int8 value represented by the 8 bits starting at the
// given bit. It panics if there are not enough bits.
func (bs *Bitstring) Int8(i uint) int8 {
	return int8(bs.Uint8(i))
}

// prints a string representing the first n bits of the base-2 representatio of x.
func printbits(x word, n uint) {
	fmt.Printf(fmt.Sprintf("%%0%db\n", n), x)
}
