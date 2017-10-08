package test

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/aurelien-rainone/evolve/factory"
	"github.com/aurelien-rainone/evolve/framework"
	"github.com/stretchr/testify/assert"
)

// Trivial fitness evaluator for integers. Used by unit tests.
type IntegerEvaluator struct{}

func (e IntegerEvaluator) Fitness(candidate framework.Candidate, population []framework.Candidate) float64 {
	return float64(candidate.(int))
}

func (e IntegerEvaluator) IsNatural() bool {
	return true
}

// Stub candidate factory for tests. Always returns zero-valued integers.
type StubIntegerFactory struct {
	factory.AbstractCandidateFactory
}

func NewStubIntegerFactory() *StubIntegerFactory {
	return &StubIntegerFactory{
		factory.AbstractCandidateFactory{
			RandomCandidateGenerator: ZeroIntegerGenerator{},
		},
	}
}

type ZeroIntegerGenerator struct{}

func (zig ZeroIntegerGenerator) GenerateRandomCandidate(rng *rand.Rand) framework.Candidate {
	return 0
}

// IntegerAdjuster is a trivial test operator that mutates all integers by
// adding a fixed offset.
type IntegerAdjuster int

func (op IntegerAdjuster) Apply(selectedCandidates []framework.Candidate, rng *rand.Rand) []framework.Candidate {
	result := make([]framework.Candidate, len(selectedCandidates))
	for i, c := range selectedCandidates {
		result[i] = c.(int) + int(op)
	}
	return result
}

//
// Utility functions used by unit tests for migration strategies.
//

func CreateTestPopulation(members ...framework.Candidate) framework.EvaluatedPopulation {
	var err error
	population := make(framework.EvaluatedPopulation, len(members))
	for i, member := range members {
		population[i], err = framework.NewEvaluatedCandidate(member, 0)
		if err != nil {
			panic(fmt.Sprintf("can't create test population: %v", err))
		}
	}
	return population
}

func AssertPopulationContents(t *testing.T, actualPopulation framework.EvaluatedPopulation,
	expectedPopulation ...string) {
	assert.Len(t, actualPopulation, len(expectedPopulation), "wrong population size after migration")
	for i := range actualPopulation {
		got := actualPopulation[i].Candidate()
		want := expectedPopulation[i]
		assert.Equalf(t, want, got, "wrong candidate at index %v, want %v, got %v", i, want, got)
	}
}
