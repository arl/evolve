package generator

import (
	"math/rand"
	"testing"

	"github.com/arl/evolve/pkg/mt19937"

	"github.com/stretchr/testify/assert"
)

func TestBinomial(t *testing.T) {
	// Check that the observed mean and standard deviation are consistent with
	// the specified distribution parameters.
	tests := []struct {
		n uint32
		p float64
	}{
		{n: 20, p: 0.163},
		{n: 2000, p: 0.163},
		{n: 20000, p: 0.163},
		{n: 20000, p: 0.1},
		{n: 20000, p: 0.9},
	}

	rng := rand.New(mt19937.New(99))
	for _, tt := range tests {
		g := NewBinomial[uint32](Const(uint32(tt.n)), Const(tt.p), rng)
		checkBinomialDistribution[uint32](t, g, tt.n, tt.p)
	}
}

func TestBinomialDynamic(t *testing.T) {
	const initn, initp = 20, 0.163

	rng := rand.New(mt19937.New(99))

	ngen := NewSwappable(Const(uint64(initn)))
	pgen := NewSwappable(Const(initp))

	g := NewBinomial[uint64](ngen, pgen, rng)
	checkBinomialDistribution[uint64](t, g, initn, initp)

	// Adjust parameters and ensure that the generator output conforms to this
	// new distribution.
	const adjustn, adjustp = 14, 0.32
	ngen.Swap(Const(uint64(adjustn)))
	pgen.Swap(Const(adjustp))

	checkBinomialDistribution[uint64](t, g, adjustn, adjustp)
}

func Test_floatToFixedBits(t *testing.T) {
	t.Run("0", func(t *testing.T) {
		bits := floatToFixedBits(0.6875)
		assert.Equal(t, "1011", bits.String(), "binary representation should be 1011")
	})

	// Ensure floatToFixedBits correctly deals with 0.
	t.Run("0", func(t *testing.T) {
		bits := floatToFixedBits(0)
		assert.Zero(t, bits.OnesCount(), "binary representation should be 0")
	})

	// floatToFixedBits creates fixed point representations without sign bit, so
	// trying to convert a negative value should panic.
	t.Run("panic when negative", func(t *testing.T) {
		assert.Panics(t, func() { floatToFixedBits(-0.5) })
	})

	// floatToFixedBits should panic when given number >= 1 since the
	// representation doesn't allow such numbers.
	t.Run("panic when >=1", func(t *testing.T) {
		assert.Panics(t, func() { floatToFixedBits(1) })
	})
}

func benchmarkBinomialEvenProbability(n int) func(b *testing.B) {
	return func(b *testing.B) {
		rng := rand.New(mt19937.New(99))
		g := NewBinomial[uint32](Const(uint32(n)), Const(0.163), rng)
		for i := 0; i < b.N; i++ {
			g.evenProbability(uint32(n))
		}
	}
}

func Benchmark_evenProbability(b *testing.B) {
	b.Run("10", benchmarkBinomialEvenProbability(10))
	b.Run("100", benchmarkBinomialEvenProbability(100))
	b.Run("1000", benchmarkBinomialEvenProbability(1000))
}
