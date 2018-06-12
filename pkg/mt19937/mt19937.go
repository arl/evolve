// mt19937.go - an implementation of the 64bit Mersenne Twister PRNG
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

// Package mt19937 implements a Mersenne Twister source of random numbers
// satisfying the rand.Source64 interface.
//
// Structs in mt19937 package are not safe for concurrent access by different
// goroutines. If more than one goroutine accesses the PRNG, the callers must
// synchronise access using sync.Mutex or similar.
//
// For random numbers suitable for security-sensitive work, see the crypto/rand
// package.
package mt19937

const (
	n = 312
	m = 156

	himask uint64 = 0xffffffff80000000
	lomask uint64 = 0x000000007fffffff

	matrixa uint64 = 0xB5026F5AA96619E9
)

// MT19937 holds the state of one instance of the Mersenne Twister PRNG.
// New instances can be allocated using the mt19937.New() function. MT19937
// implements the rand.Source64 interface and rand.New() from the math/rand
// package can be used to generate different distributions from a MT19937 PRNG.
//
// This struct is not safe for concurrent access by different goroutines. If
// more than one goroutine accesses the PRNG, the callers must synchronise
// access using sync.Mutex or similar.
type MT19937 struct {
	state []uint64
	index int
}

// TODO: provide another version of the constructor that is safe for concurrent
// use by multiple goroutines

// New returns a new instance of the 64bit Mersenne Twister with the specified
// seed.
func New(seed int64) *MT19937 {
	res := MT19937{state: make([]uint64, n)}
	res.Seed(seed)
	return &res
}

// Seed uses the given 64bit value to initialise the generator state.
// This method is part of the rand.Source interface.
func (mt *MT19937) Seed(seed int64) {
	x := mt.state
	x[0] = uint64(seed)
	for i := uint64(1); i < n; i++ {
		x[i] = 6364136223846793005*(x[i-1]^(x[i-1]>>62)) + i
	}
	mt.index = n
}

// SeedFromSlice uses the given slice of 64bit values to set the
// generator state.
func (mt *MT19937) SeedFromSlice(key []uint64) {
	mt.Seed(19650218)

	x := mt.state
	i := uint64(1)
	j := 0
	k := len(key)
	if n > k {
		k = n
	}
	for k > 0 {
		x[i] = (x[i] ^ ((x[i-1] ^ (x[i-1] >> 62)) * 3935559000370003845) +
			key[j] + uint64(j))
		i++
		if i >= n {
			x[0] = x[n-1]
			i = 1
		}
		j++
		if j >= len(key) {
			j = 0
		}
		k--
	}
	for j := uint64(0); j < n-1; j++ {
		x[i] = x[i] ^ ((x[i-1] ^ (x[i-1] >> 62)) * 2862933555777941757) - i
		i++
		if i >= n {
			x[0] = x[n-1]
			i = 1
		}
	}
	x[0] = 1 << 63
}

// Uint64 generates a (pseudo-)random 64bit value. The output can be
// used as a replacement for a sequence of independent, uniformly
// distributed samples in the range 0, 1, ..., 2^64-1.
func (mt *MT19937) Uint64() uint64 {
	x := mt.state
	if mt.index >= n {
		for i := 0; i < n-m; i++ {
			y := (x[i] & himask) | (x[i+1] & lomask)
			x[i] = x[i+m] ^ (y >> 1) ^ ((y & 1) * matrixa)
		}
		for i := n - m; i < n-1; i++ {
			y := (x[i] & himask) | (x[i+1] & lomask)
			x[i] = x[i+(m-n)] ^ (y >> 1) ^ ((y & 1) * matrixa)
		}
		y := (x[n-1] & himask) | (x[0] & lomask)
		x[n-1] = x[m-1] ^ (y >> 1) ^ ((y & 1) * matrixa)
		mt.index = 0
	}
	y := x[mt.index]
	y ^= (y >> 29) & 0x5555555555555555
	y ^= (y << 17) & 0x71D67FFFEDA60000
	y ^= (y << 37) & 0xFFF7EEE000000000
	y ^= (y >> 43)
	mt.index++
	return y
}

// Int63 generates a (pseudo-)random 63bit value.  The output can be
// used as a replacement for a sequence of independent, uniformly
// distributed samples in the range 0, 1, ..., 2^63-1.  This method is
// part of the rand.Source interface.
func (mt *MT19937) Int63() int64 {
	x := mt.state
	if mt.index >= n {
		for i := 0; i < n-m; i++ {
			y := (x[i] & himask) | (x[i+1] & lomask)
			x[i] = x[i+m] ^ (y >> 1) ^ ((y & 1) * matrixa)
		}
		for i := n - m; i < n-1; i++ {
			y := (x[i] & himask) | (x[i+1] & lomask)
			x[i] = x[i+(m-n)] ^ (y >> 1) ^ ((y & 1) * matrixa)
		}
		y := (x[n-1] & himask) | (x[0] & lomask)
		x[n-1] = x[m-1] ^ (y >> 1) ^ ((y & 1) * matrixa)
		mt.index = 0
	}
	y := x[mt.index]
	y ^= (y >> 29) & 0x5555555555555555
	y ^= (y << 17) & 0x71D67FFFEDA60000
	y ^= (y << 37) & 0xFFF7EEE000000000
	y ^= (y >> 43)
	mt.index++
	return int64(y & 0x7fffffffffffffff)
}

// Read fills `p` with (pseudo-)random bytes.  This method implements
// the io.Reader interface. The returned length `n` always equals
// `len(p)` and `err` is always nil.
func (mt *MT19937) Read(p []byte) (n int, err error) {
	for n+8 <= len(p) {
		ui64 := mt.Uint64()
		p[n] = byte(ui64)
		p[n+1] = byte(ui64 >> 8)
		p[n+2] = byte(ui64 >> 16)
		p[n+3] = byte(ui64 >> 24)
		p[n+4] = byte(ui64 >> 32)
		p[n+5] = byte(ui64 >> 40)
		p[n+6] = byte(ui64 >> 48)
		p[n+7] = byte(ui64 >> 56)
		n += 8
	}
	if n < len(p) {
		ui64 := mt.Uint64()
		for n < len(p) {
			p[n] = byte(ui64)
			ui64 >>= 8
			n++
		}
	}
	return n, nil
}
