// +build amd64 arm64 mips64 mips64le ppc64 ppc64le s390x wasm

package bitstring

import "math"

// word is a machine word
type (
	word  = uint64
	sword = int64

	// export aliases
	Word       = word
	SignedWord = sword
)

const (
	wordlen    = 64
	logwordlen = 6
	minsword   = math.MinInt64
	maxsword   = math.MaxInt64
	maxuword   = math.MaxUint64

	WordLength = wordlen
)
