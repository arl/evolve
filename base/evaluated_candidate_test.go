package base

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEvaluatedCandidateEquality(t *testing.T) {
	// Equality is determined only by fitness score, the actual candidate
	// representation is irrelevant. These two candidates should be considered
	// equal.
	var (
		candidate1, candidate2 *EvaluatedCandidate
		err                    error
	)
	candidate1, err = NewEvaluatedCandidate("AAAA", 5)
	assert.Nil(t, err)
	candidate2, err = NewEvaluatedCandidate("BBBB", 5)
	assert.Nil(t, err)

	assert.Truef(t, candidate1.Equals(candidate1), "Equality must be reflexive.")
	assert.Truef(t, candidate2.Equals(candidate2), "Equality must be reflexive.")

	//assert candidate1.hashCode() == candidate2.hashCode() : "Hash codes must be identical for equal objects.";
	//assert candidate1.compareTo(candidate2) == 0 : "compareTo() must be consistent with equals()";

	assert.Truef(t, candidate1.Equals(candidate2), "Candidates with equal fitness should be equal.")
	assert.Truef(t, candidate2.Equals(candidate1), "Equality must be symmetric.")
}

func TestEvaluatedCandidateNotEqual(t *testing.T) {
	// Equality is determined only by fitness score, the actual candidate
	// representation is irrelevant.  These two candidates should be considered
	// unequal.
	var (
		candidate1, candidate2 *EvaluatedCandidate
		err                    error
	)
	candidate1, err = NewEvaluatedCandidate("AAAA", 5)
	assert.Nil(t, err)
	candidate2, err = NewEvaluatedCandidate("AAAA", 7)
	assert.Nil(t, err)

	assert.False(t, candidate1.Equals(candidate2), "Candidates with equal fitness should be equal.")
	assert.False(t, candidate2.Equals(candidate1), "Equality must be symmetric.")

	//assert candidate1.compareTo(candidate2) != 0 : "compareTo() must be consistent with equals()";
}

func TestEvaluatedCandidateNullEquality(t *testing.T) {
	candidate, err := NewEvaluatedCandidate("AAAA", 5)
	assert.Nil(t, err)
	assert.False(t, candidate.Equals(nil), "Object must not be considered equal to nil pointer.")
}

func TestEvaluatedCandidateComparisons(t *testing.T) {
	// Only test greater than and less than comparisons here since we've already
	// done equality.

	var (
		candidate1, candidate2 *EvaluatedCandidate
		err                    error
	)
	candidate1, err = NewEvaluatedCandidate("AAAA", 5)
	assert.Nil(t, err)
	candidate2, err = NewEvaluatedCandidate("AAAA", 7)
	assert.Nil(t, err)
	assert.True(t, candidate1.CompareTo(candidate2) < 0, "Incorrect ordering.")
	assert.True(t, candidate2.CompareTo(candidate1) > 0, "Incorrect ordering.")
}

func TestEvaluatedCandidateNegativeFitness(t *testing.T) {
	// should return an error
	candidate, err := NewEvaluatedCandidate("ABC", -1)
	assert.Nil(t, candidate)
	assert.NotNil(t, err)
}
