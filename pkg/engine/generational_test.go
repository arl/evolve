package engine

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/aurelien-rainone/evolve/internal/test"
	"github.com/aurelien-rainone/evolve/pkg/api"
	"github.com/aurelien-rainone/evolve/pkg/factory"
	"github.com/aurelien-rainone/evolve/pkg/operator"
	"github.com/aurelien-rainone/evolve/pkg/operator/mutation"
	"github.com/aurelien-rainone/evolve/pkg/operator/xover"
	"github.com/aurelien-rainone/evolve/pkg/selection"
	"github.com/aurelien-rainone/evolve/pkg/termination"
	"github.com/stretchr/testify/assert"
)

func prepareEngine() api.Engine {
	return NewGenerational(
		&test.ZeroIntFactory,
		zeroIntMaker{},
		test.IntEvaluator{},
		selection.RouletteWheel,
		rand.New(rand.NewSource(99)))
}

type elitismObserver api.PopulationData

func (o *elitismObserver) PopulationUpdate(data *api.PopulationData) { *o = elitismObserver(*data) }

func (o *elitismObserver) AverageFitness() float64 { return o.Mean }

func TestGenerationalEngineElitism(t *testing.T) {
	obs := new(elitismObserver)
	engine := prepareEngine()
	engine.AddObserver(obs)
	elite := make([]interface{}, 3)
	// Add the following seed candidates, all better than any others that can possibly
	// get into the population (since every other candidate will always be zero).
	elite[0] = 7 // This candidate should be discarded by elitism.
	elite[1] = 11
	elite[2] = 13
	engine.EvolveWithSeedCandidates(10,
		2, // at least 2 generations because the first is just the initial population.
		elite,
		termination.GenerationCount(2))

	// Then when we have run the evolution, if the elite canidates were
	// preserved they will lift the average fitness above zero. The exact value
	// of the expected average fitness is easy to calculate, it is the aggregate
	// fitness divided by the population size.
	assert.Equalf(t, 24.0/10.0, obs.AverageFitness(),
		"elite candidates not preserved correctly: want %v, got %v",
		24.0/10.0, obs.AverageFitness())
	engine.RemoveObserver(obs)
}

func TestGenerationalEngineEliteCountTooHigh(t *testing.T) {
	engine := prepareEngine()
	assert.Panics(t, func() {
		engine.Evolve(10, 10, termination.GenerationCount(10))
	}, "elite count must be less than the total population size")
}

func TestGenerationalEngineNoTerminationCondition(t *testing.T) {
	engine := prepareEngine()
	assert.Panics(t, func() {
		engine.Evolve(10, 0)
	}, "some termination conditions must be set")
}

/*
func TestGenerationalEngineInterrupt(t*testing.T) {
        final long timeout = 1000L;
        final Thread requestThread = Thread.currentThread();
        engine.addObserver(new Observer<Integer>()
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

func TestGenerationalEngineSatisfiedTerminationConditions(t *testing.T) {
	engine := prepareEngine()

	cond := termination.GenerationCount(1)
	engine.Evolve(10, 0, cond)
	satisfied, err := engine.SatisfiedTerminationConditions()
	assert.NoError(t, err)
	assert.Len(t, satisfied, 1)
	assert.Equal(t, cond, satisfied[0])
}

func TestGenerationalEngineSatisfiedTerminationConditionsBeforeStart(t *testing.T) {
	engine := prepareEngine()

	// Should return an error because evolution hasn't started, let alone terminated.
	satisfied, err := engine.SatisfiedTerminationConditions()
	assert.Nil(t, satisfied)
	assert.Error(t, err)
}

// Trivial test operator that mutates all integers into zeroes.
type zeroIntMaker struct{}

func (op zeroIntMaker) Apply(selectedCandidates []interface{}, rng *rand.Rand) []interface{} {
	result := make([]interface{}, len(selectedCandidates))
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

func BenchmarkGenerationalEngine(b *testing.B) {
	const targetString = "HELLO WORLD"

	// Create a factory to generate random 11-character Strings.
	alphabet := make([]byte, 27)
	for c := byte('A'); c <= 'Z'; c++ {
		alphabet[c-'A'] = c
	}
	alphabet[26] = ' '

	fac, err := factory.NewString(string(alphabet), len(targetString))
	checkB(b, err)

	// 1st operator: string mutation
	mut := mutation.NewString(string(alphabet))
	checkB(b, mut.SetProb(0.02))

	// 2nd operator: string crossover
	xover := xover.New(xover.StringMater{})

	// Create a pipeline that applies mutation then crossover
	pipe := operator.Pipeline{mut, xover}

	engine := NewGenerational(fac,
		pipe,
		evaluator(targetString),
		selection.RouletteWheel,
		rand.New(rand.NewSource(99)))

	//engine.SetSingleThreaded(true)

	b.ResetTimer()
	var best interface{}
	for n := 0; n < b.N; n++ {
		best = engine.Evolve(100000, 5, termination.TargetFitness{Fitness: 0, Natural: false})
	}
	fmt.Println(best)

	var satisfied []api.TerminationCondition
	satisfied, err = engine.SatisfiedTerminationConditions()
	checkB(b, err)
	for i, cond := range satisfied {
		fmt.Printf("satified termination condition %v: %v\n", i, cond)
	}
}

// This 'evaluator' assigns one "fitness point" for every character in the
// candidate string that doesn't match the corresponding position in the target
// string.
type evaluator string

func (s evaluator) Fitness(
	cand interface{},
	pop []interface{}) float64 {

	var errors float64
	sc := cand.(string)
	for i := range sc {
		if sc[i] != string(s)[i] {
			errors++
		}
	}
	return errors
}

// Fitness is not natural, one fitness point represents an error, so the lower
// is better
func (evaluator) IsNatural() bool { return false }
