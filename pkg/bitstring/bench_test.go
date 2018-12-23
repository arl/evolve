package bitstring

import (
	"fmt"
	"math/rand"
	"testing"
)

var sinka interface{}

func BenchmarkBitstringCopy(b *testing.B) {
	type run struct {
		slen  uint
		human string
	}
	runs := []run{
		{1024, "1k"},
		{100 * 1024, "100k"},
		{10 * 1024 * 1024, "10M"},
	}
	for _, r := range runs {
		b.Run(fmt.Sprintf("BenchBitstringCopy-%v", r.human), func(b *testing.B) {
			var dst *Bitstring

			// create original bitstring
			rng := rand.New(rand.NewSource(99))
			org, err := Random(r.slen, rng)
			if err != nil {
				b.Error("can't create rand bitstring:", err)
			}

			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				// actual benchmark
				dst = Copy(org)
			}
			b.StopTimer()
			sinka = dst
		})
	}
}

var sinkb uint32

func benchmarkUintn(b *testing.B, nbits, i uint) {
	b.ReportAllocs()
	bs, _ := MakeFromString("0000000000000000000000000000000101000000000000000000000000000000")
	var v uint32
	for n := 0; n < b.N; n++ {
		v = bs.Uintn(nbits, i)
	}
	b.StopTimer()
	sinkb = v
}
func benchmarkUint32(b *testing.B, i uint) {
	b.ReportAllocs()
	bs, _ := MakeFromString("0000000000000000000000000000000101000000000000000000000000000000")
	var v uint32
	for n := 0; n < b.N; n++ {
		v = bs.Uint32(i)
	}
	b.StopTimer()
	sinkb = v
}
func benchmarkUint16(b *testing.B, i uint) {
	b.ReportAllocs()
	bs, _ := MakeFromString("0000000000000000000000000000000101000000000000000000000000000000")
	var v uint16
	for n := 0; n < b.N; n++ {
		v = bs.Uint16(i)
	}
	b.StopTimer()
	sinkb = uint32(v)
}
func benchmarkUint8(b *testing.B, i uint) {
	b.ReportAllocs()
	bs, _ := MakeFromString("0000000000000000000000000000000101000000000000000000000000000000")
	var v uint8
	for n := 0; n < b.N; n++ {
		v = bs.Uint8(i)
	}
	b.StopTimer()
	sinkb = uint32(v)
}

func BenchmarkUintnSameWord(b *testing.B)        { benchmarkUintn(b, 32, 32) }
func BenchmarkUintnDifferentWords(b *testing.B)  { benchmarkUintn(b, 32, 31) }
func BenchmarkUint32SameWord(b *testing.B)       { benchmarkUint32(b, 32) }
func BenchmarkUint32DifferentWords(b *testing.B) { benchmarkUint32(b, 31) }
func BenchmarkUint16SameWord(b *testing.B)       { benchmarkUint16(b, 32) }
func BenchmarkUint16DifferentWords(b *testing.B) { benchmarkUint16(b, 31) }
func BenchmarkUint8SameWord(b *testing.B)        { benchmarkUint8(b, 32) }
func BenchmarkUint8DifferentWords(b *testing.B)  { benchmarkUint8(b, 31) }

func Benchmark_genmask(b *testing.B) {
	b.ReportAllocs()

	var v uint32
	for i := 0; i < b.N; i++ {
		v = genmask(4, 27)
	}
	b.StopTimer()
	sinkb = v
}
