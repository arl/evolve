package generator

import (
	"math"
	"math/rand"
	"testing"

	"github.com/arl/evolve"
	"github.com/arl/evolve/pkg/mt19937"
	"github.com/stretchr/testify/assert"
)

func TestBinomial(t *testing.T) {
	// Check that the observed mean and standard deviation are consistent with
	// the specified distribution parameters.

	const n, p = 20, 0.163

	rng := rand.New(mt19937.New(99))

	g := NewBinomial(ConstInt(n), ConstFloat64(p), rng)
	checkBinomialDistribution(t, g, n, p)
}

func TestBinomialDynamic(t *testing.T) {
	const initn, initp = 20, 0.163

	rng := rand.New(mt19937.New(99))

	ngen := NewAdjustableInt(initn)
	pgen := NewAdjustableFloat(initp)

	g := NewBinomial(ngen, pgen, rng)
	checkBinomialDistribution(t, g, initn, initp)

	// Adjust parameters and ensure that the generator output conforms to this
	// new distribution.
	const adjustn, adjustp = 14, 0.32
	ngen.Set(adjustn)
	pgen.Set(adjustp)

	checkBinomialDistribution(t, g, adjustn, adjustp)
}

func checkBinomialDistribution(t *testing.T, gen Int, n int64, p float64) {
	t.Helper()

	const iterations = 10000

	ds := evolve.NewDataset(iterations)
	for i := 0; i < iterations; i++ {
		val := gen.Next()
		if val < 0 || val > n {
			t.Errorf("generated value out of range, got %v", val)
		}
		ds.AddValue(float64(val))
	}

	const ε = 0.02

	wantMean := float64(n) * p
	wantStdDev := math.Sqrt(float64(n) * p * (1 - p))

	assert.InEpsilon(t, wantMean, ds.ArithmeticMean(), ε, "observed mean is outside of acceptable range")
	assert.InEpsilon(t, wantStdDev, ds.SampleStandardDeviation(), ε, "observed standard deviation is outside of acceptable range")
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
