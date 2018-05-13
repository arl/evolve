package bitstring

import (
	"fmt"
	"math/rand"
	"testing"
)

var sink interface{}

func BenchmarkBitstringCopy(b *testing.B) {
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
		name := fmt.Sprintf("BenchBitstringCopy-%v", r.human)
		b.Run(name, func(b *testing.B) {
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
			sink = dst
		})
	}
}
