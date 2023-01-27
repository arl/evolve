package crossover

import "math/rand"

// CX is a Mater for slices representating permutations of a list, implementing
// the Cycle Crossover. The Cycle Crossover was first proposed by Oliver I, it
// builds offspring in such a way that cycles are kept, and copied from parents
// to offsprings.
type CX[T comparable] struct{}

// Mate performs CX crossover on a pair of parent strings and generate a pair of
// offsprings. Mate is undefined if p1 and p2 do not have the same length.
func (p CX[T]) Mate(x1, x2 []T, rng *rand.Rand) (y1, y2 []T) {
	// Create empty chromosomes.
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
