package framework

import "math/rand"

// Candidate is the inteface representing a candidate of a population of
// solutions.
type Candidate interface{}

// ShuffleCandidates shuffles a slice of candidates.
//
// The original slice is modified, though, for commodity, it is returned.
func ShuffleCandidates(slice []Candidate, rng *rand.Rand) []Candidate {
	for i := len(slice) - 1; i > 0; i-- {
		j := rng.Intn(i + 1)
		slice[i], slice[j] = slice[j], slice[i]
	}
	return slice
}
