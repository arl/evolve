// +build 386 arm nacl mips mipsle

package bitstring

import "math"

// word is a machine word
type (
	word  = uint32
	sword = int32

	// export aliases
	Word       = word
	SignedWord = sword
)

const (
	wordlen    = 32
	logwordlen = 5
	minsword   = math.MinInt32
	maxsword   = math.MaxInt32
	maxuword   = math.MaxUint32

	WordLength = wordlen
)
