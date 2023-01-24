package crossover

import "math/rand"

// The CX or Cycle Crossover, first proposed by Oliver I, builds offspring in
// such a way that cycles are kept, and copied from parents to offsprings.
type CX[T comparable] struct{}

// Mate mates 2 parents and generates a pair of offsprings with CX. the number
// of cut points is unused.
func (p CX[T]) Mate(x1, x2 []T, nxpts int, rng *rand.Rand) (y1, y2 []T) {
	if len(x1) != len(x2) {
		panic("CX cannot mate parents of different lengths")
	}

	// Create empty genotypes.
	y1 = make([]T, len(x1))
	y2 = make([]T, len(x1))

	y1[0] = x1[0]
	y2[0] = x2[0]

	// Keep track of initialized indices in offsprings
	init := make(map[int]struct{})
	init[0] = struct{}{}

	// Map x1 values to their index
	x1pos := make(map[T]int)
	for idx, val := range x1 {
		x1pos[val] = idx
	}

	// Copy cycle to children.
	i := 0
	for {
		// j is the index of x2[i] in x1
		j := x1pos[x2[i]]
		if j == 0 {
			// cycle end
			break
		}
		y1[j] = x1[j]
		y2[j] = x2[j]
		i = j
		init[i] = struct{}{}
	}

	// Copy untouched positions from p1 to c2 and from p2 to c1.
	for i := range x1 {
		if _, ok := init[i]; !ok {
			y1[i] = x2[i]
			y2[i] = x1[i]
		}
	}

	return
}
