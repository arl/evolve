package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/arl/bitstring"
	"github.com/arl/evolve"
	"github.com/arl/evolve/condition"
	"github.com/arl/evolve/crossover"
	"github.com/arl/evolve/engine"
	"github.com/arl/evolve/generator"
	"github.com/arl/evolve/mutation"
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
	xover := evolve.NewCrossover[*bitstring.Bitstring](crossover.BitstringMater{})
	xover.Probability = generator.Const(0.7)
	xover.Points = generator.Const(1)

	// Define the mutation
	mut := evolve.NewMutation[*bitstring.Bitstring](&mutation.Bitstring{
		Probability: generator.Const(0.01),
		FlipCount:   generator.Const(1),
	})

	eval := evolve.EvaluatorFunc(
		true, // natural fitness (higher is better)
		func(cand *bitstring.Bitstring) float64 {
			// our evaluator counts the ones in the bitstring
			return float64(cand.OnesCount())
		})

	epocher := engine.Generational[*bitstring.Bitstring]{
		Operator:  evolve.Pipeline[*bitstring.Bitstring]{xover, mut},
		Evaluator: eval,
		Selection: selection.RouletteWheel[*bitstring.Bitstring]{},
		NumElites: 2, // best 2 candidates gets copied to the next generation, no matter what.
	}

	// bitstring.Random(length uint, rng *rand.Rand)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	factory := evolve.FactoryFunc[*bitstring.Bitstring](func(*rand.Rand) *bitstring.Bitstring {
		return bitstring.Random(nbits, rng)
	})

	eng := &engine.Engine[*bitstring.Bitstring]{
		Factory:   factory,
		Evaluator: eval,
		Epocher:   &epocher,
		EndConditions: []evolve.Condition[*bitstring.Bitstring]{condition.TargetFitness[*bitstring.Bitstring]{
			Fitness: float64(nbits),
			Natural: true,
		}},
	}

	eng.AddObserver(
		engine.ObserverFunc(func(stats *evolve.PopulationStats[*bitstring.Bitstring]) {
			log.Printf("Generation %d: %s (%v)\n", stats.Generation, stats.Best, stats.BestFitness)
		}))

	bests, _, err := eng.Evolve(100)
	check(err)
	fmt.Println(bests.Candidates[0], "=>", bests.Fitness[0])
}
