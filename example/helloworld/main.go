package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/arl/evolve/pkg/api"
	"github.com/arl/evolve/pkg/engine"
	"github.com/arl/evolve/pkg/generator"
	"github.com/arl/evolve/pkg/operator"
	"github.com/arl/evolve/pkg/operator/mutation"
	"github.com/arl/evolve/pkg/operator/xover"
	"github.com/arl/evolve/pkg/selection"
	"github.com/arl/evolve/pkg/termination"
)

const alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ "

var target = "EVOLVE WORLD"

func check(err error) {
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

// This evaluator assigns one "fitness point" for every character in the
// candidate string that doesn't match the corresponding position in the
// target string.
type evaluator struct{}

func (evaluator) Fitness(cand interface{}, pop []interface{}) float64 {
	// count differences between candidate and target strings
	var nerrors int
	sc := cand.(string)
	for i := range sc {
		if sc[i] != target[i] {
			nerrors++
		}
	}
	return float64(nerrors)
}

// Non natural fitness, lower is better
func (evaluator) IsNatural() bool { return false }

func main() {
	if len(os.Args) == 2 {
		target = strings.ToUpper(os.Args[1])
	}

	// Setup a generator of random strings
	for _, c := range target {
		if !strings.ContainsRune(alphabet, c) {
			log.Fatalf("Target string must be solely made of \"%v\"", alphabet)
		}
	}
	gen, err := generator.NewString(alphabet, len(target))
	check(err)

	// Define our evolutionary operators, a string mutation where each rune has
	// a probability of mutation of 0.02, plus a default string crossover.
	mutation := mutation.NewString(alphabet)
	check(mutation.SetProb(0.02))
	xover := xover.New(xover.StringMater{})

	// Define a composite evolutionary operator, that is a pipeline that applies
	// to each candidate a string mutation followed by a crossover
	pipeline := operator.Pipeline{mutation, xover}

	// This evaluator assigns one "fitness point" for every character in the
	// The epocher is generational evolutionary engine.
	epocher := engine.Generational{
		Op:   pipeline,
		Eval: evaluator{},
		Sel:  selection.RouletteWheel,
	}

	// Define the components of our engine
	eng, err := engine.New(gen, evaluator{}, &epocher)
	check(err)

	// Define an observer
	eng.AddObserver(
		api.ObserverFunc(func(data *api.PopulationData) {
			log.Printf("Generation %d: %s (%v)\n",
				data.GenNumber,
				data.BestCand,
				data.BestFitness)
		}))

	// Evolution should end when a candidate with a fitness of 0 has been
	// reached (0 different chars between candidate and target string)
	cond := termination.TargetFitness{Fitness: 0, Natural: false}

	// Start evolution engine and print the best result
	bests, _, err := eng.Evolve(100, engine.Elites(5), engine.EndOn(cond))
	check(err)
	log.Println(bests[0])
}
