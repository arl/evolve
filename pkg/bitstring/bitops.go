package bitstring

// creates a mask keep the nth bit of a word.
func bitmask(n uint) uint32 { return 1 << n }

// for a given bit of a bit string, returns the offset of the word
// (from the first word of the bitstring) that contains that bit.
func wordoffset(n uint) uint32 { return uint32(n / wordlen) }

// for a given bit of a bit string, returns the offset of that bit
// from the start of the word that contains it.
func bitoffset(n uint) uint32 { return uint32(n & (wordlen - 1)) }

// returns a mask that keeps the bits in the range [l, h]
func genmask(l, h uint) uint32 { return genlomask(h) & genhimask(l) }

// returns a mask to keep the bits in the range [LSB, n]
func genlomask(n uint) uint32 { return 1<<(n+1) - 1 }

// returns a mask to keep the bits in the range [n, MSB]
func genhimask(n uint) uint32 { return maxuword << n }

// findFirstSetBit returns the offset of the first set bit in w
func findFirstSetBit(w uint32) uint {
	var num uint

	//#if BITS_PER_LONG == 64
	//if ((word & 0xffffffff) == 0) {
	//num += 32;
	//word >>= 32;
	//}
	//#endif
	if (w & 0xffff) == 0 {
		num += 16
		w >>= 16
	}
	if (w & 0xff) == 0 {
		num += 8
		w >>= 8
	}
	if (w & 0xf) == 0 {
		num += 4
		w >>= 4
	}
	if (w & 0x3) == 0 {
		num += 2
		w >>= 2
	}
	if (w & 0x1) == 0 {
		num++
	}
	return num
}
