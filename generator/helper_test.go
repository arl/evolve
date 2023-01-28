package generator

import (
	"errors"
	"fmt"
	"math"
	"testing"

	"github.com/arl/evolve/pkg/dataset"
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

	assertInEpsilon(t, wantMean, ds.ArithmeticMean(), ε, "mean")
	// Expected median is the same as expected mean.
	assertInEpsilon(t, wantMean, ds.Median(), ε, "median")
	assertInEpsilon(t, wantStdDev, ds.SampleStandardDeviation(), ε, "sample stddev")
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

	assertInEpsilon(t, wantMean, ds.ArithmeticMean(), ε, "mean")
	assertInEpsilon(t, wantStdDev, ds.SampleStandardDeviation(), ε, "sample stddev")
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

	assertInEpsilon(t, wantMean, ds.ArithmeticMean(), ε, "mean")
	assertInEpsilon(t, wantStdDev, ds.SampleStandardDeviation(), ε, "sample stddev")
	assertInEpsilon(t, wantMedian, ds.Median(), ε, "median")
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

	assertInEpsilon(t, wantMean, ds.ArithmeticMean(), ε, "mean")
	// Variance of a Possion distribution equals its mean.
	assertInEpsilon(t, math.Sqrt(wantMean), ds.SampleStandardDeviation(), ε, "sample stddev")
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

	assertInEpsilon(t, wantMean, ds.ArithmeticMean(), ε, "mean")
	assertInEpsilon(t, wantStddev, ds.SampleStandardDeviation(), ε, "sample stddev")
	assertInEpsilon(t, wantMean, ds.Median(), ε, "median")
}

// assertInEpsilon checks that the relative error between the value we want and
// the value we got is less than epsilon. observed is the name of the observed
// value and is used in the error message.
func assertInEpsilon(tb testing.TB, want, got, epsilon float64, observed string) {
	tb.Helper()

	if math.IsNaN(epsilon) {
		tb.Fatalf("assertInEpsilon(%s), epsilon must not be NaN", observed)
	}
	actualEpsilon, err := calcRelativeError(want, got)
	if err != nil {
		tb.Fatalf("assertInEpsilon(%s), %v", observed, err.Error())
	}
	if actualEpsilon > epsilon {
		tb.Fatalf("assertInEpsilon(%s), relative error is too high.\n"+
			" (actual) %f > %f (expected)", observed, actualEpsilon, epsilon)
	}
}

func calcRelativeError(want, got float64) (float64, error) {
	if math.IsNaN(want) && math.IsNaN(got) {
		return 0, nil
	}
	if math.IsNaN(want) {
		return 0, errors.New("expected value must not be NaN")
	}
	if want == 0 {
		return 0, fmt.Errorf("expected value must have a value other than zero to calculate the relative error")
	}
	if math.IsNaN(got) {
		return 0, errors.New("actual value must not be NaN")
	}

	return math.Abs(want-got) / math.Abs(want), nil
}
