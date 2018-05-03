package evolve

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/aurelien-rainone/evolve/factory"
	"github.com/aurelien-rainone/evolve/framework"
	"github.com/aurelien-rainone/evolve/internal/test"
	"github.com/aurelien-rainone/evolve/pkg/operator"
	"github.com/aurelien-rainone/evolve/pkg/operator/mutation"
	"github.com/aurelien-rainone/evolve/pkg/operator/xover"
	"github.com/aurelien-rainone/evolve/selection"
	"github.com/aurelien-rainone/evolve/termination"
	"github.com/stretchr/testify/assert"
)

func prepareEngine() framework.EvolutionEngine {
	return NewGenerationalEvolutionEngine(
		&factory.AbstractCandidateFactory{
			RandomCandidateGenerator: test.NewStubIntegerFactory(),
		},
		integerZeroMaker{},
		test.IntegerEvaluator{},
		selection.RouletteWheelSelection{},
		rand.New(rand.NewSource(99)))
}

type elitismObserver struct {
	data *framework.PopulationData
}

func (o *elitismObserver) PopulationUpdate(data *framework.PopulationData) {
	o.data = data
}

func (o *elitismObserver) AverageFitness() float64 {
	return o.data.MeanFitness()
}

func TestGenerationalEvolutionEngineElitism(t *testing.T) {
	observer := new(elitismObserver)
	engine := prepareEngine()
	engine.AddEvolutionObserver(observer)
	elite := make([]framework.Candidate, 3)
	// Add the following seed candidates, all better than any others that can possibly
	// get into the population (since every other candidate will always be zero).
	elite[0] = 7 // This candidate should be discarded by elitism.
	elite[1] = 11
	elite[2] = 13
	engine.EvolveWithSeedCandidates(10,
		2, // at least 2 generations because the first is just the initial population.
		elite,
		termination.NewGenerationCount(2))

	// Then when we have run the evolution, if the elite canidates were
	// preserved they will lift the average fitness above zero. The exact value
	// of the expected average fitness is easy to calculate, it is the aggregate
	// fitness divided by the population size.
	assert.Equalf(t, 24.0/10.0, observer.AverageFitness(),
		"Elite candidates not preserved correctly: want %v, got %v",
		24.0/10.0, observer.AverageFitness())
	engine.RemoveEvolutionObserver(observer)
}

func TestGenerationalEvolutionEngineEliteCountTooHigh(t *testing.T) {
	engine := prepareEngine()
	assert.Panics(t, func() { engine.Evolve(10, 10, termination.NewGenerationCount(10)) },
		"Should panic because elite count must be less than the total population size")
}

func TestGenerationalEvolutionEngineEliteCountNegative(t *testing.T) {
	engine := prepareEngine()
	assert.Panics(t, func() { engine.EvolvePopulation(10, -1, termination.NewGenerationCount(10)) },
		"Should panic because elite count must not be negative",
	)
}

func TestGenerationalEvolutionEngineNoTerminationCondition(t *testing.T) {
	engine := prepareEngine()
	assert.Panics(t, func() { engine.Evolve(10, 0) },
		"Should panic because there are no termination conditions")
}

/*
func TestGenerationalEvolutionEngineInterrupt(t*testing.T) {
        final long timeout = 1000L;
        final Thread requestThread = Thread.currentThread();
        engine.addEvolutionObserver(new EvolutionObserver<Integer>()
        {
            public void populationUpdate(PopulationData<? extends Integer> populationData)
            {
                if (populationData.getElapsedTime() > timeout / 2)
                {
                    requestThread.interrupt();
                }
            }
        });
        long startTime = System.currentTimeMillis();
        engine.evolve(10, 0, new ElapsedTime(timeout));
        long elapsedTime = System.currentTimeMillis() - startTime;
        assert Thread.interrupted() : "Thread was not interrupted before timeout.";
        assert elapsedTime < timeout : "Engine did not respond to interrupt before timeout.";
        assert engine.getSatisfiedTerminationConditions().isEmpty()
            : "Interrupted engine should have no satisfied termination conditions.";
    }
*/

func TestGenerationalEvolutionEngineSatisfiedTerminationConditions(t *testing.T) {
	engine := prepareEngine()

	generationsCondition := termination.NewGenerationCount(1)
	engine.Evolve(10, 0, generationsCondition)
	satisfiedConditions, err := engine.SatisfiedTerminationConditions()
	assert.NoError(t, err)
	assert.Len(t, satisfiedConditions, 1)
	assert.Equal(t, generationsCondition, satisfiedConditions[0])
}

func TestGenerationalEvolutionEngineSatisfiedTerminationConditionsBeforeStart(t *testing.T) {
	engine := prepareEngine()

	// Should return an error because evolution hasn't started, let alone terminated.
	conditions, err := engine.SatisfiedTerminationConditions()
	assert.Nil(t, conditions)
	assert.Error(t, err)
}

// Trivial test operator that mutates all integers into zeroes.
type integerZeroMaker struct{}

func (op integerZeroMaker) Apply(selectedCandidates []framework.Candidate, rng *rand.Rand) []framework.Candidate {
	result := make([]framework.Candidate, len(selectedCandidates))
	for i := range selectedCandidates {
		result[i] = 0
	}
	return result
}

func checkB(b *testing.B, err error) {
	if err != nil {
		b.Fatalf("error: %v", err)
	}
}

func BenchmarkGenerationalEvolutionEngine(b *testing.B) {
	const targetString = "HELLO WORLD"

	// Create a factory to generate random 11-character Strings.
	alphabet := make([]byte, 27)
	for c := byte('A'); c <= 'Z'; c++ {
		alphabet[c-'A'] = c
	}
	alphabet[26] = ' '

	stringFactory, err := factory.NewStringFactory(string(alphabet), len(targetString))
	checkB(b, err)

	// 1st operator: string mutation
	mut := mutation.NewStringMutation(string(alphabet))
	checkB(b, mut.SetProb(0.02))

	// 2nd operator: string crossover
	xover := xover.NewCrossover(xover.StringMater{})

	// Create a pipeline that applies mutation then crossover
	pipe := operator.Pipeline{mut, xover}

	fitnessEvaluator := newStringEvaluator(targetString)

	var strategy = &selection.RouletteWheelSelection{}
	rng := rand.New(rand.NewSource(99))

	engine := NewGenerationalEvolutionEngine(stringFactory,
		pipe,
		fitnessEvaluator,
		strategy,
		rng)

	//engine.SetSingleThreaded(true)

	b.ResetTimer()
	var result framework.Candidate
	for n := 0; n < b.N; n++ {
		result = engine.Evolve(100000, 5, termination.NewTargetFitness(0, false))
	}
	fmt.Println(result)

	var conditions []framework.TerminationCondition
	conditions, err = engine.SatisfiedTerminationConditions()
	checkB(b, err)
	for i, condition := range conditions {
		fmt.Printf("satified termination condition %v %T: %v\n",
			i, condition, condition)
	}
}

type stringEvaluator struct{ targetString string }

func newStringEvaluator(targetString string) stringEvaluator {
	return stringEvaluator{targetString: targetString}
}

// Assigns one "fitness point" for every character in the candidate string that
// doesn't match the corresponding position in the target string.
func (se stringEvaluator) Fitness(
	candidate framework.Candidate,
	population []framework.Candidate) float64 {

	var errors float64
	sc := candidate.(string)
	for i := range sc {
		if sc[i] != se.targetString[i] {
			errors++
		}
	}
	return errors
}

func (se stringEvaluator) IsNatural() bool {
	return false
}
