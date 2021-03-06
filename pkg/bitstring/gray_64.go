// +build amd64 arm64 mips64 mips64le ppc64 ppc64le s390x wasm

package bitstring

// Grayn returns the n-bit unsigned integer value represented by the n
// gray-coded bits starting at the bit index i. It panics if there are not
// enough bits or if n is greater than the size of a machine word.
func (bs *Bitstring) Grayn(nbits, i uint) uint {
	v := bs.Uintn(nbits, i)
	v ^= v >> 32
	v ^= v >> 16
	v ^= v >> 8
	v ^= v >> 4
	v ^= v >> 2
	v ^= v >> 1
	return v
}
