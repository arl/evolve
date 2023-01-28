package generator

import (
	"math/rand"
	"testing"

	"github.com/arl/evolve/pkg/mt19937"
)

func TestGaussian(t *testing.T) {
	const mean, stddev float64 = 147, 17

	rng := rand.New(mt19937.New())

	g := NewGaussian(Const(mean), Const(stddev), rng)
	checkGaussianDistribution(t, g, mean, stddev)
}

func TestGaussianDynamic(t *testing.T) {
	const mean, stddev float64 = 147, 17

	rng := rand.New(mt19937.New())

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
