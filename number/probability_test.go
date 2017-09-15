package number

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProbability(t *testing.T) {
	t.Run("a negative probability value is invalid", func(*testing.T) {
		_, err := NewProbability(-1)
		assert.Error(t, err)
	})

	t.Run("a probability value greater than 1 is invalid", func(*testing.T) {
		_, err := NewProbability(1.01)
		assert.Error(t, err)
	})

	t.Run("a probability of 0 is valid", func(*testing.T) {
		p, err := NewProbability(0)
		assert.NoError(t, err)
		assert.Equal(t, p, ProbabilityZero)
	})

	t.Run("a probability of 1 is valid", func(*testing.T) {
		p, err := NewProbability(1)
		assert.NoError(t, err)
		assert.Equal(t, p, ProbabilityOne)
	})
}

func TestProbabilityEvents(t *testing.T) {
	const iterations = 1000
	rng := rand.New(rand.NewSource(99))

	t.Run("an event with a probability of 0 never happens", func(*testing.T) {
		for cur := 0; cur < iterations; cur++ {
			assert.False(t, ProbabilityZero.NextEvent(rng))
		}
	})

	t.Run("an event with a probability of 1 always happens", func(*testing.T) {
		for cur := 0; cur < iterations; cur++ {
			assert.True(t, ProbabilityOne.NextEvent(rng))
		}
	})

	t.Run("an event with a probability of 0.5 happen half of the time on average", func(*testing.T) {
		var trues, falses int
		for cur := 0; cur < iterations; cur++ {
			if ProbabilityEven.NextEvent(rng) {
				trues++
			} else {
				falses++
			}
		}
		trueAvg := float32(trues) / float32(iterations)
		assert.Truef(t, trueAvg > 0.45 && trueAvg < 0.55,
			"want on average, 0.45 < e < 0.55, got e = %f", trueAvg)
	})
}

func TestProbabilityComplement(t *testing.T) {
	assert.Equal(t, ProbabilityZero.Complement(), ProbabilityOne,
		"if p = 0, want complement = 1, got %v", ProbabilityZero.Complement())
	assert.Equal(t, ProbabilityOne.Complement(), ProbabilityZero,
		"if p = 1, want complement = 0, got %v", ProbabilityOne.Complement())
	assert.Equal(t, ProbabilityEven.Complement(), ProbabilityEven,
		"if p = 0.5, want complement = 0.5, got %v", ProbabilityEven.Complement())
}
