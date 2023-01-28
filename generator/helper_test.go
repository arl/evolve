package generator

import (
	"math"
	"testing"

	"github.com/arl/evolve/pkg/dataset"
	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/constraints"
)

func checkGaussianDistribution(t *testing.T, g Float, wantMean, wantStdDev float64) {
	t.Helper()

	const iterations = 50000
	ds := dataset.New(iterations)
	for i := 0; i < iterations; i++ {
		ds.AddValue(g.Next())
	}

	const ε = 0.02

	assert.InEpsilon(t, wantMean, ds.ArithmeticMean(), ε, "observed mean is outside of acceptable range")
	// Expected median is the same as expected mean.
	assert.InEpsilon(t, wantMean, ds.Median(), ε, "observed median is outside of acceptable range")
	assert.InEpsilon(t, wantStdDev, ds.SampleStandardDeviation(), ε, "observed standard deviation is outside of acceptable range")
}

func checkBinomialDistribution[T constraints.Integer | constraints.Float](t *testing.T, g Generator[T], n T, p float64) {
	t.Helper()

	const iterations = 10000

	ds := dataset.New(iterations)
	for i := 0; i < iterations; i++ {
		val := g.Next()
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

func checkExponentialDistribution(t *testing.T, g Float, rate float64) {
	t.Helper()

	const iterations = 50000
	ds := dataset.New(iterations)
	for i := 0; i < iterations; i++ {
		ds.AddValue(g.Next())
	}

	// Exponential distribution appears to be a bit more volatile than the
	// others in terms of conforming to expectations, so use a 4% epsilon here,
	// instead of the 2% used for other distributions, to avoid too many false
	// positives.
	ε := 0.04

	wantMean := 1 / rate
	wantStdDev := math.Sqrt(1 / (rate * rate))
	wantMedian := math.Log(2) / rate

	assert.InEpsilon(t, wantMean, ds.ArithmeticMean(), ε,
		"observed mean is outside of acceptable range")
	assert.InEpsilon(t, wantStdDev, ds.SampleStandardDeviation(), ε,
		"observed standard deviation is outside of acceptable range")
	assert.InEpsilon(t, wantMedian, ds.Median(), ε,
		"observed median is outside of acceptable range")
}

func checkPoissonDistribution[U constraints.Unsigned](t *testing.T, g *Poisson[U], wantMean float64) {
	t.Helper()

	const iterations = 50000
	ds := dataset.New(iterations)
	for i := 0; i < iterations; i++ {
		val := g.Next()
		if val < 0 {
			t.Errorf("generated value must be non-negative, got %v", val)
		}
		ds.AddValue(float64(val))
	}

	ε := 0.02

	assert.InEpsilon(t, wantMean, ds.ArithmeticMean(), ε,
		"observed mean is outside of acceptable range")

	// Variance of a Possion distribution equals its mean.
	assert.InEpsilon(t, math.Sqrt(wantMean), ds.SampleStandardDeviation(), ε,
		"observed standard deviation is outside of acceptable range")
}

func checkUniformDistribution[T constraints.Integer | constraints.Float](t *testing.T, g Generator[T], wantMean, wantStddev float64) {
	t.Helper()

	const iterations = 50000
	ds := dataset.New(iterations)
	for i := 0; i < iterations; i++ {
		val := g.Next()
		if val < 0 {
			t.Fatalf("generated value must be non-negative, got %v", val)
		}
		ds.AddValue(float64(val))
	}

	ε := 0.02

	assert.InEpsilon(t, wantMean, ds.ArithmeticMean(), ε,
		"observed mean is outside of acceptable range")

	assert.InEpsilon(t, wantStddev, ds.SampleStandardDeviation(), ε,
		"observed standard deviation is outside of acceptable range")

	assert.InEpsilon(t, wantMean, ds.Median(), ε,
		"observed median outside of acceptable range")
}
