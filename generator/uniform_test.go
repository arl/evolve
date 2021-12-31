package generator

import (
	"math"
	"math/rand"
	"testing"

	"github.com/arl/evolve"
	"github.com/arl/evolve/pkg/mt19937"
	"github.com/stretchr/testify/assert"
)

func TestUniformInt(t *testing.T) {
	const (
		min, max = 150, 500
		rang     = max - min
	)

	const mean = (rang / 2) + min

	// For a uniformly distributed [0 n] range, standard deviation should be n/sqrt(12).
	stddev := float64(rang) / math.Sqrt(12)

	rng := rand.New(mt19937.New(0))
	g := NewUniformtInt[uint16](min, max, rng)

	checkUniformDistribution(t, g, mean, stddev)
}

// func TestUniformAcrossZero(t *testing.T) {
// 	const (
// 		min, max = -150, 500
// 		rang     = max - min
// 	)

// 	const mean = (rang / 2) + min

// 	// For a uniformly distributed [0 n] range, standard deviation should be n/sqrt(12).
// 	stddev := float64(rang) / math.Sqrt(12)

// 	rng := rand.New(mt19937.New(0))
// 	g := NewUniformtInt[int16](min, max, rng)

// 	checkUniformDistribution(t, g, mean, stddev)
// }

func checkUniformDistribution[U Unsigned](t *testing.T, g *UniformInt[U], wantMean, wantStddev float64) {
	t.Helper()

	const iterations = 10000
	ds := evolve.NewDataset(iterations)
	for i := 0; i < iterations; i++ {
		val := g.Next()
		// if val < 0 {
		// 	t.Errorf("generated value must be non-negative, got %v", val)
		// }
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
