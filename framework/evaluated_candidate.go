package framework

import (
	"fmt"
	"math"
)

// EvaluatedCandidate is an immutable wrapper for associating a candidate
// solution with its fitness score.
type EvaluatedCandidate struct {
	candidate Candidate
	fitness   float64
}

// NewEvaluatedCandidate returns an EvaluatedCandidate
func NewEvaluatedCandidate(candidate Candidate, fitness float64) (*EvaluatedCandidate, error) {
	if fitness < 0 {
		return nil, fmt.Errorf("fitness score must be >= 0, got %v", fitness)
	}
	return &EvaluatedCandidate{
		candidate: candidate,
		fitness:   fitness,
	}, nil
}

// Candidate returns the evolved candidate solution.
func (ec *EvaluatedCandidate) Candidate() Candidate {
	return ec.candidate
}

// Fitness returns the fitness score for the associated candidate.
func (ec *EvaluatedCandidate) Fitness() float64 {
	return ec.fitness
}

// Equals returns true If this object is logically equivalent to {code o}.
func (ec *EvaluatedCandidate) Equals(o *EvaluatedCandidate) bool {
	if ec == o {
		return true
	} else if o == nil {
		return false
	}
	return ec.Fitness() == o.Fitness()
}

// Hash returns the candidate's hash code (consistent with Equals)
func (ec *EvaluatedCandidate) Hash() int64 {
	var temp uint64
	if ec.fitness != 0 {
		temp = math.Float64bits(ec.fitness)
	}
	return int64(temp ^ (temp >> 32))
}

// CompareTo compares this candidate's fitness score with that of the specified
// candidate.
//
// Returns -1, 0 or 1 if this candidate's score is less than, equal to, or
// greater than that of the specified candidate. The comparison applies to the
// raw numerical score and does not consider whether that score is a natural
// fitness score or not.
func (ec *EvaluatedCandidate) CompareTo(o *EvaluatedCandidate) int {
	switch {
	case ec.Fitness() < o.Fitness():
		return -1
	case ec.Fitness() > o.Fitness():
		return 1
	}
	return 0
}

// EvaluatedPopulation represents a slice of pointers to EvaluatedCandidate.
type EvaluatedPopulation []*EvaluatedCandidate

// Len is the number of elements in the collection.
func (s EvaluatedPopulation) Len() int {
	return len(s)
}

// Less reports whether the element with
// index a should sort before the element with index b.
func (s EvaluatedPopulation) Less(i, j int) bool {
	return s[i].CompareTo(s[j]) == -1
}

// Swap swaps the elements with indexes i and j.
func (s EvaluatedPopulation) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
