package base

import "math/rand"

type Candidate interface{}

func ShuffleCandidates(slice []Candidate, rng *rand.Rand) {
	for i := len(slice) - 1; i > 0; i-- {
		j := rng.Intn(i + 1)
		slice[i], slice[j] = slice[j], slice[i]
	}
}
