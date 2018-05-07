package test

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/aurelien-rainone/evolve/factory"
	"github.com/aurelien-rainone/evolve/pkg/api"
	"github.com/stretchr/testify/assert"
)

// Trivial fitness evaluator for integers. Used by unit tests.
type IntegerEvaluator struct{}

func (e IntegerEvaluator) Fitness(candidate api.Candidate, population []api.Candidate) float64 {
	return float64(candidate.(int))
}

func (e IntegerEvaluator) IsNatural() bool {
	return true
}

// Stub candidate factory for tests. Always returns zero-valued integers.
type StubIntegerFactory struct {
	factory.BaseFactory
}

func NewStubIntegerFactory() *StubIntegerFactory {
	return &StubIntegerFactory{
		factory.BaseFactory{
			CandidateGenerator: ZeroIntegerGenerator{},
		},
	}
}

type ZeroIntegerGenerator struct{}

func (zig ZeroIntegerGenerator) GenerateCandidate(rng *rand.Rand) api.Candidate {
	return 0
}

// IntegerAdjuster is a trivial test operator that mutates all integers by
// adding a fixed offset.
type IntegerAdjuster int

func (op IntegerAdjuster) Apply(selectedCandidates []api.Candidate, rng *rand.Rand) []api.Candidate {
	result := make([]api.Candidate, len(selectedCandidates))
	for i, c := range selectedCandidates {
		result[i] = c.(int) + int(op)
	}
	return result
}

//
// Utility functions used by unit tests for migration strategies.
//

func CreateTestPopulation(members ...api.Candidate) api.EvaluatedPopulation {
	var err error
	population := make(api.EvaluatedPopulation, len(members))
	for i, member := range members {
		population[i], err = api.NewEvaluatedCandidate(member, 0)
		if err != nil {
			panic(fmt.Sprintf("can't create test population: %v", err))
		}
	}
	return population
}

func AssertPopulationContents(t *testing.T, actualPopulation api.EvaluatedPopulation,
	expectedPopulation ...string) {
	t.Helper() // mark current function as helper in case of error
	assert.Len(t, actualPopulation, len(expectedPopulation), "wrong population size after migration")
	for i, evCand := range actualPopulation {
		got := evCand.Candidate()
		want := expectedPopulation[i]
		assert.Equalf(t, want, got, "wrong candidate at index %v, want %v, got %v", i, want, got)
	}
}
