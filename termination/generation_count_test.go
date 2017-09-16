package termination

import (
	"testing"
	"time"

	"github.com/aurelien-rainone/evolve/framework"
	"github.com/stretchr/testify/assert"
)

func TestGenerationCount(t *testing.T) {
	condition := NewGenerationCount(5)
	data := framework.NewPopulationData(struct{}{}, 0, 0, 0, true, 2, 0, 3, 100*time.Millisecond)

	// Generation number 3 is the 4th generation (generation numbers are zero-frameworkd).
	assert.False(t, condition.ShouldTerminate(data), "Should not terminate after 4th generation.")
	data = framework.NewPopulationData(struct{}{}, 0, 0, 0, true, 2, 0, 4, 100*time.Millisecond)
	assert.True(t, condition.ShouldTerminate(data), "Should terminate after 5th generation.")

	// The generation count must be greater than zero to be useful. This test
	// ensures that NewGenerationCount panics so that it's impossible to create
	// an invalid terminate condition.
	assert.Panics(t, func() { NewGenerationCount(0) })
}
