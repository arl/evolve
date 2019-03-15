package bitstring

import (
	"fmt"
	"strconv"
	"testing"
)

func atobin(s string) uint {
	i, err := strconv.ParseUint(s, 2, uintsize)
	if err != nil {
		panic(fmt.Sprintf("Can't convert %s to base 2: %s", s, err))
	}
	return uint(i)
}

func Test_genlomask(t *testing.T) {
	tests := []struct {
		n    uint
		want uint
	}{
		{n: 0, want: atobin("00000000000000000000000000000000")},
		{n: 1, want: atobin("00000000000000000000000000000001")},
		{n: 2, want: atobin("00000000000000000000000000000011")},
		{n: uintsize - 2, want: maxuint >> 2},
		{n: uintsize - 1, want: maxuint >> 1},
		{n: uintsize, want: maxuint},
	}
	for _, tt := range tests {
		if got := genlomask(tt.n); got != tt.want {
			t.Errorf("genlomask(%d) got %s, want %s", tt.n,
				sprintubits(got, 32), sprintubits(tt.want, 32))
		}
	}
}

func Test_genhimask(t *testing.T) {
	tests := []struct {
		n    uint
		want uint
	}{
		{n: 0, want: maxuint},
		{n: 1, want: maxuint - 1},
		{n: uintsize - 2, want: 1<<(uintsize-1) + 1<<(uintsize-2)},
		{n: uintsize - 1, want: 1 << (uintsize - 1)},
	}
	for _, tt := range tests {
		if got := genhimask(tt.n); got != tt.want {
			t.Errorf("genhimask(%d) got %s, want %s", tt.n,
				sprintubits(got, 32), sprintubits(tt.want, 32))
		}
	}
}

func Test_genmask(t *testing.T) {
	tests := []struct {
		l, h uint
		want uint
	}{
		{l: 0, h: 0, want: atobin("00000000000000000000000000000000")},
		{l: 0, h: 1, want: atobin("00000000000000000000000000000001")},
		{l: 0, h: 2, want: atobin("00000000000000000000000000000011")},
		{l: 1, h: 1, want: atobin("00000000000000000000000000000000")},
		{l: 1, h: 2, want: atobin("00000000000000000000000000000010")},
		{l: 0, h: 31, want: atobin("01111111111111111111111111111111")},
		{l: 1, h: 31, want: atobin("01111111111111111111111111111110")},
		{l: 0, h: 30, want: atobin("00111111111111111111111111111111")},
	}
	for _, tt := range tests {
		if got := genmask(tt.l, tt.h); got != tt.want {
			t.Errorf("genmask(%d, %d) got %s, want %s", tt.l, tt.h,
				sprintubits(got, 32), sprintubits(tt.want, 32))
		}
	}
}

func Test_findFirstSetBit(t *testing.T) {
	tests := []struct {
		w    uint
		want uint
	}{
		{w: atobin("00000000000000000000000000000001"), want: 0},
		{w: atobin("00000000000000000000000000000010"), want: 1},
		{w: atobin("10000000000000000000000000000001"), want: 0},
		{w: atobin("00000000000001111111000000000100"), want: 2},
		{w: atobin("00000000000001111111000000000000"), want: 12},
		{w: atobin("10000000000000000000000000000000"), want: 31},
		{w: atobin("00000000000000000000000000000000"), want: uintsize - 1},
		{w: atobin("11111111111111111111111111111111"), want: 0},
	}
	for _, tt := range tests {
		if got := findFirstSetBit(tt.w); got != tt.want {
			t.Errorf("findFirstSetBit(%s) got %d, want %d",
				sprintubits(tt.w, 32), got, tt.want)
		}
	}
}
