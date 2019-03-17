// +build !unchecked

package bitstring

// mustExist panics if i is not a valid bit index for bs, that is if i is
// greater than bs.length.
func (bs *Bitstring) mustExist(i uint) {
	if i >= bs.length {
		panic(ErrIndexOutOfRange)
	}
}
