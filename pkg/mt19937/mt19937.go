// Adapted from an implementation of the 64bit Mersenne Twister PRNG
// original copyright:
// Copyright (C) 2013  Jochen Voss <voss@seehuhn.de>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

// Package mt19937 implements a b4-bit Mersenne Twister source of random numbers
// satisfying the rand.Source64 interface.
//
// A single mt19937.MT19937 instance is not safe for concurrent access by different
// goroutines.
//
// For random numbers suitable for security-sensitive work, see the crypto/rand
// package.
package mt19937

import (
	"encoding/binary"
	"hash/maphash"
)

const (
	n = 312
	m = 156

	himask uint64 = 0xffffffff80000000
	lomask uint64 = 0x000000007fffffff

	matrixa uint64 = 0xB5026F5AA96619E9
)

// MT19937 holds the state of a 64-bit Mersenne Twister PRNG.
//
// Use New to create a new, randomly seeded MT19937 instance. The zero-value of
// MT19937 is not valid instance.
//
// This struct is not safe for concurrent access by different goroutines. If
// more than one goroutine accesses the PRNG, callers must synchronise access
// using sync.Mutex or similar.
type MT19937 struct {
	state [n]uint64
	index int
}

// New creates and initializes a 64bit Mersenne Twister PRNG, seeded with a
// source of randomness.
func New() *MT19937 {
	var mt MT19937
	mt.Seed(int64(new(maphash.Hash).Sum64()))
	return &mt
}

// Seed uses the provided seed value to initialize the generator to a
// deterministic state. Seed should not be called concurrently with any other
// MT19937 method.
func (mt *MT19937) Seed(seed int64) {
	mt.state[0] = uint64(seed)
	for i := uint64(1); i < n; i++ {
		mt.state[i] = 6364136223846793005*(mt.state[i-1]^(mt.state[i-1]>>62)) + i
	}
	mt.index = n
}

// SeedFromSlice uses the given slice of 64bit values to set the
// generator state.
func (mt *MT19937) SeedFromSlice(key []uint64) {
	mt.Seed(19650218)

	i := uint64(1)
	j := 0
	k := len(key)
	if n > k {
		k = n
	}
	for k > 0 {
		mt.state[i] = (mt.state[i] ^ ((mt.state[i-1] ^ (mt.state[i-1] >> 62)) * 3935559000370003845) +
			key[j] + uint64(j))
		i++
		if i >= n {
			mt.state[0] = mt.state[n-1]
			i = 1
		}
		j++
		if j >= len(key) {
			j = 0
		}
		k--
	}
	for j := uint64(0); j < n-1; j++ {
		mt.state[i] = mt.state[i] ^ ((mt.state[i-1] ^ (mt.state[i-1] >> 62)) * 2862933555777941757) - i
		i++
		if i >= n {
			mt.state[0] = mt.state[n-1]
			i = 1
		}
	}
	mt.state[0] = 1 << 63
}

// Uint64 returns a pseudo-random 64-bit value as a uint64.
func (mt *MT19937) Uint64() uint64 {
	if mt.index >= n {
		for i := 0; i < n-m; i++ {
			y := (mt.state[i] & himask) | (mt.state[i+1] & lomask)
			mt.state[i] = mt.state[i+m] ^ (y >> 1) ^ ((y & 1) * matrixa)
		}
		for i := n - m; i < n-1; i++ {
			y := (mt.state[i] & himask) | (mt.state[i+1] & lomask)
			mt.state[i] = mt.state[i+(m-n)] ^ (y >> 1) ^ ((y & 1) * matrixa)
		}
		y := (mt.state[n-1] & himask) | (mt.state[0] & lomask)
		mt.state[n-1] = mt.state[m-1] ^ (y >> 1) ^ ((y & 1) * matrixa)
		mt.index = 0
	}
	y := mt.state[mt.index]
	y ^= (y >> 29) & 0x5555555555555555
	y ^= (y << 17) & 0x71D67FFFEDA60000
	y ^= (y << 37) & 0xFFF7EEE000000000
	y ^= (y >> 43)
	mt.index++
	return y
}

// Int63 returns a non-negative pseudo-random 63-bit integer as an int64.
func (mt *MT19937) Int63() int64 {
	if mt.index >= n {
		for i := 0; i < n-m; i++ {
			y := (mt.state[i] & himask) | (mt.state[i+1] & lomask)
			mt.state[i] = mt.state[i+m] ^ (y >> 1) ^ ((y & 1) * matrixa)
		}
		for i := n - m; i < n-1; i++ {
			y := (mt.state[i] & himask) | (mt.state[i+1] & lomask)
			mt.state[i] = mt.state[i+(m-n)] ^ (y >> 1) ^ ((y & 1) * matrixa)
		}
		y := (mt.state[n-1] & himask) | (mt.state[0] & lomask)
		mt.state[n-1] = mt.state[m-1] ^ (y >> 1) ^ ((y & 1) * matrixa)
		mt.index = 0
	}
	y := mt.state[mt.index]
	y ^= (y >> 29) & 0x5555555555555555
	y ^= (y << 17) & 0x71D67FFFEDA60000
	y ^= (y << 37) & 0xFFF7EEE000000000
	y ^= (y >> 43)
	mt.index++
	return int64(y & 0x7fffffffffffffff)
}

// Read generates len(p) random bytes and writes them into p. It
// always returns len(p) and a nil error.
// Read should not be called concurrently with any other MT19937 method.
func (mt *MT19937) Read(p []byte) (n int, err error) {
	for n+8 <= len(p) {
		binary.LittleEndian.PutUint64(p[n:], mt.Uint64())
		n += 8
	}
	if n < len(p) {
		for n < len(p) {
			ui64 := mt.Uint64()
			p[n] = byte(ui64)
			ui64 >>= 8
			n++
		}
	}
	return n, nil
}
