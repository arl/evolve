package api

import "math/rand"

// Candidate is the inteface representing a candidate of a population of
// solutions.
type Candidate interface{}

// ShuffleCandidates shuffles a slice of candidates.
//
// The original slice is modified, though, for commodity, it is returned.
//
// TODO: instead of having ShuffleCandidates, ShuffleEvaluatedPopulation,
// etc... we should have a Shuffler interface and let those types implement
// Shuffler. To do this though, we need to create a Population type representing
// a slice of []Candidate.
func ShuffleCandidates(slice []Candidate, rng *rand.Rand) []Candidate {
	for i := len(slice) - 1; i > 0; i-- {
		j := rng.Intn(i + 1)
		slice[i], slice[j] = slice[j], slice[i]
	}
	return slice
}

// ShuffleEvaluatedPopulation shuffles a slice of evaluated population.
//
// The original slice is modified, though, for commodity, it is returned.
//
// TODO: instead of having ShuffleCandidates, ShuffleEvaluatedPopulation,
// etc... we should have a Shuffler interface and let those types implement
// Shuffler. To do this though, we need to create a Population type representing
// a slice of []Candidate.
func ShuffleEvaluatedPopulation(slice EvaluatedPopulation, rng *rand.Rand) EvaluatedPopulation {
	for i := len(slice) - 1; i > 0; i-- {
		j := rng.Intn(i + 1)
		slice[i], slice[j] = slice[j], slice[i]
	}
	return slice
}
