package operator

import (
	"math/rand"
	"testing"

	"github.com/aurelien-rainone/evolve/internal/test"
	"github.com/aurelien-rainone/evolve/pkg/api"
)

func TestEvolutionPipeline(t *testing.T) {
	// Make sure that multiple operators in a pipeline are applied correctly
	// to the population and validate the cumulative effects.
	rng := rand.New(rand.NewSource(99))
	pop := make([]api.Candidate, 0, 10)
	for i := 10; i <= 100; i += 10 {
		pop = append(pop, i)
	}
	// Increment 30% of the numbers and decrement the other 70%.
	pipe := Pipeline{test.IntegerAdjuster(1), test.IntegerAdjuster(3)}
	pop = pipe.Apply(pop, rng)
	// Net result should be each candidate increased by 4.
	var aggregate int
	for _, c := range pop {
		ic := c.(int)
		aggregate += ic
		if ic%10 != 4 {
			t.Error("candidate value should have increased by 4, got", ic)
		}
	}
	if aggregate != 590 {
		t.Error("want aggregate to be 590, got", aggregate)
	}
}
