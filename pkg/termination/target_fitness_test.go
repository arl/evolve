package termination

import (
	"testing"

	"github.com/aurelien-rainone/evolve/pkg/api"
	"github.com/stretchr/testify/assert"
)

// Unit test for termination condition that checks the best fitness attained so far
// against a pre-determined target.
func TestTargetFitness(t *testing.T) {

	t.Run("natural fitness", func(t *testing.T) {
		var condition api.TerminationCondition = NewTargetFitness(10.0, true)
		data := api.NewPopulationData(struct{}{}, 5.0, 4.0, 0, true, 2, 0, 0, 100)
		assert.False(t, condition.ShouldTerminate(data), "Should not terminate before target fitness is reached")
		data = api.NewPopulationData(struct{}{}, 10.0, 8.0, 0, true, 2, 0, 0, 100)
		assert.True(t, condition.ShouldTerminate(data), "Should terminate once target fitness is reached")
	})

	t.Run("non-natural fitness", func(t *testing.T) {
		var condition api.TerminationCondition = NewTargetFitness(1.0, false)
		data := api.NewPopulationData(struct{}{}, 5.0, 4.0, 0, true, 2, 0, 0, 100)
		assert.False(t, condition.ShouldTerminate(data), "Should not terminate before target fitness is reached")
		data = api.NewPopulationData(struct{}{}, 1.0, 3.1, 0, true, 2, 0, 0, 100)
		assert.True(t, condition.ShouldTerminate(data), "Should terminate once target fitness is reached")
	})
}
