package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"

	"github.com/arl/evolve"
	"github.com/arl/evolve/condition"
	"github.com/arl/evolve/crossover"
	"github.com/arl/evolve/engine"
	"github.com/arl/evolve/generator"
	"github.com/arl/evolve/mutation"
	"github.com/arl/evolve/selection"
)

const alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ "

var target = "EVOLVE WORLD"

func check(err error) {
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

func main() {
	if len(os.Args) == 2 {
		target = strings.ToUpper(os.Args[1])
	}

	// Setup a generator of random strings
	for _, c := range target {
		if !strings.ContainsRune(alphabet, c) {
			log.Fatalf("Target string must only be made of runes in %q", alphabet)
		}
	}

	// Our factory generates random strings that have the same length as the
	// target string.
	factory := func(rng *rand.Rand) string {
		b := make([]byte, len(target))
		for i := 0; i < len(target); i++ {
			b[i] = alphabet[rng.Intn(len(target))]
		}
		return string(b)
	}

	// Our candidate evaluator assigns one "fitness point" for every character
	// in the candidate string that doesn't match the corresponding position in
	// the target string.
	evaluator := evolve.EvaluatorFunc[string](false, func(cand string) float64 {
		var n int
		for i := range cand {
			if cand[i] != target[i] {
				n++
			}
		}
		return float64(n)
	})

	// Define our evolutionary operators, a string mutation where each rune has
	// a probability of mutation of 0.02, plus a default string crossover.
	mutation := evolve.NewMutation[string](&mutation.String{
		Alphabet:    alphabet,
		Probability: generator.Const(0.02),
	})
	xover := evolve.NewCrossover[string](crossover.StringMater{})
	xover.Points = generator.Const(1)
	xover.Probability = generator.Const(1.0)

	// Define a composite evolutionary operator, that is a pipeline that applies
	// to each candidate a string mutation followed by a crossover
	pipeline := evolve.Pipeline[string]{mutation, xover}

	generational := &engine.Generational[string]{
		Operator:  pipeline,
		Evaluator: evaluator,
		Selection: selection.RouletteWheel[string]{},
		NumElites: 5,
	}

	// Define the components of our engine
	eng := engine.Engine[string]{
		Factory:   evolve.FactoryFunc[string](factory),
		Evaluator: evaluator,
		Epocher:   generational,
		EndConditions: []evolve.Condition[string]{
			// Evolution terminates when a candidate reach fitness 0 (0 chars
			// are different from the target string).
			condition.TargetFitness[string]{Fitness: 0, Natural: false},
		},
		Concurrency: 1,
	}

	// Define an observer
	eng.AddObserver(
		engine.ObserverFunc(func(stats *evolve.PopulationStats[string]) {
			log.Printf("Generation %d: %s (%v)\n", stats.Generation, stats.Best, stats.BestFitness)
		}))

	// Start evolution engine and print the best result
	bests, _, err := eng.Evolve(100)
	check(err)
	log.Println(bests.Candidates[0], bests.Fitness[0])
}
