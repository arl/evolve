package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"

	"github.com/arl/evolve"
	"github.com/arl/evolve/condition"
	"github.com/arl/evolve/engine"
	"github.com/arl/evolve/generator"
	"github.com/arl/evolve/operator"
	"github.com/arl/evolve/operator/mutation"
	"github.com/arl/evolve/operator/xover"
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

	evaluator := evolve.EvaluatorFunc[string](false, func(cand string, _ []string) float64 {
		// Assigns one "fitness point" for every character in the candidate
		// string that doesn't match the corresponding position in the
		// target string.
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
	mutation := mutation.New[string](&mutation.String{
		Alphabet:    alphabet,
		Probability: generator.Const(0.02),
	})
	xover := xover.New[string](xover.StringMater{})
	xover.Points = generator.Const(1)
	xover.Probability = generator.Const(1.0)

	// Define a composite evolutionary operator, that is a pipeline that applies
	// to each candidate a string mutation followed by a crossover
	pipeline := operator.Pipeline[string]{mutation, xover}

	// This evaluator assigns one "fitness point" for every character in the
	// The epocher is generational evolutionary engine.
	epocher := engine.Generational[string]{
		Op:   pipeline,
		Eval: evaluator,
		Sel:  selection.RouletteWheel[string]{},
	}

	nchars := len(target)
	// Define the components of our engine
	eng, err := engine.New[string](evolve.FactoryFunc[string](func(rng *rand.Rand) string {
		b := make([]byte, nchars)
		for i := 0; i < nchars; i++ {
			b[i] = alphabet[rng.Int31n(int32(nchars))]
		}
		return string(b)
	}), evaluator, &epocher)
	check(err)

	// Define an observer
	eng.AddObserver(
		engine.ObserverFunc(func(stats *evolve.PopulationStats[string]) {
			log.Printf("Generation %d: %s (%v)\n",
				stats.GenNumber,
				stats.BestCand,
				stats.BestFitness)
		}))

	// Evolution should end when a candidate with a fitness of 0 has been
	// reached (0 different chars between candidate and target string)
	cond := engine.EndOn[string](condition.TargetFitness[string]{Fitness: 0, Natural: false})

	// Start evolution engine and print the best result
	bests, _, err := eng.Evolve(
		100,                      // population zize
		engine.Elites[string](5), // number of elites, if any
		cond,
	)
	check(err)
	log.Println(bests[0])
}
