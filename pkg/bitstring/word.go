package bitstring

// word is a machine word
type word = uint

const (
	uintsize = 32 << (^uint(0) >> 32 & 1) // 32 or 64
	maxuword = ^uint(0)
)
