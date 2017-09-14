package base

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
