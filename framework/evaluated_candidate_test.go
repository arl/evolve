package framework

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEvaluatedCandidateEquality(t *testing.T) {
	// Equality is determined only by fitness score, the actual candidate
	// representation is irrelevant. These two candidates should be considered
	// equal.
	candidate1, err1 := NewEvaluatedCandidate("AAAA", 5)
	candidate2, err2 := NewEvaluatedCandidate("BBBB", 5)
	if assert.NoError(t, err1) && assert.NoError(t, err2) {
		assert.Truef(t, candidate1.Equals(candidate1), "Equality must be reflexive.")
		assert.Truef(t, candidate2.Equals(candidate2), "Equality must be reflexive.")

		assert.Equalf(t, candidate1.Hash(), candidate2.Hash(), "Hash codes must be identical for equal objects")
		assert.Zerof(t, candidate1.CompareTo(candidate2), "CompareTo() must be consistent with Equals()")

		assert.Truef(t, candidate1.Equals(candidate2), "Candidates with equal fitness should be equal.")
		assert.Truef(t, candidate2.Equals(candidate1), "Equality must be symmetric.")
	}
}

func TestEvaluatedCandidateNotEqual(t *testing.T) {
	// Equality is determined only by fitness score, the actual candidate
	// representation is irrelevant.  These two candidates should be considered
	// unequal.
	candidate1, err1 := NewEvaluatedCandidate("AAAA", 5)
	candidate2, err2 := NewEvaluatedCandidate("AAAA", 7)
	if assert.NoError(t, err1) && assert.Nil(t, err2) {
		assert.False(t, candidate1.Equals(candidate2), "Candidates with equal fitness should be equal.")
		assert.False(t, candidate2.Equals(candidate1), "Equality must be symmetric.")
		assert.NotZerof(t, candidate1.CompareTo(candidate2), "CompareTo() must be consistent with Equals()")
	}
}

func TestEvaluatedCandidateNullEquality(t *testing.T) {
	candidate, err := NewEvaluatedCandidate("AAAA", 5)
	if assert.NoError(t, err) {
		assert.False(t, candidate.Equals(nil), "Object must not be considered equal to nil pointer.")
	}
}

func TestEvaluatedCandidateComparisons(t *testing.T) {
	// Only test greater than and less than comparisons here since we've already
	// done equality.
	candidate1, err1 := NewEvaluatedCandidate("AAAA", 5)
	candidate2, err2 := NewEvaluatedCandidate("AAAA", 7)
	if assert.NoError(t, err1) && assert.NoError(t, err2) {
		assert.True(t, candidate1.CompareTo(candidate2) < 0, "Incorrect ordering.")
		assert.True(t, candidate2.CompareTo(candidate1) > 0, "Incorrect ordering.")
	}
}

func TestEvaluatedCandidateNegativeFitness(t *testing.T) {
	// should return an error
	candidate, err := NewEvaluatedCandidate("ABC", -1)
	if assert.Error(t, err) {
		assert.Nil(t, candidate)
	}
}
