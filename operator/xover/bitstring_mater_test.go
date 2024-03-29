package xover

import (
	"math/rand"
	"testing"

	"github.com/arl/bitstring"
	"github.com/arl/evolve"
	"github.com/arl/evolve/factory"
	"github.com/arl/evolve/generator"

	"github.com/stretchr/testify/assert"
)

func TestBitstringCrossover(t *testing.T) {
	rng := rand.New(rand.NewSource(99))
	xover := New[*bitstring.Bitstring](BitstringMater{})
	xover.Probability = generator.Const(1.0)
	xover.Points = generator.Const(1)
	f := factory.Bitstring(50)

	pop := evolve.GeneratePopulation[*bitstring.Bitstring](f, 2, rng)
	// Test to make sure that crossover correctly preserves all genetic material
	// originally present in the population and does not introduce anything new.
	want := pop[0].OnesCount() + pop[1].OnesCount()
	for i := 0; i < 50; i++ {
		// Test several generations.
		pop = xover.Apply(pop, rng)

		got := pop[0].OnesCount() + pop[1].OnesCount()
		assert.Equal(t, got, want, "bitstring crossover should not change the total number of set bits in population")
	}
}

func TestBitstringCrossoveWithDifferentLengthParents(t *testing.T) {
	rng := rand.New(rand.NewSource(99))
	xover := New[*bitstring.Bitstring](BitstringMater{})
	xover.Probability = generator.Const(1.0)
	xover.Points = generator.Const(1)

	bs1 := bitstring.Random(32, rng)
	bs2 := bitstring.Random(33, rng)
	pop := []*bitstring.Bitstring{bs1, bs2}

	assert.Panics(t, func() {
		xover.Apply(pop, rng)
	})
}
