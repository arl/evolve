package operators

import (
	"math/rand"
	"testing"

	"github.com/aurelien-rainone/evolve/framework"
	"github.com/aurelien-rainone/evolve/internal/test"
	"github.com/stretchr/testify/assert"
)

func TestEvolutionPipeline(t *testing.T) {
	// Make sure that multiple operators in a pipeline are applied correctly
	// to the population and validate the cumulative effects.
	rng := rand.New(rand.NewSource(99))
	population := make([]framework.Candidate, 0, 10)
	for i := 10; i <= 100; i += 10 {
		population = append(population, i)
	}
	// Increment 30% of the numbers and decrement the other 70%.
	operators := make([]framework.EvolutionaryOperator, 2)
	operators[0] = test.IntegerAdjuster(1)
	operators[1] = test.IntegerAdjuster(3)
	evolutionScheme, err := NewEvolutionPipeline(operators...)
	if assert.NoError(t, err) {
		population = evolutionScheme.Apply(population, rng)
		// Net result should be each candidate increased by 4.
		var aggregate int
		for _, c := range population {
			ic := c.(int)
			aggregate += ic
			assert.Equal(t, ic%10, 4, "Candidate should have increased by 4, is", ic)
		}
		assert.Equal(t, aggregate, 590, "Aggregate should be 590 after mutations, is", aggregate)
	}
}

func TestEvolutionPipelineEmptyPipeline(t *testing.T) {
	// An empty pipeline is not allowed, an error should be returned.
	ep, err := NewEvolutionPipeline()
	assert.Nil(t, ep)
	assert.Error(t, err)
}
