package termination

import (
	"testing"

	"github.com/aurelien-rainone/evolve/framework"
	"github.com/stretchr/testify/assert"
)

func TestUserAbort(t *testing.T) {
	condition := NewUserAbort()

	// This population data should be irrelevant.
	data := framework.NewPopulationData(struct{}{}, 0, 0, 0, true, 2, 0, 0, 100)
	assert.False(t, condition.ShouldTerminate(data), "Should not terminate without user abort.")
	assert.False(t, condition.IsAborted(), "Should not be aborted without user intervention.")
	condition.Abort()
	assert.True(t, condition.ShouldTerminate(data), "Should terminate after user abort.")
	assert.True(t, condition.IsAborted(), "Should be aborted after user intervention.")
}
