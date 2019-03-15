package bitstring

const (
	uintsize = 32 << (^uint(0) >> 32 & 1) // 32 or 64
	maxuint  = ^uint(0)
)
