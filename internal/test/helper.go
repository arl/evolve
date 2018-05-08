package test

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/aurelien-rainone/evolve/pkg/api"
	"github.com/aurelien-rainone/evolve/pkg/factory"
	"github.com/stretchr/testify/assert"
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

func CreateTestPopulation(members ...interface{}) api.EvaluatedPopulation {
	var err error
	pop := make(api.EvaluatedPopulation, len(members))
	for i, member := range members {
		pop[i], err = api.NewEvaluatedCandidate(member, 0)
		if err != nil {
			panic(fmt.Sprintf("can't create test pop: %v", err))
		}
	}
	return pop
}

func AssertPopulationContents(t *testing.T, actualpop api.EvaluatedPopulation, expected ...string) {
	t.Helper() // mark current function as helper in case of error
	assert.Len(t, actualpop, len(expected), "wrong population size after migration")
	for i, cand := range actualpop {
		got := cand.Candidate()
		want := expected[i]
		assert.Equalf(t, want, got, "wrong candidate at index %v, want %v, got %v", i, want, got)
	}
}
