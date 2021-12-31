package generator

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/arl/evolve"
	"github.com/arl/evolve/pkg/mt19937"
)

func TestGaussian(t *testing.T) {
	const mean, stddev float64 = 147, 17

	rng := rand.New(mt19937.New(99))

	g := NewGaussian(Const(mean), Const(stddev), rng)
	checkGaussianDistribution(t, g, mean, stddev)
}

func TestGaussianDynamic(t *testing.T) {
	const mean, stddev float64 = 147, 17

	rng := rand.New(mt19937.New(99))

	gmean := NewSwappable(Const(mean))
	gstddev := NewSwappable(Const(stddev))
	g := NewGaussian(gmean, gstddev, rng)
	checkGaussianDistribution(t, g, gmean.Next(), gstddev.Next())

	// Change parameters and verify that the generator output is in line
	// with the new distribution.
	gmean.Swap(Const(float64(73)))
	gstddev.Swap(Const(float64(9)))

	checkGaussianDistribution(t, g, gmean.Next(), gstddev.Next())
}

func checkGaussianDistribution(t *testing.T, g Float, wantMean, wantStdDev float64) {
	t.Helper()

	const iterations = 10000
	ds := evolve.NewDataset(iterations)
	for i := 0; i < iterations; i++ {
		ds.AddValue(g.Next())
	}

	const ε = 0.02

	assert.InEpsilon(t, wantMean, ds.ArithmeticMean(), ε, "observed mean is outside of acceptable range")
	// Expected median is the same as expected mean.
	assert.InEpsilon(t, wantMean, ds.Median(), ε, "observed median is outside of acceptable range")
	assert.InEpsilon(t, wantStdDev, ds.SampleStandardDeviation(), ε, "observed standard deviation is outside of acceptable range")
}
