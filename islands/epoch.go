package islands

import "github.com/aurelien-rainone/evolve/framework"

type epoch struct {
	island                framework.EvolutionEngine
	populationSize        int
	eliteCount            int
	seedCandidates        []framework.Candidate
	terminationConditions []framework.TerminationCondition
}

func newEpoch(
	island framework.EvolutionEngine,
	populationSize int,
	eliteCount int,
	seedCandidates []framework.Candidate,
	terminationConditions ...framework.TerminationCondition) *epoch {

	return &epoch{
		island:                island,
		populationSize:        populationSize,
		eliteCount:            eliteCount,
		seedCandidates:        seedCandidates,
		terminationConditions: terminationConditions,
	}
}

func (e *epoch) Work() interface{} {
	return e.island.EvolvePopulationWithSeedCandidates(e.populationSize, e.eliteCount, e.seedCandidates, e.terminationConditions...)
}
