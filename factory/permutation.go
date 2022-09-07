package factory

import (
	"math/rand"
)

// A permutation factory generates shuffled variations of an original slice.
type Permutation[T any] []T

func (f Permutation[T]) New(rng *rand.Rand) []T {
	// We need to create a copy of the original before shuffling it.
	cpy := make([]T, len(f))
	copy(cpy, f)
	rng.Shuffle(len(cpy), func(i, j int) {
		cpy[i], cpy[j] = cpy[j], cpy[i]
	})

	return cpy
}
