package islands

import "github.com/aurelien-rainone/evolve/pkg/api"

type epoch struct {
	island                api.EvolutionEngine
	populationSize        int
	eliteCount            int
	seedCandidates        []api.Candidate
	terminationConditions []api.TerminationCondition
}

func newEpoch(
	island api.EvolutionEngine,
	populationSize int,
	eliteCount int,
	seedCandidates []api.Candidate,
	terminationConditions ...api.TerminationCondition) *epoch {

	return &epoch{
		island:                island,
		populationSize:        populationSize,
		eliteCount:            eliteCount,
		seedCandidates:        seedCandidates,
		terminationConditions: terminationConditions,
	}
}

func (e *epoch) Work() (interface{}, error) {
	return e.island.EvolvePopulationWithSeedCandidates(
		e.populationSize,
		e.eliteCount,
		e.seedCandidates,
		e.terminationConditions...), nil
}
