package operator

import (
	"math/rand"
	"testing"

	"github.com/arl/evolve"
)

// adjustInt mutates integers candidates by adding a fixed offset.
type adjustInt int

func (op adjustInt) Apply(pop *evolve.Population[int], rng *rand.Rand) {
	for i, c := range pop.Candidates {
		pop.Candidates[i] = c + int(op)
	}
}

func TestEvolutionPipeline(t *testing.T) {
	// Make sure that multiple operators in a pipeline are applied correctly
	// to the population and validate the cumulative effects.
	rng := rand.New(rand.NewSource(99))
	pop := evolve.NewPopulation[int](10, nil)
	for i := range pop.Candidates {
		pop.Candidates[i] = 10 + i*10
	}

	pipe := Pipeline[int]{adjustInt(1), adjustInt(3)}
	pipe.Apply(pop, rng)

	// Net result should be each candidate increased by 4.
	var sum int
	for _, c := range pop.Candidates {
		ic := c
		sum += ic
		if ic%10 != 4 {
			t.Error("candidate value should have increased by 4, got", ic)
		}
	}
	if sum != 590 {
		t.Error("want sum = 590, got", sum)
	}
}
