package bitstring

import (
	"fmt"
)

// Uintn returns the n bits unsigned integer value represented by the n bits
// starting at the bit index i. It panics if there aren't enough bits in bs or
// if n is greater than the size of a machine word.
// TODO: reverse order of nbits and i params
func (bs *Bitstring) Uintn(n, i uint) uint {
	if n > uintsize || n < 1 {
		panic(fmt.Sprintf("Uintn supports unsigned integers from 1 to %d bits long", uintsize))
	}
	bs.mustExist(i + n - 1)

	j := wordoffset(i)
	k := wordoffset(i + n - 1)
	looff := bitoffset(i)
	loword := bs.data[j]
	if j == k {
		// fast path: same word
		return (loword >> looff) & genlomask(n)
	}
	hioff := bitoffset(i + n)
	hiword := bs.data[k] & genlomask(uint(hioff))
	loword = genhimask(uint(looff)) & loword >> looff
	return loword | hiword<<(uintsize-looff)
}

// Uint16 returns the uint16 value represented by the 16 bits starting at the
// given bit. It panics if there are not enough bits.
func (bs *Bitstring) Uint16(i uint) uint16 {
	bs.mustExist(i + 15)

	off := bitoffset(i)
	loword := bs.data[wordoffset(i)] >> off
	hiword := bs.data[wordoffset(i+15)] & ((1 << off) - 1)
	return uint16(loword | hiword<<(uintsize-off))
}

// Uint8 returns the uint8 value represented by the 8 bits starting at the
// given bit. It panics if there are not enough bits.
func (bs *Bitstring) Uint8(i uint) uint8 {
	bs.mustExist(i + 7)

	off := bitoffset(i)
	loword := bs.data[wordoffset(i)] >> off
	hiword := bs.data[wordoffset(i+7)] & ((1 << off) - 1)
	return uint8(loword | hiword<<(uintsize-off))
}

// Intn returns the n-bit signed integer value represented by the n bits
// starting at the i. It panics if there are not enough bits or if n is greater
// than the size of a machine word.
func (bs *Bitstring) Intn(nbits, i uint) int32 { return int32(bs.Uintn(nbits, i)) }

// Int64 returns the int64 value represented by the 64 bits starting at the
// given bit. It panics if there are not enough bits.
func (bs *Bitstring) Int64(i uint) int64 { return int64(bs.Uint64(i)) }

// Int32 returns the int32 value represented by the 32 bits starting at the
// given bit. It panics if there are not enough bits.
func (bs *Bitstring) Int32(i uint) int32 { return int32(bs.Uint32(i)) }

// Int16 returns the int16 value represented by the 16 bits starting at the
// given bit. It panics if there are not enough bits.
func (bs *Bitstring) Int16(i uint) int16 { return int16(bs.Uint16(i)) }

// Int8 returns the int8 value represented by the 8 bits starting at the
// given bit. It panics if there are not enough bits.
func (bs *Bitstring) Int8(i uint) int8 { return int8(bs.Uint8(i)) }

// SetUintn sets the n bits starting at i with the first n bits of value x.
// It panics if there aren't enough bits in bs or if n is greater than
// the size of a machine word.
func (bs *Bitstring) SetUintn(n, i uint, x uint) {
	if n > uintsize || n < 1 {
		panic(fmt.Sprintf("SetUintn supports unsigned integers from 1 to %d bits long", uintsize))
	}
	bs.mustExist(i + n - 1)

	lobit := uint(bitoffset(i))
	j := wordoffset(i)
	k := wordoffset(i + n - 1)
	if j == k {
		// fast path: same word
		x := (x & genlomask(n)) << lobit
		bs.data[j] = transferbits(bs.data[j], x, genmask(lobit, lobit+n))
		return
	}
	// slow path: first and last bits are on different words
	// transfer bits to low word
	lon := uintsize - lobit // how many bits of n we transfer to loword
	bs.data[j] = transferbits(bs.data[j], x<<lobit, genhimask(lon))

	// transfer bits to high word
	bs.data[k] = transferbits(bs.data[k], x>>lon, genlomask(n-lon))
}

// SetUint8 sets the 8 bits starting at i with the value of x. It panics if
// there are not enough bits.
func (bs *Bitstring) SetUint8(i uint, x uint8) {
	bs.mustExist(i + 7)

	lobit := bitoffset(i)
	j := wordoffset(i)
	k := wordoffset(i + 7)
	if j == k {
		// fast path: same word
		neww := uint(x) << lobit
		mask := genmask(lobit, lobit+8)
		bs.data[j] = transferbits(bs.data[j], neww, mask)
		return
	}
	// transfer bits to low word
	bs.data[j] = transferbits(bs.data[j], uint(x)<<lobit, genhimask(lobit))
	// transfer bits to high word
	lon := uintsize - lobit
	bs.data[k] = transferbits(bs.data[k], uint(x)>>lon, genlomask(8-lon))
}

// SetUint16 sets the 16 bits starting at i with the value of x. It panics if
// there are not enough bits.
func (bs *Bitstring) SetUint16(i uint, x uint16) {
	bs.mustExist(i + 15)

	lobit := uint(bitoffset(i))
	j := wordoffset(i)
	k := wordoffset(i + 15)
	if j == k {
		// fast path: same word
		neww := uint(x) << lobit
		mask := genmask(lobit, lobit+16)
		bs.data[j] = transferbits(bs.data[j], neww, mask)
		return
	}
	// transfer bits to low word
	bs.data[j] = transferbits(bs.data[j], uint(x)<<lobit, genhimask(lobit))
	// transfer bits to high word
	lon := uintsize - lobit
	bs.data[k] = transferbits(bs.data[k], uint(x)>>lon, genlomask(16-lon))
}

// SetIntn sets the n bits starting at i with the first n bits of value x.
// It panics if there aren't enough bits in bs or if n is greater than
// the size of a machine word.
func (bs *Bitstring) SetIntn(n, i uint, x uint) { bs.SetIntn(n, i, x) }

// SetInt8 sets the 8 bits starting at i with the value of x. It panics if
// there are not enough bits.
func (bs *Bitstring) SetInt8(i uint, x int8) { bs.SetUint8(i, uint8(x)) }

// SetInt16 sets the 16 bits starting at i with the value of x. It panics if
// there are not enough bits.
func (bs *Bitstring) SetInt16(i uint, x int16) { bs.SetUint16(i, uint16(x)) }

// SetInt32 sets the 32 bits starting at i with the value of x. It panics if
// there are not enough bits.
func (bs *Bitstring) SetInt32(i uint, x int32) { bs.SetUint32(i, uint32(x)) }

// SetInt64 sets the 64 bits starting at i with the value of x. It panics if
// there are not enough bits.
func (bs *Bitstring) SetInt64(i uint, x int64) { bs.SetUint64(i, uint64(x)) }

// prints a string representing the first n bits of the base-2 representatio of x.
func printbits(x, n uint) {
	fmt.Printf(fmt.Sprintf("%%0%db\n", n), x)
}
