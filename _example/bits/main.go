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
	xover := xover.New[*bitstring.Bitstring](xover.BitstringMater{})
	xover.Probability = generator.Const(0.7)
	xover.Points = generator.Const(1)

	// Define the mutation
	mut := mutation.New[*bitstring.Bitstring](&mutation.Bitstring{
		Probability: generator.Const(0.01),
		FlipCount:   generator.Const(1),
	})

	eval := evolve.EvaluatorFunc(
		true, // natural fitness (higher is better)
		func(cand *bitstring.Bitstring, pop []*bitstring.Bitstring) float64 {
			// our evaluator counts the ones in the bitstring
			return float64(cand.OnesCount())
		})

	epocher := engine.Generational[*bitstring.Bitstring]{
		Op:   operator.Pipeline[*bitstring.Bitstring]{xover, mut},
		Eval: eval,
		Sel:  selection.RouletteWheel[*bitstring.Bitstring]{},
	}

	eng, err := engine.New[*bitstring.Bitstring](generator.Bitstring(nbits), eval, &epocher)
	check(err)

	eng.AddObserver(
		engine.ObserverFunc(func(stats *evolve.PopulationStats[*bitstring.Bitstring]) {
			log.Printf("Generation %d: %s (%v)\n",
				stats.GenNumber,
				stats.BestCand,
				stats.BestFitness)
		}))

	bests, _, err := eng.Evolve(
		100,                                    // 100 candidates in the population
		engine.Elites[*bitstring.Bitstring](2), // best 2 are put into next population
		engine.EndOn[*bitstring.Bitstring](condition.TargetFitness[*bitstring.Bitstring]{
			Fitness: float64(nbits),
			Natural: true,
		}),
	)
	check(err)
	fmt.Println(bests[0])
}
