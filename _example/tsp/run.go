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

	"evolve/example/tsp/internal/tsp"
)

type config struct {
	cities []tsp.Point2D
	maxgen int
}

func runTSP(cfg config, obs engine.Observer[[]int]) (*evolve.Population[[]int], *evolve.PopulationStats[[]int], error) {
	var pipeline operator.Pipeline[[]int]

	// Define the crossover operator.
	pmx := xover.New[[]int](xover.PMX[int]{})
	pmx.Points = generator.Const(2) // unused for cycle crossover
	pmx.Probability = generator.Const(1.0)

	pipeline = append(pipeline, pmx)

	const mutationRate = 0.05

	// Define the mutation operator.
	rng := rand.New(mt19937.New(time.Now().UnixNano()))
	mut := operator.NewSwitch[[]int](
		&mutation.SliceOrder[int]{
			Count:       generator.Const(1),
			Amount:      generator.Uniform(1, len(cfg.cities), rng),
			Probability: generator.Const(mutationRate),
		},
		&mutation.SRS[int]{
			Probability: generator.Const(mutationRate),
		},
		&mutation.CIM[int]{
			Probability: generator.Const(mutationRate),
		},
	)
	pipeline = append(pipeline, mut)

	indices := make([]int, len(cfg.cities))
	for i := 0; i < len(cfg.cities); i++ {
		indices[i] = i
	}

	eval := newRouteEvaluator(cfg.cities)

	generational := engine.Generational[[]int]{
		Operator:  pipeline,
		Evaluator: eval,
		Selection: &selection.RouletteWheel[[]int]{},
		Elites:    2,
	}

	eng := engine.Engine[[]int]{
		Factory:   factory.Permutation[int](indices),
		Evaluator: eval,
		Epocher:   &generational,
		// Concurrency: runtime.NumCPU() * 2,
		Concurrency: 1,
	}
	var userAbort condition.UserAbort[[]int]
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		<-c
		userAbort.Abort()
	}()

	eng.EndConditions = append(eng.EndConditions, &userAbort)
	if cfg.maxgen != 0 {
		eng.EndConditions = append(eng.EndConditions, condition.GenerationCount[[]int](cfg.maxgen))
	}

	var latestStats *evolve.PopulationStats[[]int]
	eng.AddObserver(engine.ObserverFunc(func(stats *evolve.PopulationStats[[]int]) {
		obs.Observe(stats)
		latestStats = stats
	}))

	pop, cond, err := eng.Evolve(150)
	fmt.Printf("TSP ended, reason: %v\n", cond)

	// eng.
	return pop, latestStats, err
}

func printStatsToCli() engine.Observer[[]int] {
	start := time.Now()
	last := start
	prevFitness := 0.0
	const refreshInterval = 1 * time.Second

	return engine.ObserverFunc(func(stats *evolve.PopulationStats[[]int]) {
		now := time.Now()
		if now.Sub(last) < refreshInterval && (stats.Generation%100 != 0 || prevFitness == stats.BestFitness) {
			return
		}
		last = now

		fmt.Printf("[%d]: distance: %.2f stddev: %.2f elapsed: %v\n", stats.Generation, stats.BestFitness, stats.StdDev, stats.Elapsed.Round(time.Millisecond))
		prevFitness = stats.BestFitness
	})
}
