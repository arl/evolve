package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/aurelien-rainone/evolve/pkg/api"
	"github.com/aurelien-rainone/evolve/pkg/engine"
	"github.com/aurelien-rainone/evolve/pkg/generator"
	"github.com/aurelien-rainone/evolve/pkg/operator"
	"github.com/aurelien-rainone/evolve/pkg/operator/mutation"
	"github.com/aurelien-rainone/evolve/pkg/operator/xover"
	"github.com/aurelien-rainone/evolve/pkg/selection"
	"github.com/aurelien-rainone/evolve/pkg/termination"
)

// This 'evaluator' assigns one "fitness point" for every character in the
// candidate string that doesn't match the corresponding position in the target
// string.
type evaluator string

func (s evaluator) Fitness(cand interface{}, pop []interface{}) float64 {
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

var alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ "

// Create a generator that creates random string
func createGenerator(target string) *generator.String {
	for _, c := range target {
		if !strings.ContainsRune(alphabet, c) {
			fmt.Printf("All runes must be exist in the alphabet ('%v'), that's not the case of %c\n", alphabet, c)
			os.Exit(1)
		}
	}

	fac, err := generator.NewString(alphabet, len(target))
	check(err)
	return fac
}

func main() {
	var target = "HELLO WORLD"
	if len(os.Args) == 2 {
		target = strings.ToUpper(os.Args[1])
	}

	// create the generator that will generate random candidates
	fac := createGenerator(target)

	// create an evolutionary operator pipeline that will apply to each
	// candidate, first a string mutation and then a crossover
	mutation := mutation.NewString(alphabet)
	check(mutation.SetProb(0.02))
	xover := xover.New(xover.StringMater{})
	pipeline := operator.Pipeline{mutation, xover}

	// the evaluator will assign fitness to the candidates
	eval := evaluator(target)

	// choose a selection strategy
	var selector = selection.RouletteWheel
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	// we can now define our evolutionary engine
	engine := engine.NewGenerational(fac, pipeline, eval, selector, rng)

	// define an observer
	engine.AddObserver(
		api.ObserverFunc(func(data *api.PopulationData) {
			fmt.Printf("Generation %d: %s (%v)\n",
				data.GenNumber,
				data.BestCand,
				data.BestFitness)
		}))

	// we want evolution to end when a fitness of 0 has been reached (0
	// differences between candidate and target string)
	condition := termination.TargetFitness{Fitness: 0, Natural: false}

	// start evolution engine and print the best result
	fmt.Println(engine.Evolve(100, 5, condition))
}

func check(err error) {
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
