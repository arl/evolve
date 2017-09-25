package test

import (
	"math/rand"

	"github.com/aurelien-rainone/evolve/framework"
)

// Trivial fitness evaluator for integers. Used by unit tests.
type IntegerEvaluator struct{}

func (e IntegerEvaluator) Fitness(candidate framework.Candidate, population []framework.Candidate) float64 {
	return float64(candidate.(int))
}

func (e IntegerEvaluator) IsNatural() bool {
	return true
}

// Stub candidate factory for tests.  Always returns zero-valued integers.
type StubIntegerFactory struct{}

func (f StubIntegerFactory) GenerateRandomCandidate(rng *rand.Rand) framework.Candidate {
	return 0
}
