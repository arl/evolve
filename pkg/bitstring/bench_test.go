package bitstring

import (
	"math/rand"
	"testing"
)

var sink interface{}

func benchmarkUintn(b *testing.B, nbits, i uint) {
	b.ReportAllocs()
	bs, _ := MakeFromString("0000000000000000000000000000000101000000000000000000000000000000")
	var v word
	for n := 0; n < b.N; n++ {
		v = bs.Uintn(nbits, i)
	}
	b.StopTimer()
	sink = v
}

func benchmarkUint32(b *testing.B, i uint) {
	b.ReportAllocs()
	bs, _ := MakeFromString("0000000000000000000000000000000101000000000000000000000000000000")
	var v uint32
	for n := 0; n < b.N; n++ {
		v = bs.Uint32(i)
	}
	b.StopTimer()
	sink = v
}

func benchmarkUint16(b *testing.B, i uint) {
	b.ReportAllocs()
	bs, _ := MakeFromString("0000000000000000000000000000000101000000000000000000000000000000")
	var v uint16
	for n := 0; n < b.N; n++ {
		v = bs.Uint16(i)
	}
	b.StopTimer()
	sink = v
}

func benchmarkUint8(b *testing.B, i uint) {
	b.ReportAllocs()
	bs, _ := MakeFromString("0000000000000000000000000000000101000000000000000000000000000000")
	var v uint8
	for n := 0; n < b.N; n++ {
		v = bs.Uint8(i)
	}
	b.StopTimer()
	sink = v
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

	var v word
	for i := 0; i < b.N; i++ {
		v = genmask(4, 27)
	}
	b.StopTimer()
	sink = v
}

func Benchmark_lomask(b *testing.B) {
	b.ReportAllocs()

	var v word
	for i := 0; i < b.N; i++ {
		v = genlomask(27)
	}
	b.StopTimer()
	sink = v
}

func BenchmarkSwapRange(b *testing.B) {
	x := New(1026)
	y := New(1026)

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		SwapRange(x, y, 1, 1024)
	}
}

func BenchmarkRandom(b *testing.B) {
	var x *Bitstring

	rng := rand.New(rand.NewSource(99))
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		x = Random(1026, rng)
	}
	b.StopTimer()
	sink = x
}

func BenchmarkEquals(b *testing.B) {
	rng := rand.New(rand.NewSource(99))
	x := Random(1026, rng)
	y := Copy(x)
	var res bool
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		res = x.Equals(y)
	}
	b.StopTimer()
	sink = res
}
