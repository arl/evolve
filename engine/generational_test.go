package engine

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/arl/evolve"
	"github.com/arl/evolve/condition"
	"github.com/arl/evolve/factory"
	"github.com/arl/evolve/generator"
	"github.com/arl/evolve/operator"
	"github.com/arl/evolve/operator/mutation"
	"github.com/arl/evolve/operator/xover"
	"github.com/arl/evolve/selection"
)

// Trivial test operator that mutates all integers into zeroes.
type zeroIntMaker struct{}

func (op zeroIntMaker) Apply(selectedCandidates []int, rng *rand.Rand) []int {
	result := make([]int, len(selectedCandidates))
	for i := range selectedCandidates {
		result[i] = 0
	}
	return result
}

type intEvaluator struct{}

func (intEvaluator) Fitness(cand int, pop []int) float64 {
	return float64(cand)
}

func (intEvaluator) IsNatural() bool { return true }

func check(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Error(err)
	}
}

var zeroFactory = evolve.FactoryFunc[int](func(_ *rand.Rand) int { return 0 })

func TestGenerationalEngineElitism(t *testing.T) {
	epocher := Generational[int]{
		Op:   zeroIntMaker{},
		Eval: intEvaluator{},
		Sel:  selection.RouletteWheel[int]{},
	}

	eng, _ := New[int](zeroFactory, intEvaluator{}, &epocher)

	var avgfitness float64
	// add an observer that record the mean fitness at each generation
	obs := ObserverFunc(func(stats *evolve.PopulationStats[int]) {
		avgfitness = stats.Mean
	})
	eng.AddObserver(obs)

	seeds := make([]int, 3)
	// Add the following seed candidates, all better than any others that can possibly
	// get into the population (since every other candidate will always be zero).
	seeds[0] = 7 // This candidate should be discarded by elitism.
	seeds[1] = 11
	seeds[2] = 13
	cond := condition.GenerationCount[int](3)
	eng.Evolve(10, Elites[int](2), Seeds[int](seeds), EndOn[int](cond))

	// Then when we have run the evolution, if the elite canidates were
	// preserved they will lift the average fitness above zero. The exact value
	// of the expected average fitness is easy to calculate, it is the aggregate
	// fitness divided by the population size.
	assert.Equalf(t, 24.0/10.0, avgfitness,
		"elite candidates not preserved correctly: want %v, got %v",
		24.0/10.0, avgfitness)
	eng.RemoveObserver(obs)
}

func TestGenerationalEngineSatisfiedConditions(t *testing.T) {
	epocher := Generational[int]{
		Op:   zeroIntMaker{},
		Eval: intEvaluator{},
		Sel:  selection.RouletteWheel[int]{},
	}

	eng, _ := New[int](zeroFactory, intEvaluator{}, &epocher)

	cond := condition.GenerationCount[int](1)
	_, satisfied, err := eng.Evolve(10, EndOn[int](cond))
	check(t, err)
	if len(satisfied) != 1 {
		t.Errorf("want len(satisfied) = 1, got %v", len(satisfied))
	}
	if satisfied[0] != cond {
		t.Errorf("want satisfied[0] == cond, got != (cond[0]: %v)", satisfied[0])
	}
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
	const alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

	// create the target string
	var target string
	for len(target) < strlen {
		target = fmt.Sprintf("%s%c", target, 'A'+byte(rand.Intn(int('Z'-'A'))))
	}

	// Create a string generator
	fac, err := factory.NewString(alphabet, len(target))
	checkB(b, err)

	// 1st operator: string mutation
	mut := mutation.New[string](&mutation.String{
		Alphabet:    alphabet,
		Probability: generator.Const(0.02),
	})

	// 2nd operator: string crossover
	xover := xover.New[string](xover.StringMater{})

	// Create a pipeline that applies mutation then crossover
	pipe := operator.Pipeline[string]{mut, xover}

	epocher := Generational[string]{
		Op:   pipe,
		Eval: evaluator(target),
		Sel:  selection.RouletteWheel[string]{},
	}
	eng, err := New[string](fac, evaluator(target), &epocher)
	checkB(b, err)

	// TODO: add option function for singlethread
	// engine.SetSingleThreaded(!multithread)
	cond := condition.TargetFitness[string]{Fitness: 0, Natural: false}

	b.ResetTimer()
	var best interface{}
	for n := 0; n < b.N; n++ {
		best, _, err = eng.Evolve(100000, Elites[string](5), EndOn[string](cond))
		checkB(b, err)
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
// TODO: rename to charMatchEvaluator or something (maybe generalize for byteseq (~string | ~[]byte) , just maybe...)
type evaluator string

func (s evaluator) Fitness(cand string, pop []string) float64 {
	var errors float64
	for i := 0; 0 < len(cand); i++ {
		if cand[i] != string(s)[i] {
			errors++
		}
	}
	return errors
}

// Fitness is not natural, one fitness point represents an error, so the lower
// is better
func (evaluator) IsNatural() bool { return false }
