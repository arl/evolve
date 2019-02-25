package bitstring

// creates a mask keep the nth bit of a word.
func bitmask(n uint) word { return 1 << n }

// for a given bit of a bit string, returns the offset of the word
// (from the first word of the bitstring) that contains that bit.
func wordoffset(n uint) word { return word(n / wordlen) }

// for a given bit of a bit string, returns the offset of that bit
// from the start of the word that contains it.
func bitoffset(n uint) word { return word(n & (wordlen - 1)) }

// returns a mask that keeps the bits in the range [l, h[
// behaviour undefined if any argument is greater than wordlen.
func genmask(l, h uint) word { return genlomask(h) & genhimask(l) }

// returns a mask to keep the n LSB (least significant bits).
// Undefined behaviour if n is greater than wordlen.
func genlomask(n uint) word { return maxuword >> (wordlen - n) }

// returns a mask to keep the n MSB (most significan bits).
// Undefined behaviour if n is greater than wordlen.
func genhimask(n uint) word { return maxuword << n }

// findFirstSetBit returns the offset of the first set bit in w
func findFirstSetBit(w word) uint {
	var num uint

	if wordlen == 64 {
		if (w & 0xffffffff) == 0 {
			num += 32
			w >>= 32
		}
	}
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
