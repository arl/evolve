package xover

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestByteSliceMater(t *testing.T) {
	rng := rand.New(rand.NewSource(99))

	xover := New(ByteSliceMater{})
	pop := make([]interface{}, 4)
	pop[0] = []byte{1, 2, 3, 4, 5}
	pop[1] = []byte{6, 7, 8, 9, 10}
	pop[2] = []byte{11, 12, 13, 14, 15}
	pop[3] = []byte{16, 17, 18, 19, 20}

	for i := 0; i < 20; i++ {
		values := make(map[byte]struct{}, 20)
		pop = xover.Apply(pop, rng)
		if len(pop) != 4 {
			t.Error("population size changed, want 4, got", len(pop))
		}

		for _, individual := range pop {
			s := individual.([]byte)
			if len(s) != 5 {
				t.Error("wrong candidate length, want 5, got", len(s))
			}
			for _, value := range s {
				values[value] = struct{}{}
			}
		}
		// All of the individual elements should still be present, just jumbled
		// up between individuals.
		if len(values) != 20 {
			t.Error("wrong number of candidates, want 20, got", len(values))
		}
	}
}

// ByteSliceMater must operate on []byte of equal length. It should panic if
// different length slices are used.
func TestByteSliceMaterDifferentLength(t *testing.T) {
	rng := rand.New(rand.NewSource(99))

	xover := New(ByteSliceMater{})
	pop := make([]interface{}, 2)
	pop[0] = []byte{1, 2, 3, 4, 5}
	pop[1] = []byte{2, 4, 8, 10, 12, 14, 16}

	assert.Panics(t, func() { xover.Apply(pop, rng) })
}

var sink interface{}

func createRandSlice(l int) ([]byte, error) {
	s := make([]byte, l)
	if _, err := rand.Read(s); err != nil {
		return nil, err
	}
	return s, nil
}

func BenchmarkByteSliceAppend(b *testing.B) {
	type run struct {
		slen  int
		human string
	}
	runs := []run{
		{1024, "1k"},
		{100 * 1024, "100k"},
		{10 * 1024 * 1024, "10M"},
	}
	for _, r := range runs {
		name := fmt.Sprintf("BenchByteSliceAppend-%v", r.human)
		b.Run(name, func(b *testing.B) {
			var dst []byte

			// allocate original slice
			org, err := createRandSlice(r.slen)
			if err != nil {
				b.Error("can't create rand slice:", err)
			}

			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				// actual benchmark
				dst = append([]byte{}, org...)
			}
			b.StopTimer()
			sink = dst
		})
	}
}

func BenchmarkByteSliceCopy(b *testing.B) {
	type run struct {
		slen  int
		human string
	}
	runs := []run{
		{1024, "1k"},
		{100 * 1024, "100k"},
		{10 * 1024 * 1024, "10M"},
	}
	for _, r := range runs {
		name := fmt.Sprintf("BenchByteSliceCopy-%v", r.human)
		b.Run(name, func(b *testing.B) {
			var dst []byte

			// allocate original slice
			org, err := createRandSlice(r.slen)
			if err != nil {
				b.Error("can't create rand slice:", err)
			}

			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				// actual benchmark
				dst = make([]byte, len(org))
				copy(dst, org)
			}
			b.StopTimer()
			sink = dst
		})
	}
}
