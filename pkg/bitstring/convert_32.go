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
	return uint32(loword | hiword<<(wordlen-off))
}
