package generator

import (
	"math"
	"math/rand"
	"testing"

	"github.com/arl/evolve"
	"github.com/arl/evolve/pkg/mt19937"

	"github.com/stretchr/testify/assert"
)

func TestExponential(t *testing.T) {
	rng := rand.New(mt19937.New(23))
	const rate float64 = 3.2

	g := NewExponential(Const(rate), rng)
	checkExponentialDistribution(t, g, rate)
}

func TestExponentialDynamic(t *testing.T) {
	const initRate = 0.75

	rng := rand.New(mt19937.New(23))

	grate := NewSwappable(Const(initRate))
	g := NewExponential(grate, rng)
	checkExponentialDistribution(t, g, initRate)

	const adjustRate = 1.05
	grate.Swap(Const(adjustRate))

	checkExponentialDistribution(t, g, adjustRate)
}

func checkExponentialDistribution(t *testing.T, g *Exponential, rate float64) {
	t.Helper()

	const iterations = 10000
	ds := evolve.NewDataset(iterations)
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
