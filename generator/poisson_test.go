package generator

import (
	"math/rand"
	"testing"

	"github.com/arl/evolve/pkg/mt19937"
)

func TestPoisson(t *testing.T) {
	rng := rand.New(mt19937.New(23))
	const mean = 19
	g := NewPoisson[uint](Const[float64](mean), rng)
	checkPoissonDistribution(t, g, mean)
}

func TestPoissonDynamic(t *testing.T) {
	const initMean float64 = 19

	rng := rand.New(mt19937.New(23))

	gmean := NewSwappable(Const(initMean))
	g := NewPoisson[uint32](gmean, rng)
	checkPoissonDistribution(t, g, initMean)

	const adjustMean float64 = 13
	gmean.Swap(Const(adjustMean))

	checkPoissonDistribution(t, g, adjustMean)
}
