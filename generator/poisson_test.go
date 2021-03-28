package generator

import (
	"math"
	"math/rand"
	"testing"

	"github.com/arl/evolve"
	"github.com/arl/evolve/pkg/mt19937"

	"github.com/stretchr/testify/assert"
)

func TestPoisson(t *testing.T) {
	rng := rand.New(mt19937.New(23))
	const mean = 19
	g := NewPoisson(ConstFloat64(mean), rng)
	checkPoissonDistribution(t, g, mean)
}

func TestPoissonDynamic(t *testing.T) {
	const initMean = 19

	rng := rand.New(mt19937.New(23))

	gmean := NewAdjustableFloat(initMean)
	g := NewPoisson(gmean, rng)
	checkPoissonDistribution(t, g, initMean)

	const adjustMean = 13
	gmean.Set(adjustMean)

	checkPoissonDistribution(t, g, adjustMean)
}

func checkPoissonDistribution(t *testing.T, g Int, wantMean float64) {
	t.Helper()

	const iterations = 10000
	ds := evolve.NewDataset(iterations)
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
