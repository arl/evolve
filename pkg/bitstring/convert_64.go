// +build amd64 arm64 mips64 mips64le ppc64 ppc64le s390x wasm

package bitstring

// Uint32 returns the uint32 value represented by the 32 bits starting at the
// given bit. It panics if there are not enough bits.
func (bs *Bitstring) Uint32(i uint) uint32 {
	bs.mustExist(i + 31)

	off := bitoffset(i)
	loword := bs.data[wordoffset(i)] >> off
	hiword := bs.data[wordoffset(i+31)] & ((1 << off) - 1)
	return uint32(loword | hiword<<(uintsize-off))
}

// Uint64 returns the uint64 value represented by the 64 bits starting at the
// given bit. It panics if there are not enough bits.
func (bs *Bitstring) Uint64(i uint) uint64 {
	bs.mustExist(i + 63)

	// fast path: i is a multiple of 64
	if i&((1<<6)-1) == 0 {
		return uint64(bs.data[i>>6])
	}

	w := wordoffset(i)
	off := bitoffset(i)
	loword := bs.data[w] >> off
	hiword := bs.data[w+1] & ((1 << off) - 1)
	return uint64(loword | hiword<<(uintsize-off))
}

// SetUint32 sets the 32 bits starting at i with the value of x. It panics if
// there are not enough bits.
func (bs *Bitstring) SetUint32(i uint, x uint32) {
	bs.mustExist(i + 31)

	lobit := uint(bitoffset(i))
	j := wordoffset(i)
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
	bs.mustExist(i + 63)

	lobit := uint(bitoffset(i))
	j := wordoffset(i)

	// fast path: i is a multiple of 64
	if i&((1<<6)-1) == 0 {
		bs.data[i>>6] = uint(x)
		return
	}

	k := wordoffset(i + 63)
	if j == k {
		// fast path: same word
		neww := uint(x) << lobit
		mask := genmask(lobit, lobit+64)
		bs.data[j] = transferbits(bs.data[j], neww, mask)
		return
	}
	// transfer bits to low word
	bs.data[j] = transferbits(bs.data[j], uint(x)<<lobit, genhimask(lobit))
	// transfer bits to high word
	lon := (uintsize - lobit)
	bs.data[k] = transferbits(bs.data[k], uint(x)>>lon, genlomask(64-lon))
}
