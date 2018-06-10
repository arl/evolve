package random

import (
	"math"
	"math/rand"

	"github.com/arl/evolve/pkg/api"
)

// Provides methods used for testing the operation of RNG implementations.
// This is a rudimentary check to ensure that the output of a given RNG is
// approximately uniformly distributed. If the RNG output is not uniformly
// distributed, this method will return a poor estimate for the
// value of pi.
// - rng is the RNG to test.
// - niter is the number of random points to generate for use in the
// calculation. This value needs to be sufficiently large in order to
// produce a reasonably accurate result (assuming the RNG is uniform).
// Less than 10,000 is not particularly useful.  100,000 should be sufficient.
//
// Returns an approximation of pi generated using the provided RNG.
func monteCarloValueForPi(rng *rand.Rand, niter int) float64 {
	// Assumes a quadrant of a circle of radius 1, bounded by a box with
	// sides of length 1.  The area of the square is therefore 1 square unit
	// and the area of the quadrant is (pi * r^2) / 4.
	var tot int
	// Generate the specified number of random points and count how many fall
	// within the quadrant and how many do not. We expect the number of points
	// in the quadrant (expressed as a fraction of the total number of points)
	// to be pi/4. Therefore pi = 4 * ratio.
	for i := 0; i < niter; i++ {
		x := rng.Float64()
		y := rng.Float64()
		if isInQuadrant(x, y) {
			tot++
		}
	}
	// From these figures we can deduce an approximate value for Pi.
	return 4 * float64(tot) / float64(niter)
}

// Uses Pythagoras' theorem to determine whether the specified coordinates fall
// within the area of the quadrant of a circle of radius 1 that is centered on
// the origin.
// - x, y are the coordinates of the point (must be between 0 and 1).
//
// Returns true if the point is within the quadrant, false otherwise.
func isInQuadrant(x, y float64) bool {
	distance := math.Sqrt((x * x) + (y * y))
	return distance <= 1
}

// Generates a sequence of values from a given random number generator and
// then calculates the standard deviation of the sample.
// - rng is the RNG to use.
// - max is the maximum value for generated integers (values will be in the
// range [0, max)).
// - niter is the number of values to generate and use in the standard deviation
// calculation.
//
// Returns the standard deviation of the generated sample.
func calculareSampleStdDev(rng *rand.Rand, max int64, niter int) float64 {
	dataset := api.NewDataset(niter)
	for i := 0; i < niter; i++ {
		dataset.AddValue(float64(rng.Int63n(max)))
	}
	return dataset.SampleStandardDeviation()
}
