package islands

import (
	"fmt"
	"testing"

	"github.com/aurelien-rainone/evolve/framework"
	"github.com/stretchr/testify/assert"
)

//
// Utility functions used by unit tests for migration strategies.
//

func createTestPopulation(members ...framework.Candidate) framework.EvaluatedPopulation {
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

func testPopulationContents(t *testing.T, actualPopulation framework.EvaluatedPopulation,
	expectedPopulation ...string) {
	assert.Len(t, actualPopulation, len(expectedPopulation), "wrong population size after migration")
	for i := range actualPopulation {
		got := actualPopulation[i].Candidate()
		want := expectedPopulation[i]
		assert.Equalf(t, got, want, "wrong candidate at index %v, got %v, want %v", i, got, want)
	}
}
