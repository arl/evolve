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
