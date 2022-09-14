package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"time"

	"github.com/arl/evolve"
	"github.com/arl/evolve/condition"
	"github.com/arl/evolve/engine"
	"github.com/arl/evolve/factory"
	"github.com/arl/evolve/generator"
	"github.com/arl/evolve/operator"
	"github.com/arl/evolve/operator/mutation"
	"github.com/arl/evolve/operator/xover"
	"github.com/arl/evolve/pkg/mt19937"
	"github.com/arl/evolve/selection"
)

const (
	numCities  = 26
	plotEach   = 200
	xmax, ymax = 200, 200
)

func runTSP(cities []point, obs engine.Observer[[]int]) (*evolve.Population[[]int], error) {
	rng := rand.New(mt19937.New(time.Now().UnixNano()))

	// Define the crossover operator.
	xover := xover.New[[]int](xover.PMX[int]{})
	xover.Points = generator.Const(2)
	xover.Probability = generator.Const(1.0)

	// Define the mutation operator.

	mut := &mutation.SliceOrder[int]{
		Count:       generator.NewPoisson[int](generator.Const(2.0), rng),
		Amount:      generator.NewPoisson[int](generator.Const(4.0), rng),
		Probability: generator.Const(0.1),
	}

	indices := make([]int, len(cities))
	for i := 0; i < len(cities); i++ {
		indices[i] = i
	}

	eval := newRouteEvaluator(cities)

	generational := engine.Generational[[]int]{
		Operator:  operator.Pipeline[[]int]{xover, mut},
		Evaluator: eval,
		Selection: &selection.RouletteWheel[[]int]{},
		Elites:    4,
	}

	eng := engine.Engine[[]int]{
		Factory:   factory.Permutation[int](indices),
		Evaluator: eval,
		Epocher:   &generational,
	}
	var userAbort condition.UserAbort[[]int]
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		<-c
		userAbort.Abort()
	}()

	eng.EndConditions = append(eng.EndConditions, &userAbort)

	eng.AddObserver(obs)

	pop, cond, err := eng.Evolve(100)
	fmt.Printf("TSP ended, reason: %v\n", cond)

	return pop, err
}
