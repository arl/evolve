package engine

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/aurelien-rainone/evolve/pkg/api"
	"github.com/aurelien-rainone/evolve/pkg/generator"
	"github.com/aurelien-rainone/evolve/pkg/operator"
	"github.com/aurelien-rainone/evolve/pkg/operator/mutation"
	"github.com/aurelien-rainone/evolve/pkg/operator/xover"
	"github.com/aurelien-rainone/evolve/pkg/selection"
	"github.com/aurelien-rainone/evolve/pkg/termination"
	"github.com/stretchr/testify/assert"
)

// Trivial test operator that mutates all integers into zeroes.
type zeroIntMaker struct{}

func (op zeroIntMaker) Apply(selectedCandidates []interface{}, rng *rand.Rand) []interface{} {
	result := make([]interface{}, len(selectedCandidates))
	for i := range selectedCandidates {
		result[i] = 0
	}
	return result
}

type intEvaluator struct{}

func (intEvaluator) Fitness(cand interface{}, pop []interface{}) float64 {
	return float64(cand.(int))
}

func (intEvaluator) IsNatural() bool { return true }

type zeroGenerator struct{}

func (zeroGenerator) GenerateCandidate(rng *rand.Rand) interface{} { return 0 }

func prepareEngine() api.Engine {
	return NewGenerational(
		zeroGenerator{},
		zeroIntMaker{},
		intEvaluator{},
		selection.RouletteWheel,
		rand.New(rand.NewSource(99)))
}

func TestGenerationalEngineElitism(t *testing.T) {
	engine := prepareEngine()

	var avgfitness float64
	// add an observer that record the mean fitness at each generation
	obs := api.ObserverFunc(func(data *api.PopulationData) {
		avgfitness = data.Mean
	})
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
	assert.Equalf(t, 24.0/10.0, avgfitness,
		"elite candidates not preserved correctly: want %v, got %v",
		24.0/10.0, avgfitness)
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

func checkB(b *testing.B, err error) {
	if err != nil {
		b.Fatalf("error: %v", err)
	}
}

// XXX: to prove useful in order to measure the difference between single and
// multithreaded modes, the fitness evaluation must take a `long` time to
// perform the job, otherwise the overhead of concurrent execution hides the
// eventual performance gain.
func benchmarkGenerationalEngine(b *testing.B, multithread bool, strlen int) {
	// Create a factory to generate random 11-character Strings.
	const alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

	// create the target string
	var target string
	for len(target) < strlen {
		target = fmt.Sprintf("%s%c", target, 'A'+byte(rand.Intn(int('Z'-'A'))))
	}

	fac, err := generator.NewString(alphabet, len(target))
	checkB(b, err)

	// 1st operator: string mutation
	mut := mutation.NewString(alphabet)
	checkB(b, mut.SetProb(0.02))

	// 2nd operator: string crossover
	xover := xover.New(xover.StringMater{})

	// Create a pipeline that applies mutation then crossover
	pipe := operator.Pipeline{mut, xover}

	engine := NewGenerational(fac,
		pipe,
		evaluator(target),
		selection.RouletteWheel,
		rand.New(rand.NewSource(99)))

	engine.SetSingleThreaded(!multithread)
	cond := termination.TargetFitness{Fitness: 0, Natural: false}

	b.ResetTimer()
	var best interface{}
	for n := 0; n < b.N; n++ {
		best = engine.Evolve(100000, 5, cond)
	}
	if best.(string) != target {
		b.Errorf("want target string \"%v\", got \"%v\"", target, best.(string))
	}
}

func BenchmarkGenerationalEngineSingleThread10(b *testing.B) {
	benchmarkGenerationalEngine(b, false, 10)
}
func BenchmarkGenerationalEngineMultithread10(b *testing.B) {
	benchmarkGenerationalEngine(b, true, 10)
}
func BenchmarkGenerationalEngineSingleThread100(b *testing.B) {
	benchmarkGenerationalEngine(b, false, 100)
}
func BenchmarkGenerationalEngineMultithread100(b *testing.B) {
	benchmarkGenerationalEngine(b, true, 100)
}
func BenchmarkGenerationalEngineSingleThread1000(b *testing.B) {
	benchmarkGenerationalEngine(b, false, 1000)
}
func BenchmarkGenerationalEngineMultithread1000(b *testing.B) {
	benchmarkGenerationalEngine(b, true, 1000)
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
