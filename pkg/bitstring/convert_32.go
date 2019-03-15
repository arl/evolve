// +build 386 arm nacl mips mipsle

package bitstring

// Uint32 returns the uint32 value represented by the 32 bits starting at the
// given bit. It panics if there are not enough bits.
func (bs *Bitstring) Uint32(i uint) uint32 {
	bs.mustExist(i + 31)

	// fast path: i is a multiple of 32
	if i&((1<<5)-1) == 0 {
		return uint32(bs.data[i>>5])
	}

	w := wordoffset(i)
	off := bitoffset(i)
	loword := bs.data[w] >> off
	hiword := bs.data[w+1] & ((1 << off) - 1)
	return uint32(loword | hiword<<(uintsize-off))
}

// Uint64 returns the uint64 value represented by the 64 bits starting at the
// given bit. It panics if there are not enough bits.
func (bs *Bitstring) Uint64(i uint) uint64 {
	lohw := bs.Uint32(i)
	hihw := bs.Uint32(i + 32)
	return uint64(hihw)<<32 | uint64(lohw)
}

// SetUint32 sets the 32 bits starting at i with the value of x. It panics if
// there are not enough bits.
func (bs *Bitstring) SetUint32(i uint, x uint32) {
	bs.mustExist(i + 31)

	lobit := uint(bitoffset(i))
	j := wordoffset(i)

	// fast path: i is a multiple of 32
	if i&((1<<5)-1) == 0 {
		bs.data[i>>5] = x
		return
	}

	k := wordoffset(i + 31)
	if j == k {
		// fast path: same word
		neww := uint(x) << lobit
		mask := genmask(lobit, lobit+32)
		bs.data[j] = transferbits(bs.data[j], neww, mask)
		return
	}
	// transfer bits to low word
	bs.data[j] = transferbits(bs.data[j], uint(x)<<lobit, genhimask(lobit))
	// transfer bits to high word
	lon := uintsize - lobit
	bs.data[k] = transferbits(bs.data[k], uint(x)>>lon, genlomask(32-lon))
}

// SetUint64 sets the 64 bits starting at i with the value of x. It panics if
// there are not enough bits.
func (bs *Bitstring) SetUint64(i uint, x uint64) {
	bs.SetUint32(i, uint32(x))
	bs.SetUint32(i+32, uint32(x>>32))
}
