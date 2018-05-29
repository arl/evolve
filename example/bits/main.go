package main

import (
	"fmt"
	"log"

	"github.com/aurelien-rainone/evolve/pkg/api"
	"github.com/aurelien-rainone/evolve/pkg/bitstring"
	"github.com/aurelien-rainone/evolve/pkg/engine"
	"github.com/aurelien-rainone/evolve/pkg/generator"
	"github.com/aurelien-rainone/evolve/pkg/operator"
	"github.com/aurelien-rainone/evolve/pkg/operator/mutation"
	"github.com/aurelien-rainone/evolve/pkg/operator/xover"
	"github.com/aurelien-rainone/evolve/pkg/selection"
	"github.com/aurelien-rainone/evolve/pkg/termination"
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
	// define the crossover
	xover := xover.New(xover.BitstringMater{})
	check(xover.SetPoints(1))
	check(xover.SetProb(0.7))

	// define the mutation
	mut := mutation.NewBitstring()
	check(mut.SetProb(0.01))

	eval := api.EvaluatorFunc(true, // natural fitness (higher is better)
		// our evaluator counts the ones in the bitstring
		func(cand interface{}, pop []interface{}) float64 {
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
		api.ObserverFunc(func(data *api.PopulationData) {
			log.Printf("Generation %d: %s (%v)\n",
				data.GenNumber,
				data.BestCand,
				data.BestFitness)
		}))

	bests, _, err := eng.Evolve(
		100,              // 100 candidates in the population
		engine.Elites(2), // best 2 are put into next population
		engine.EndOn(termination.TargetFitness{
			Fitness: float64(nbits),
			Natural: true,
		}),
	)
	check(err)
	fmt.Println(bests[0])
}
