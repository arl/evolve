package xover

import (
	"math/rand"
	"testing"

	"github.com/arl/evolve"
	"github.com/arl/evolve/generator"
	"github.com/arl/evolve/pkg/bitstring"

	"github.com/stretchr/testify/assert"
)

func TestBitstringCrossover(t *testing.T) {
	rng := rand.New(rand.NewSource(99))
	xover := New(BitstringMater{})
	xover.Probability = generator.ConstFloat64(1)
	xover.Points = generator.ConstInt(1)
	f := generator.Bitstring(50)

	pop := evolve.GeneratePopulation(f, 2, rng)
	// Test to make sure that crossover correctly preserves all genetic material
	// originally present in the population and does not introduce anything new.
	want := pop[0].(*bitstring.Bitstring).OnesCount() +
		pop[1].(*bitstring.Bitstring).OnesCount()
	for i := 0; i < 50; i++ {
		// Test several generations.
		pop = xover.Apply(pop, rng)

		got := pop[0].(*bitstring.Bitstring).OnesCount() +
			pop[1].(*bitstring.Bitstring).OnesCount()
		assert.Equal(t, got, want, "bitstring crossover should not change the total number of set bits in population")
	}
}

func TestBitstringCrossoveWithDifferentLengthParents(t *testing.T) {
	// The BitstringCrossover operator is only defined to work on populations
	// containing Strings of equal lengths. Any attempt to apply the operation
	// to populations that contain different length strings should panic. Not
	// panicking should be considered a bug since it could lead to hard to trace
	// bugs elsewhere.
	rng := rand.New(rand.NewSource(99))
	xover := New(BitstringMater{})
	xover.Probability = generator.ConstFloat64(1)
	xover.Points = generator.ConstInt(1)

	bs1 := bitstring.Random(32, rng)
	bs2 := bitstring.Random(33, rng)
	pop := []interface{}{bs1, bs2}

	assert.Panics(t, func() {
		// This should panic since the parents are different lengths.
		xover.Apply(pop, rng)
	})
}
