package crossover

import "math/rand"

// PMX is a Mater for slices representing permutations of a list, implementing
// the Partially Mapped Crossover algorithm or PMX.
//
// It creates offsprings that are permutations of the parents by choosing 2
// random crossover points and exchanging elements positions.
type PMX[T comparable] struct{}

// Mate performs PMX crossover on a pair of parent strings and generate a pair
// of offsprings. Mate is undefined if p1 and p2 do not have the same length.
func (p PMX[T]) Mate(p1, p2 []T, rng *rand.Rand) (off1, off2 []T) {
	// Create identical copies.
	off1 = make([]T, len(p1))
	off2 = make([]T, len(p1))
	copy(off1, p1)
	copy(off2, p2)

	// Generate the 2 cut points.
	pt1, pt2 := rng.Intn(len(p1)), rng.Intn(len(p1))
	mapBasedPMX(p1, p2, off1, off2, pt1, pt2)
	return
}

func mapBasedPMX[T comparable](p1, p2, off1, off2 []T, pt1, pt2 int) {
	length := pt2 - pt1
	if length < 0 {
		length += len(p1)
	}

	m1 := make(map[T]T, len(p1)*2)
	m2 := make(map[T]T, len(p1)*2)

	for i := 0; i < length; i++ {
		index := (i + pt1) % len(p1)
		item1 := off1[index]
		item2 := off2[index]
		off1[index] = item2
		off2[index] = item1
		m1[item1] = item2
		m2[item2] = item1
	}

	checkUnmappedElements(off1, m2, pt1, pt2)
	checkUnmappedElements(off2, m1, pt1, pt2)
}

// checks elements that are outside of the partially mapped section to see if
// there are any duplicate items in the list. If there are, they are mapped
// appropriately.
func checkUnmappedElements[T comparable](offspring []T, m map[T]T, start, end int) {
	for i := range offspring {
		if isInsideMappedRegion(i, start, end) {
			continue
		}
		mapped := offspring[i]
		for {
			_, ok := m[mapped]
			if !ok {
				break
			}
			mapped = m[mapped]
		}
		offspring[i] = mapped
	}
}

// checks whether a given list position is within the partially mapped region
// used for pmx. pos is the position to check start is the (inclusive) start
// index of the mapped region end is the (exclusive) end index of the mapped
// region.
func isInsideMappedRegion(pos, start, end int) bool {
	enclosed := pos < end && pos >= start
	wrapAround := start > end && (pos >= start || pos < end)
	return enclosed || wrapAround
}
