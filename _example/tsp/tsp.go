package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"runtime"
	"time"

	"evolve/example/tsp/internal/tsp"

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

type algorithm struct {
	cfg  config
	eng  engine.Engine[[]int]
	scsv *evolve.StatsToCSV[[]int]
}

func (a *algorithm) setup(obs engine.Observer[[]int]) error {
	var pipeline operator.Pipeline[[]int]

	// Define the crossover operator.
	pmx := operator.NewCrossover[[]int](xover.PMX[int]{})
	pmx.Points = generator.Const(2) // unused for cycle crossover
	pmx.Probability = generator.Const(1.0)

	pipeline = append(pipeline, pmx)

	const mutationRate = 0.05

	// Define the mutation operator.
	rng := rand.New(mt19937.New(time.Now().UnixNano()))
	mut := operator.NewSwitch[[]int](
		&mutation.SliceOrder[int]{
			Count:       generator.Const(1),
			Amount:      generator.Uniform(1, len(a.cfg.cities), rng),
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

	indices := make([]int, len(a.cfg.cities))
	for i := 0; i < len(a.cfg.cities); i++ {
		indices[i] = i
	}

	eval := newRouteEvaluator(a.cfg.cities)

	generational := engine.Generational[[]int]{
		Operator:  pipeline,
		Evaluator: eval,
		Selection: &selection.RouletteWheel[[]int]{},
		Elites:    2,
	}

	a.eng = engine.Engine[[]int]{
		Factory:     factory.Permutation[int](indices),
		Evaluator:   eval,
		Epocher:     &generational,
		Concurrency: runtime.NumCPU(),
	}
	var userAbort condition.UserAbort[[]int]
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		<-c
		userAbort.Abort()
	}()

	a.eng.EndConditions = append(a.eng.EndConditions, &userAbort)
	if a.cfg.maxgen != 0 {
		a.eng.EndConditions = append(a.eng.EndConditions, condition.GenerationCount[[]int](a.cfg.maxgen))
	}

	if a.cfg.csvpath != "" {
		csv, err := evolve.NewStatsToCSV[[]int](a.cfg.csvpath)
		if err != nil {
			return err
		}
		a.scsv = csv
		a.eng.AddObserver(csv)
	}

	a.eng.AddObserver(obs)
	return nil
}

func (a *algorithm) run() error {
	_, cond, err := a.eng.Evolve(150)
	fmt.Printf("TSP ended, reason: %v\n", cond)
	return err
}

type config struct {
	cities  []tsp.Point2D
	maxgen  int
	csvpath string
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