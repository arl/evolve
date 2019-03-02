// +build amd64 arm64 mips64 mips64le ppc64 ppc64le s390x wasm

package bitstring

// Uint32 returns the uint32 value represented by the 32 bits starting at the
// given bit. It panics if there are not enough bits.
func (bs *Bitstring) Uint32(i uint) uint32 {
	bs.mustExist(i + 31)

	off := bitoffset(i)
	loword := bs.data[wordoffset(i)] >> off
	hiword := bs.data[wordoffset(i+31)] & ((1 << off) - 1)
	return uint32(loword | hiword<<(wordlen-off))
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
	return uint64(loword | hiword<<(wordlen-off))

}
