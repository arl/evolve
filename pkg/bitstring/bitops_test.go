package bitstring

import (
	"strconv"
	"testing"
)

func bintoa(s string) uint32 {
	i, err := strconv.ParseInt(s, 2, 64)
	if err != nil {
		panic(err)
	}
	return uint32(i)
}

func Test_genlomask(t *testing.T) {
	tests := []struct {
		n    uint
		want uint32
	}{
		{n: 0, want: bintoa("00000000000000000000000000000001")},
		{n: 1, want: bintoa("00000000000000000000000000000011")},
		{n: wordlen - 2, want: bintoa("01111111111111111111111111111111")},
		{n: wordlen - 1, want: bintoa("11111111111111111111111111111111")},
		{n: wordlen, want: bintoa("11111111111111111111111111111111")},
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
		want uint32
	}{
		{n: 0, want: bintoa("11111111111111111111111111111111")},
		{n: 1, want: bintoa("11111111111111111111111111111110")},
		{n: wordlen - 2, want: bintoa("11000000000000000000000000000000")},
		{n: wordlen - 1, want: bintoa("10000000000000000000000000000000")},
		{n: wordlen, want: bintoa("00000000000000000000000000000000")},
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
		want uint32
	}{
		{l: 0, h: 0, want: bintoa("00000000000000000000000000000001")},
		{l: 0, h: 1, want: bintoa("00000000000000000000000000000011")},
		{l: 0, h: 2, want: bintoa("00000000000000000000000000000111")},
		{l: 1, h: 1, want: bintoa("00000000000000000000000000000010")},
		{l: 1, h: 2, want: bintoa("00000000000000000000000000000110")},
		{l: 0, h: 31, want: bintoa("11111111111111111111111111111111")},
		{l: 1, h: 31, want: bintoa("11111111111111111111111111111110")},
		{l: 0, h: 30, want: bintoa("01111111111111111111111111111111")},
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
		w    uint32
		want uint
	}{
		{w: uint32(bintoa("00000000000000000000000000000001")), want: 0},
		{w: uint32(bintoa("00000000000000000000000000000010")), want: 1},
		{w: uint32(bintoa("10000000000000000000000000000001")), want: 0},
		{w: uint32(bintoa("00000000000001111111000000000100")), want: 2},
		{w: uint32(bintoa("00000000000001111111000000000000")), want: 12},
		{w: uint32(bintoa("10000000000000000000000000000000")), want: 31},
		{w: uint32(bintoa("00000000000000000000000000000000")), want: 31},
		{w: uint32(bintoa("11111111111111111111111111111111")), want: 0},
	}
	for _, tt := range tests {
		if got := findFirstSetBit(tt.w); got != tt.want {
			t.Errorf("findFirstSetBit(%s) got %d, want %d",
				sprintubits(tt.w, 32), got, tt.want)
		}
	}
}
