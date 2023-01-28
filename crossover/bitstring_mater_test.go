package crossover

import (
	"math/rand"
	"testing"

	"github.com/arl/bitstring"
	"github.com/arl/evolve"
	"github.com/arl/evolve/factory"
	"github.com/arl/evolve/generator"
)

func TestBitstringCrossover(t *testing.T) {
	rng := rand.New(rand.NewSource(99))
	xover := evolve.Crossover[*bitstring.Bitstring]{
		Probability: generator.Const(1.0),
		Mater: &BitstringMater{
			Points: generator.Const(1),
		},
	}

	f := factory.Bitstring(50)

	pop := evolve.GeneratePopulation[*bitstring.Bitstring](2, f, nil, rng)

	// Test to make sure that crossover correctly preserves all genetic material
	// originally present in the population and does not introduce anything new.
	want := pop.Candidates[0].OnesCount() + pop.Candidates[1].OnesCount()
	for i := 0; i < 50; i++ {
		// Test several generations.
		xover.Apply(pop, rng)

		if got := pop.Candidates[0].OnesCount() + pop.Candidates[1].OnesCount(); got != want {
			t.Errorf("bitstring crossover should not change the total number of set bits in population")
		}
	}
}
