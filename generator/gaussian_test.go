package generator

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/arl/evolve"
	"github.com/arl/evolve/pkg/mt19937"
)

func TestGaussian(t *testing.T) {
	const mean, stddev = 147, 17

	rng := rand.New(mt19937.New(99))

	gmean := ConstFloat64(mean)
	gstddev := ConstFloat64(stddev)

	g := NewGaussian(&gmean, &gstddev, rng)
	checkGaussianDistribution(t, g, mean, stddev)
}

func TestGaussianDynamic(t *testing.T) {
	const mean, stddev = 147, 17

	rng := rand.New(mt19937.New(99))

	gmean := NewAdjustableFloat(mean)
	gstddev := NewAdjustableFloat(stddev)
	g := NewGaussian(gmean, gstddev, rng)
	checkGaussianDistribution(t, g, gmean.Next(), gstddev.Next())

	// Change parameters and verify that the generator output is in line
	// with the new distribution.
	gmean.Set(73)
	gstddev.Set(9)

	checkGaussianDistribution(t, g, gmean.Next(), gstddev.Next())
}

func checkGaussianDistribution(t *testing.T, gen Float, wantMean, wantStdDev float64) {
	t.Helper()

	const iterations = 10000
	ds := evolve.NewDataset(iterations)
	for i := 0; i < iterations; i++ {
		ds.AddValue(gen.Next())
	}

	const ε = 0.02

	assert.InEpsilon(t, wantMean, ds.ArithmeticMean(), ε, "observed mean is outside of acceptable range")
	// Expected median is the same as expected mean.
	assert.InEpsilon(t, wantMean, ds.Median(), ε, "observed median is outside of acceptable range")
	assert.InEpsilon(t, wantStdDev, ds.SampleStandardDeviation(), ε, "observed standard deviation is outside of acceptable range")
}
