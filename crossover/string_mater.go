package crossover

import (
	"math/rand"

	"github.com/arl/evolve/generator"
)

// StringMater is a crossover helper that mates pairs of parent strings and
// produces pairs of offsprings.
type StringMater struct {
	// Points generator decided the number of cut points to apply, for each mating.
	Points generator.Generator[int]
}

// Mate performs crossover on a pair of parent strings and generate a pair of
// offsprings. Mate is undefined if p1 and p2 do not have the same length (in bytes).
func (m *StringMater) Mate(p1, p2 string, rng *rand.Rand) (string, string) {
	// Decide the number of cut points.
	npts := int(m.Points.Next())

	off1, off2 := []byte(p1), []byte(p2)

	// Apply as many crossovers as required.
	for i := 0; i < npts; i++ {
		// Cross-over index is always greater than zero and less than the length
		// of the parent so that we always pick a point that will result in a
		// meaningful crossover.
		xidx := (1 + rng.Intn(len(p1)-1))
		for j := 0; j < xidx; j++ {
			// swap elements j of both offsprings
			off1[j], off2[j] = off2[j], off1[j]
		}
	}
	return string(off1), string(off2)
}
