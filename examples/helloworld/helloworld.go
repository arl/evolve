package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/aurelien-rainone/evolve"
	"github.com/aurelien-rainone/evolve/factory"
	"github.com/aurelien-rainone/evolve/framework"

	"github.com/aurelien-rainone/evolve/pkg/operator"
	"github.com/aurelien-rainone/evolve/pkg/operator/mutation"
	"github.com/aurelien-rainone/evolve/pkg/operator/xover"
	"github.com/aurelien-rainone/evolve/pkg/selection"
	"github.com/aurelien-rainone/evolve/termination"
)

func main() {
	var targetString = "HELLO WORLD"
	if len(os.Args) == 2 {
		targetString = strings.ToUpper(os.Args[1])
	}

	// Create a factory to generate random 11-character Strings.
	alphabet := make([]byte, 27)
	for c := byte('A'); c <= 'Z'; c++ {
		alphabet[c-'A'] = c
	}
	alphabet[26] = ' '

	for _, c := range targetString {
		if !strings.ContainsRune(string(alphabet), c) {
			fmt.Printf("Rune %c is not contained in the alphabet\n", c)
			os.Exit(1)
		}
	}

	var (
		stringFactory *factory.StringFactory
		err           error
	)
	stringFactory, err = factory.NewStringFactory(string(alphabet), len(targetString))
	check(err)

	// 1st operator: string mutation
	mutation := mutation.NewStringMutation(string(alphabet))
	check(mutation.SetProb(0.02))

	// 2nd operator: string crossover
	xover := xover.NewCrossover(xover.StringMater{})

	// Create a en operator pipeline applying first mutation, then crossover
	pipeline := operator.Pipeline{mutation, xover}

	fitnessEvaluator := newStringEvaluator(targetString)

	var selectionStrategy = selection.RouletteWheelSelection
	rng := rand.New(rand.NewSource(randomSeed()))

	engine := evolve.NewGenerationalEvolutionEngine(stringFactory,
		pipeline,
		fitnessEvaluator,
		selectionStrategy,
		rng)

	//engine.SetSingleThreaded(true)
	engine.AddEvolutionObserver(observer{})
	result := engine.Evolve(100, 5, termination.NewTargetFitness(0, false))
	fmt.Println(result)

	var conditions []framework.TerminationCondition
	conditions, err = engine.SatisfiedTerminationConditions()
	check(err)
	for i, condition := range conditions {
		fmt.Printf("satified termination condition %v %T: %v\n",
			i, condition, condition)
	}
}

func randomSeed() int64 {
	return time.Now().UnixNano()
}

func check(err error) {
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

type observer struct{}

func (o observer) PopulationUpdate(data *framework.PopulationData) {
	fmt.Printf("Generation %d: %s (%v)\n", data.GenerationNumber(), data.BestCandidate(),
		data.BestCandidateFitness())
}
