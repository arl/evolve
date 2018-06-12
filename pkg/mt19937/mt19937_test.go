package mt19937

import (
	"math"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Test to ensure that two distinct RNGs with the same seed return the same
// sequence of numbers.
func TestMT19937Repeatability(t *testing.T) {

	// create common seed
	seed := int64(time.Now().UnixNano())

	// create first mersenne twister rng
	rng := rand.New(New(seed))

	// create second one using same seed.
	other := rand.New(New(seed))

	// read sequences of 100 bytes on both generators
	data1 := make([]byte, 100)
	rng.Read(data1)

	data2 := make([]byte, 100)
	other.Read(data2)
	assert.EqualValues(t, data1, data2, "generated sequences should match")
}

// Test to ensure that the output from the RNG is broadly as expected. This will
// not detect the subtle statistical anomalies that would be picked up by
// Diehard, but it provides a simple check for major problems with the output.
func TestMT19937Distribution(t *testing.T) {
	seed := int64(time.Now().UnixNano())
	rng := rand.New(New(seed))

	pi := monteCarloValueForPi(rng, 1000000)
	assert.InDeltaf(t, pi, math.Pi, 0.01, "Monte Carlo value for Pi is outside acceptable range")
}

// Test to ensure that the output from the RNG is broadly as expected. This will
// not detect the subtle statistical anomalies that would be picked up by
// Diehard, but it provides a simple check for major problems with the output.
func TestMT19937StdDev(t *testing.T) {
	seed := int64(time.Now().UnixNano())
	rng := rand.New(New(seed))

	// Expected standard deviation for a uniformly distributed
	// population of values in the range 0..n approaches n/sqrt(12).
	const n = 100
	got := calculareSampleStdDev(rng, n, 10000000)
	want := 100 / math.Sqrt(12)
	assert.InDeltaf(t, got, want, 0.02, "standard deviation outside acceptable range")
}
