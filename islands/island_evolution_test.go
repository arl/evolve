package islands

import (
	"math/rand"
	"testing"

	"github.com/aurelien-rainone/evolve/internal/test"
	"github.com/aurelien-rainone/evolve/pkg/api"
	"github.com/aurelien-rainone/evolve/pkg/selection"
	"github.com/aurelien-rainone/evolve/termination"
	"github.com/stretchr/testify/assert"
)

type countObserver struct {
	observedEpochCount       int
	observedGenerationCounts []int
}

func (obs *countObserver) PopulationUpdate(data *api.PopulationData) {
	obs.observedEpochCount++
}

func (obs *countObserver) IslandPopulationUpdate(islandIndex int, data *api.PopulationData) {
	obs.observedGenerationCounts[islandIndex]++
}

func TestIslandEvolutionObservers(t *testing.T) {
	// This test makes sure that the evolution observer global method only gets
	// invoked at the end of each epoch, and that the island method gets invoked
	// for each generation on each island.
	var (
		islandCount     = 3
		epochCount      = 2
		generationCount = 5
	)

	islandEvolution := NewIslandEvolution(islandCount,
		RingMigration{},
		test.NewStubIntegerFactory(),
		test.IntegerAdjuster(2),
		dummyFitnessEvaluator{},
		selection.RouletteWheelSelection,
		rand.New(rand.NewSource(99)),
	)

	obs := countObserver{
		observedGenerationCounts: make([]int, islandCount),
	}
	islandEvolution.AddEvolutionObserver(&obs)

	islandEvolution.Evolve(5, 0, 5, 0, termination.NewGenerationCount(2))
	assert.Equal(t, obs.observedEpochCount, 2, "want observer to have been notified twice, got ", obs.observedEpochCount)
	for i := 0; i < islandCount; i++ {
		expected := epochCount * generationCount
		assert.Equalf(t,
			obs.observedGenerationCounts[i],
			expected, "want generation count for island %v = %v, got %v",
			i, expected, obs.observedGenerationCounts[i])
	}
}

/*
   @Test
   public void testInterrupt()
   {
       IslandEvolution<Integer> islandEvolution = new IslandEvolution<Integer>(2,
                                                                               new RingMigration(),
                                                                               new StubIntegerFactory(),
                                                                               new IntegerAdjuster(2),
                                                                               new DummyFitnessEvaluator(),
                                                                               new RouletteWheelSelection(),
                                                                               apiTestUtils.getRNG());
       final long timeout = 1000L;
       final Thread requestThread = Thread.currentThread();
       islandEvolution.addEvolutionObserver(new IslandEvolutionObserver<Integer>()
       {
           public void populationUpdate(PopulationData<? extends Integer> populationData)
           {
               if (populationData.getElapsedTime() > timeout / 2)
               {
                   requestThread.interrupt();
               }
           }


           public void islandPopulationUpdate(int islandIndex, PopulationData<? extends Integer> populationData){}
       });
       long startTime = System.currentTimeMillis();
       islandEvolution.evolve(10, 0, 10, 0, new ElapsedTime(timeout));
       long elapsedTime = System.currentTimeMillis() - startTime;
       assert Thread.interrupted() : "Thread was not interrupted before timeout.";
       assert elapsedTime < timeout : "Engine did not respond to interrupt before timeout.";
       assert islandEvolution.getSatisfiedTerminationConditions().isEmpty()
           : "Interrupted islands should have no satisfied termination conditions.";
   }

*/

func TestGetSatisfiedTerminationConditionsBeforeStart(t *testing.T) {
	islandEvolution := NewIslandEvolution(3,
		RingMigration{},
		test.NewStubIntegerFactory(),
		test.IntegerAdjuster(2),
		dummyFitnessEvaluator{},
		selection.RouletteWheelSelection,
		rand.New(rand.NewSource(99)),
	)
	// Should return an error and nil as the slice of satisfied termination
	// conditions because evolution hasn't started, let alone terminated.
	cond, err := islandEvolution.SatisfiedTerminationConditions()
	assert.Nil(t, cond, "want nil for satisfied termination conditions, evolution hasn't started, got:", cond)
	assert.Error(t, err, "want error when retrieving satisfied termination conditions, evolution hasn't started, got nil")
}

type dummyFitnessEvaluator struct{}

func (dfe dummyFitnessEvaluator) Fitness(candidate api.Candidate, population []api.Candidate) float64 {
	return 0
}

func (dfe dummyFitnessEvaluator) IsNatural() bool {
	return true
}
