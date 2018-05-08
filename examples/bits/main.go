package main

import (
	"fmt"
	"log"
	"math/rand"

	"github.com/aurelien-rainone/evolve/pkg/api"
	"github.com/aurelien-rainone/evolve/pkg/bitstring"
	"github.com/aurelien-rainone/evolve/pkg/engine"
	"github.com/aurelien-rainone/evolve/pkg/factory"
	"github.com/aurelien-rainone/evolve/pkg/operator"
	"github.com/aurelien-rainone/evolve/pkg/operator/mutation"
	"github.com/aurelien-rainone/evolve/pkg/operator/xover"
	"github.com/aurelien-rainone/evolve/pkg/selection"
	"github.com/aurelien-rainone/evolve/pkg/termination"
	"github.com/aurelien-rainone/evolve/random"
)

const nbits = 20

// fitness evaluator that simply counts the number of ones in a string
type evaluator struct{}

func (evaluator) Fitness(cand interface{}, pop []interface{}) float64 {
	return float64(cand.(*bitstring.Bitstring).OnesCount())
}

func (evaluator) IsNatural() bool { return true }

func check(err error) {
	if err != nil {
		log.Fatalln("quitting with error:", err)
	}
}

// An implementation of the first exercise (page 31) from the book An
// Introduction to Genetic Algorithms, by Melanie Mitchell.  The algorithm
// evolves bit strings and the fitness function simply counts the number of ones
// in the bit string.  The evolution should therefore converge on strings that
// consist only of ones.
func main() {
	// define the crossover
	xover := xover.New(xover.BitstringMater{})
	check(xover.SetPoints(1))
	check(xover.SetProb(0.7))

	// define the mutation
	mut := mutation.NewBitstring()
	check(mut.SetProb(0.01))
	pipeline := operator.Pipeline{xover, mut}

	mt19937 := rand.New(random.NewMT19937(0))

	eng := engine.NewGenerational(factory.NewBitstring(nbits),
		pipeline,
		evaluator{},
		selection.RouletteWheel,
		mt19937)

	eng.AddObserver(observer{})

	best := eng.Evolve(
		100, // 100 candidates in the population
		0,   // no elitism
		termination.TargetFitness{Fitness: float64(nbits), Natural: true})

	fmt.Println(best)
}

type observer struct{}

func (o observer) PopulationUpdate(data *api.PopulationData) {
	fmt.Printf("Generation %d: %s (%v)\n",
		data.GenNumber,
		data.BestCand,
		data.BestFitness)
}
