package random

import (
	"math"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMersenneTwister(t *testing.T) {

	t.Run("repeatability", func(t *testing.T) {
		// Test to ensure that two distinct RNGs with the same seed return the
		// same sequence of numbers.

		// create common seed
		var seed = randomSeed()

		// create first mersenne twister rng
		rng := rand.New(NewMT19937(seed))

		// create second one using same seed.
		other := rand.New(NewMT19937(seed))

		// read sequences of 100 bytes on both generators
		data1 := make([]byte, 100)
		rng.Read(data1)

		data2 := make([]byte, 100)
		other.Read(data2)
		assert.EqualValues(t, data1, data2, "generated sequences should match")
	})

	t.Run("distribution", func(t *testing.T) {
		// Test to ensure that the output from the RNG is broadly as
		// expected. This will not detect the subtle statistical
		// anomalies that would be picked up by Diehard, but it provides
		// a simple check for major problems with the output.
		rng := rand.New(NewMT19937(randomSeed()))

		pi := calculateMonteCarloValueForPi(rng, 100000)
		assert.InDeltaf(t, pi, math.Pi, 0.01, "Monte Carlo value for Pi is outside acceptable range: %v", pi)
	})

	t.Run("standard deviation", func(t *testing.T) {
		// Test to ensure that the output from the RNG is broadly as
		// expected. This will not detect the subtle statistical
		// anomalies that would be picked up by Diehard, but it provides
		// a simple check for major problems with the output.
		rng := rand.New(NewMT19937(randomSeed()))

		// Expected standard deviation for a uniformly distributed
		// population of values in the range 0..n approaches n/sqrt(12).
		const n = 100
		observedSD := calculateSampleStandardDeviation(rng, n, 1000000)
		expectedSD := 100 / math.Sqrt(12)
		assert.InDeltaf(t, observedSD, expectedSD, 0.02, "standard deviation outside acceptable range: %f", observedSD)
	})
}

func randomSeed() int64 {
	return int64(time.Now().UnixNano())
}
