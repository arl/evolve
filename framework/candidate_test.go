package framework

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShuffleCandidates(t *testing.T) {
	candidatesSequence := func(count int) []Candidate {
		// create an ordered sequence of integer candidates
		seq := make([]Candidate, count)
		for i := 0; i < count; i++ {
			seq[i] = i
		}
		return seq
	}

	rng := rand.New(rand.NewSource(99))

	org := candidatesSequence(10)
	shuf := candidatesSequence(10)
	// ensure the slices values are the same
	assert.True(t, assert.ObjectsAreEqualValues(org, shuf))

	// perform shuffling
	ShuffleCandidates(shuf, rng)

	// ensure the slices values are still the same
	assert.Subset(t, org, shuf)
	assert.Subset(t, shuf, org)

	// ensure values are ordered differently
	assert.False(t, assert.ObjectsAreEqualValues(org, shuf))

}
