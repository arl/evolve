package test

import (
	"math/rand"
	"testing"

	"github.com/aurelien-rainone/evolve/pkg/api"
	"github.com/aurelien-rainone/evolve/pkg/factory"
)

// Trivial fitness evaluator for integers. Used by unit tests.
type IntEvaluator struct{}

func (IntEvaluator) Fitness(cand interface{}, pop []interface{}) float64 {
	return float64(cand.(int))
}

func (IntEvaluator) IsNatural() bool { return true }

// Stub candidate factory for tests. Always returns zero-valued integers.
var ZeroIntFactory = factory.BaseFactory{CandidateGenerator: ZeroIntGenerator{}}

type ZeroIntGenerator struct{}

func (ZeroIntGenerator) GenerateCandidate(rng *rand.Rand) interface{} { return 0 }

// IntAdjuster is a trivial test operator that mutates all integers by
// adding a fixed offset.
type IntAdjuster int

func (op IntAdjuster) Apply(cands []interface{}, rng *rand.Rand) []interface{} {
	result := make([]interface{}, len(cands))
	for i, c := range cands {
		result[i] = c.(int) + int(op)
	}
	return result
}

//
// Utility functions used by unit tests for migration strategies.
//

func CreateTestPopulation(members ...interface{}) api.Population {
	pop := make(api.Population, len(members))
	for i, member := range members {
		pop[i] = &api.Individual{Candidate: member, Fitness: 0}
	}
	return pop
}

func AssertPopulationContents(t *testing.T, actualpop api.Population, expected ...string) {
	t.Helper() // mark current function as helper in case of error
	if len(actualpop) != len(expected) {
		t.Errorf("wrong population size, want %v got %v", len(expected), len(actualpop))
	}

	for i, cand := range actualpop {
		got := cand.Candidate.(string)
		want := expected[i]
		if want != got {
			t.Errorf("wrong candidate at index %v, want %v, got %v", i, want, got)
		}
	}
}
