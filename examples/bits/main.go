package main

import (
	"fmt"
	"math/rand"

	"github.com/aurelien-rainone/evolve"
	"github.com/aurelien-rainone/evolve/pkg/api"
	"github.com/aurelien-rainone/evolve/pkg/bitstring"
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

func (evaluator) Fitness(cand api.Candidate, pop []api.Candidate) float64 {
	return float64(cand.(*bitstring.Bitstring).OnesCount())
}

func (evaluator) IsNatural() bool { return true }

// An implementation of the first exercise (page 31) from the book An
// Introduction to Genetic Algorithms, by Melanie Mitchell.  The algorithm
// evolves bit strings and the fitness function simply counts the number of ones
// in the bit string.  The evolution should therefore converge on strings that
// consist only of ones.
func main() {
	xover := xover.New(xover.BitstringMater{})
	xover.SetPoints(1)
	xover.SetProb(0.7)
	mut := mutation.NewBitstring()
	mut.SetProb(0.01)
	pipeline := operator.Pipeline{xover, mut}

	mt19937 := rand.New(random.NewMT19937(0))

	engine := evolve.NewGenerationalEngine(factory.NewBitstring(nbits),
		pipeline,
		evaluator{},
		selection.RouletteWheelSelection,
		mt19937)

	engine.AddObserver(observer{})

	best := engine.Evolve(
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
