package termination

import (
	"testing"
	"time"

	"github.com/aurelien-rainone/evolve/pkg/api"
	"github.com/stretchr/testify/assert"
)

func TestElapsedTime(t *testing.T) {
	cond := 1 * ElapsedTime(time.Second)
	data := api.NewPopulationData(struct{}{}, 0, 0, 0, true, 2, 0, 0, 100*time.Millisecond)

	assert.False(t, cond.ShouldTerminate(data), "Should not terminate before timeout.")
	data = api.NewPopulationData(struct{}{}, 0, 0, 0, true, 2, 0, 0, time.Second)
	assert.True(t, cond.ShouldTerminate(data), "Should terminate after timeout.")
}
