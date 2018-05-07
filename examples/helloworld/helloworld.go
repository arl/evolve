package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/aurelien-rainone/evolve"

	"github.com/aurelien-rainone/evolve/pkg/api"
	"github.com/aurelien-rainone/evolve/pkg/factory"
	"github.com/aurelien-rainone/evolve/pkg/operator"
	"github.com/aurelien-rainone/evolve/pkg/operator/mutation"
	"github.com/aurelien-rainone/evolve/pkg/operator/xover"
	"github.com/aurelien-rainone/evolve/pkg/selection"
	"github.com/aurelien-rainone/evolve/pkg/termination"
)

// This 'evaluator' assigns one "fitness point" for every character in the
// candidate string that doesn't match the corresponding position in the target
// string.
type evaluator string

func (s evaluator) Fitness(
	cand api.Candidate,
	pop []api.Candidate) float64 {

	var errors float64
	sc := cand.(string)
	for i := range sc {
		if sc[i] != string(s)[i] {
			errors++
		}
	}
	return errors
}

// Fitness is not natural, one fitness point represents an error, so the lower
// is better
func (evaluator) IsNatural() bool { return false }

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
		stringFactory *factory.String
		err           error
	)
	stringFactory, err = factory.NewString(string(alphabet), len(targetString))
	check(err)

	// 1st operator: string mutation
	mutation := mutation.NewString(string(alphabet))
	check(mutation.SetProb(0.02))

	// 2nd operator: string crossover
	xover := xover.New(xover.StringMater{})

	// Create a en operator pipeline applying first mutation, then crossover
	pipeline := operator.Pipeline{mutation, xover}

	eval := evaluator(targetString)

	var selector = selection.RouletteWheelSelection
	rng := rand.New(rand.NewSource(randomSeed()))

	engine := evolve.NewGenerationalEvolutionEngine(stringFactory,
		pipeline,
		eval,
		selector,
		rng)

	//engine.SetSingleThreaded(true)
	engine.AddEvolutionObserver(observer{})

	condition := termination.TargetFitness{Fitness: 0, Natural: false}

	result := engine.Evolve(100, 5, condition)
	fmt.Println(result)

	var conditions []api.TerminationCondition
	conditions, err = engine.SatisfiedTerminationConditions()
	check(err)
	for i, condition := range conditions {
		fmt.Printf("satified termination condition %v: %v\n",
			i, condition)
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

func (o observer) PopulationUpdate(data *api.PopulationData) {
	fmt.Printf("Generation %d: %s (%v)\n",
		data.GenNumber,
		data.BestCand,
		data.BestFitness)
}
