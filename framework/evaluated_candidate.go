package framework

import "fmt"

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
