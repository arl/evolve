package generator

import (
	"math/rand"
	"testing"

	"github.com/arl/evolve/pkg/mt19937"
)

func TestExponential(t *testing.T) {
	rng := rand.New(mt19937.New())
	const rate float64 = 3.2

	g := NewExponential(Const(rate), rng)
	checkExponentialDistribution(t, g, rate)
}

func TestExponentialDynamic(t *testing.T) {
	const initRate = 0.75

	rng := rand.New(mt19937.New())

	grate := NewSwappable(Const(initRate))
	g := NewExponential(grate, rng)
	checkExponentialDistribution(t, g, initRate)

	const adjustRate = 1.05
	grate.Swap(Const(adjustRate))

	checkExponentialDistribution(t, g, adjustRate)
}
