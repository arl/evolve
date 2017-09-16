package termination

import (
	"testing"
	"time"

	"github.com/aurelien-rainone/evolve/framework"
	"github.com/stretchr/testify/assert"
)

func TestElapsedTime(t *testing.T) {
	condition := NewElapsedTime(1 * time.Second)
	data := framework.NewPopulationData(struct{}{}, 0, 0, 0, true, 2, 0, 0, 100*time.Millisecond)

	assert.False(t, condition.ShouldTerminate(data), "Should not terminate before timeout.")
	data = framework.NewPopulationData(struct{}{}, 0, 0, 0, true, 2, 0, 0, time.Second)
	assert.True(t, condition.ShouldTerminate(data), "Should terminate after timeout.")

	// The duration must be greater than zero to be useful. This test ensures
	// that NewElapsedTime panics so that it's impossible to create an invalid
	// termination condition.
	assert.Panics(t, func() { NewElapsedTime(0 * time.Millisecond) })
}
