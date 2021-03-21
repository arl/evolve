package main

import (
	"fmt"
	"log"

	"github.com/arl/evolve"
	"github.com/arl/evolve/condition"
	"github.com/arl/evolve/engine"
	"github.com/arl/evolve/generator"
	"github.com/arl/evolve/operator"
	"github.com/arl/evolve/operator/mutation"
	"github.com/arl/evolve/operator/xover"
	"github.com/arl/evolve/pkg/bitstring"
	"github.com/arl/evolve/selection"
)

const nbits = 20

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// An implementation of the first exercise (page 31) from the book An
// Introduction to Genetic Algorithms, by Melanie Mitchell. The algorithm
// evolves bit strings and the fitness function simply counts the number of ones
// in the bit string. The evolution should therefore converge on strings that
// consist only of ones.
func main() {
	// Define the crossover
	xover := xover.New(xover.BitstringMater{})
	xover.Probability = generator.ConstFloat64(0.7)
	xover.Points = generator.ConstInt(1)
	// define the mutation
	mut := mutation.New(&mutation.Bitstring{
		Probability: generator.ConstFloat64(0.01),
		FlipCount:   generator.ConstInt(1),
	})

	eval := evolve.EvaluatorFunc(
		true, // natural fitness (higher is better)
		func(cand interface{}, pop []interface{}) float64 {
			// our evaluator counts the ones in the bitstring
			return float64(cand.(*bitstring.Bitstring).OnesCount())
		})

	epocher := engine.Generational{
		Op:   operator.Pipeline{xover, mut},
		Eval: eval,
		Sel:  selection.RouletteWheel,
	}

	eng, err := engine.New(generator.Bitstring(nbits), eval, &epocher)
	check(err)

	eng.AddObserver(
		engine.ObserverFunc(func(stats *evolve.PopulationStats) {
			log.Printf("Generation %d: %s (%v)\n",
				stats.GenNumber,
				stats.BestCand,
				stats.BestFitness)
		}))

	bests, _, err := eng.Evolve(
		100,              // 100 candidates in the population
		engine.Elites(2), // best 2 are put into next population
		engine.EndOn(condition.TargetFitness{
			Fitness: float64(nbits),
			Natural: true,
		}),
	)
	check(err)
	fmt.Println(bests[0])
}
