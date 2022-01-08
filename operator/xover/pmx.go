package xover

import "math/rand"

// PMX implements the partially mapped crossover algorithm.
//
// This crossover is indicated when chromomes are lists of a predefined set of
// elements. It creates offsprings that are non-repeating permutations of the
// parents by choosing 2 random crossover points and exchanging elements
// positions.
type PMX[T comparable] struct{}

// Mate mates 2 parents and generates a pair of offsprings.
//
// parent1 and parent2 are the two individuals that provides the source
// material for generating offspring.
func (p PMX[T]) Mate(p1, p2 []T, nxpts int, rng *rand.Rand) [][]T {
	if nxpts != 2 {
		panic("PMX is only defined for 2 cut points")
	}

	if len(p1) != len(p2) {
		panic("PMX cannot mate parents of different lengths")
	}

	plen := len(p1)

	offsp1 := make([]T, plen)
	offsp2 := make([]T, plen)
	copy(offsp1, p1)
	copy(offsp2, p2)

	pt1, pt2 := rng.Intn(plen), rng.Intn(plen)
	length := pt2 - pt1
	if length < 0 {
		length += len(p1)
	}

	m1 := make(map[T]T, plen*2)
	m2 := make(map[T]T, plen*2)

	for i := 0; i < length; i++ {
		index := (i + pt1) % plen
		item1 := offsp1[index]
		item2 := offsp2[index]
		offsp1[index] = item2
		offsp2[index] = item1
		m1[item1] = item2
		m2[item2] = item1
	}

	p.checkUnmappedElements(offsp1, m2, pt1, pt2)
	p.checkUnmappedElements(offsp2, m1, pt1, pt2)

	return [][]T{offsp1, offsp2}
}

// checks elements that are outside of the partially mapped section to see if
// there are any duplicate items in the list. If there are, they are mapped
// appropriately.
func (p PMX[T]) checkUnmappedElements(offspring []T, mapping map[T]T, mapStart, mapEnd int) {
	for i := range offspring {
		if !p.isInsideMappedRegion(i, mapStart, mapEnd) {
			mapped := offspring[i]
			for {
				_, ok := mapping[mapped]
				if !ok {
					break
				}

				mapped = mapping[mapped]
			}
			offspring[i] = mapped
		}
	}
}

// checks whether a given list position is within the partially mapped region used for crossover.
// pos is the position to check
// start is the (inclusive) start index of the mapped region
// end is the (exclusive) end index of the mapped region
func (p PMX[T]) isInsideMappedRegion(pos, start, end int) bool {
	enclosed := pos < end && pos >= start
	wrapAround := start > end && (pos >= start || pos < end)
	return enclosed || wrapAround
}
