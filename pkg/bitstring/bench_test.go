package bitstring

import (
	"math/rand"
	"testing"
)

var sink interface{}

func benchmarkUintn(b *testing.B, nbits, i uint) {
	b.ReportAllocs()
	bs, _ := MakeFromString("0000000000000000000000000000000101000000000000000000000000000000")
	var v uint
	for n := 0; n < b.N; n++ {
		v = bs.Uintn(nbits, i)
	}
	b.StopTimer()
	sink = v
}

func benchmarkUint64(b *testing.B, i uint) {
	b.ReportAllocs()
	bs, _ := MakeFromString("00000000000000000000000000000001010000000000000000000000000000000")
	var v uint64
	for n := 0; n < b.N; n++ {
		v = bs.Uint64(i)
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
func BenchmarkUint64SameWord(b *testing.B)       { benchmarkUint64(b, 0) }
func BenchmarkUint64DifferentWords(b *testing.B) { benchmarkUint64(b, 1) }
func BenchmarkUint32SameWord(b *testing.B)       { benchmarkUint32(b, 32) }
func BenchmarkUint32DifferentWords(b *testing.B) { benchmarkUint32(b, 31) }
func BenchmarkUint16SameWord(b *testing.B)       { benchmarkUint16(b, 32) }
func BenchmarkUint16DifferentWords(b *testing.B) { benchmarkUint16(b, 31) }
func BenchmarkUint8SameWord(b *testing.B)        { benchmarkUint8(b, 32) }
func BenchmarkUint8DifferentWords(b *testing.B)  { benchmarkUint8(b, 31) }

func Benchmark_mask(b *testing.B) {
	b.ReportAllocs()

	var v uint
	for i := 0; i < b.N; i++ {
		v = mask(4, 27)
	}
	b.StopTimer()
	sink = v
}

func Benchmark_lomask(b *testing.B) {
	b.ReportAllocs()

	var v uint
	for i := 0; i < b.N; i++ {
		v = lomask(27)
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

func BenchmarkSetUint8(b *testing.B) {
	bs := New(67)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		bs.SetUint8(59, 255)
	}
	b.StopTimer()
	sink = bs
}

func BenchmarkSetUintn(b *testing.B) {
	bs := New(117)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		bs.SetUintn(64, 35, 0x9cfbeb71ee3fcf5f)
	}
	b.StopTimer()
	sink = bs
}
