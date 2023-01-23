package engine

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/arl/evolve"
	"github.com/arl/evolve/condition"
	"github.com/arl/evolve/factory"
	"github.com/arl/evolve/generator"
	"github.com/arl/evolve/operator/mutation"
	"github.com/arl/evolve/operator/xover"
	"github.com/arl/evolve/selection"
)

// Dummy operator acting on populations of numbers by setting all candidates to 0.
type zeroMaker struct{}

func (op zeroMaker) Apply(pop *evolve.Population[int], rng *rand.Rand) {
	for i := range pop.Candidates {
		pop.Candidates[i] = 0
	}
}

type intEvaluator struct{}

func (intEvaluator) Fitness(cand int) float64 {
	return float64(cand)
}

func (intEvaluator) IsNatural() bool { return true }

var zeroFactory = evolve.FactoryFunc[int](func(_ *rand.Rand) int { return 0 })

func TestGenerationalEngineElitism(t *testing.T) {
	eng := Engine[int]{
		Factory:   zeroFactory,
		Evaluator: intEvaluator{},
		Epocher: &Generational[int]{
			Operator:  zeroMaker{},
			Evaluator: intEvaluator{},
			Selection: selection.RouletteWheel[int]{},
			NumElites: 2,
		},
		// Seed candidates, all better than any others that can possibly get
		// into the population (since every other candidate will always be
		// zero). Though elitism should discard 7.
		Seeds: []int{7, 11, 13},
		EndConditions: []evolve.Condition[int]{
			condition.GenerationCount[int](3),
		},
	}

	// Add an observer that records the mean fitness at each generation.
	var avgfitness float64
	eng.AddObserver(ObserverFunc(func(stats *evolve.PopulationStats[int]) {
		avgfitness = stats.Mean
	}))

	if _, _, err := eng.Evolve(10); err != nil {
		t.Fatal(err)
	}

	// Run the evolution. We verify that the 2 best candidates have been kept by
	// using the expected average fitness.
	const want = (11 + 13) / 10.
	if avgfitness != want {
		t.Errorf("average fitness = %v, want %v", avgfitness, want)
	}
}

func TestGenerationalEngineSatisfiedConditions(t *testing.T) {
	eng := Engine[int]{
		Factory:   zeroFactory,
		Evaluator: intEvaluator{},
		Epocher: &Generational[int]{
			Operator:  zeroMaker{},
			Evaluator: intEvaluator{},
			Selection: selection.RouletteWheel[int]{},
		},
		EndConditions: []evolve.Condition[int]{
			condition.GenerationCount[int](1),
		},
	}

	_, satisfied, err := eng.Evolve(10)
	if err != nil {
		t.Fatal(err)
	}

	if len(satisfied) != 1 {
		t.Errorf("len(satisfied) = %d, want 1", len(satisfied))
	}
}

func checkB(b *testing.B, err error) {
	if err != nil {
		b.Fatalf("error: %v", err)
	}
}

// XXX: To be proven useful in order to measure the difference between single and
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
	factory, err := factory.NewString(alphabet, len(target))
	checkB(b, err)

	eng := Engine[string]{
		Factory:   factory,
		Evaluator: evaluator(target),
		Epocher: &Generational[string]{
			// Create a operator pipeline that first apply a string muration then a crossover.
			Operator: evolve.Pipeline[string]{
				evolve.NewMutation[string](&mutation.String{
					Alphabet:    alphabet,
					Probability: generator.Const(0.02),
				}),
				evolve.NewCrossover[string](xover.StringMater{}),
			},
			Evaluator: evaluator(target),
			Selection: selection.RouletteWheel[string]{},
			NumElites: 5,
		},
		EndConditions: []evolve.Condition[string]{
			condition.TargetFitness[string]{Fitness: 0, Natural: false},
		},
	}

	b.ResetTimer()
	var best interface{}
	for n := 0; n < b.N; n++ {
		best, _, err = eng.Evolve(100000)
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

func (s evaluator) Fitness(cand string) float64 {
	var errors float64
	for i := 0; i < len(cand); i++ {
		if cand[i] != string(s)[i] {
			errors++
		}
	}
	return errors
}

// Fitness is not natural, one fitness point represents an error, so the lower
// is better
func (evaluator) IsNatural() bool { return false }
