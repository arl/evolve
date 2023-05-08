package main

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"time"

	"golang.org/x/exp/constraints"

	"github.com/arl/evolve"
	"github.com/arl/evolve/condition"
	"github.com/arl/evolve/crossover"
	"github.com/arl/evolve/engine"
	"github.com/arl/evolve/factory"
	"github.com/arl/evolve/generator"
	"github.com/arl/evolve/mutation"
	"github.com/arl/evolve/pkg/tsp"
	"github.com/arl/evolve/selection"
)

type algorithm struct {
	cfg  config
	eng  engine.Engine[[]byte]
	scsv *evolve.StatsToCSV[[]byte]
}

func (a *algorithm) setup(obs engine.Observer[[]byte]) error {
	eval := tsp.NewSymmetricEvaluator[byte](a.cfg.cities)

	var pipeline evolve.Pipeline[[]byte]

	// Define the crossover operator.
	xover := &evolve.Crossover[[]byte]{
		Mater:       crossover.PMX[byte]{},
		Probability: generator.Const(1.),
	}

	// Collision crossover
	// xover := &crossover.Collision2[byte]{
	// 	Probability: generator.Const(1.),
	// 	EdgeWeight: func(i, j int) float64 {
	// 		return eval.Distances[i][j]
	// 	},
	// }

	pipeline = append(pipeline, xover)

	const mutationRate = 0.05

	// Define the mutation operator.
	// rng := rand.New(mt19937.New())
	mut := evolve.NewSwitch[[]byte](
		// &mutation.Permutation[byte]{
		// 	Count:       generator.Const(1),
		// 	Amount:      generator.Uniform(1, 4 /*len(a.cfg.cities)*/, rng),
		// 	Probability: generator.Const(mutationRate),
		// },
		&mutation.SRS[byte]{
			Probability: generator.Const(mutationRate),
		},
		&mutation.CIM[byte]{
			Probability: generator.Const(mutationRate),
		},
	)
	pipeline = append(pipeline, mut)

	indices := make([]byte, len(a.cfg.cities))
	for i := 0; i < len(a.cfg.cities); i++ {
		indices[i] = byte(i)
	}

	generational := engine.Generational[[]byte]{
		Operator:  pipeline,
		Evaluator: eval,
		Selection: &selection.RouletteWheel[[]byte]{},
		// Selection: &selection.SUS[[]byte]{},
		NumElites: 4,
	}

	a.eng = engine.Engine[[]byte]{
		Factory:     factory.Permutation[byte](indices),
		Evaluator:   eval,
		Epocher:     &generational,
		Concurrency: runtime.NumCPU(),
	}
	var userAbort condition.UserAbort[[]byte]
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		<-c
		userAbort.Abort()
	}()

	a.eng.EndConditions = append(a.eng.EndConditions, &userAbort)
	if a.cfg.maxgen != 0 {
		a.eng.EndConditions = append(a.eng.EndConditions, condition.GenerationCount[[]byte](a.cfg.maxgen))
	}

	if a.cfg.csvpath != "" {
		csv, err := evolve.NewStatsToCSV[[]byte](a.cfg.csvpath)
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

func printStatsToCli[T constraints.Integer]() engine.Observer[[]T] {
	start := time.Now()
	last := start
	prevFitness := 0.0
	const refreshInterval = 1 * time.Second

	return engine.ObserverFunc(func(stats *evolve.PopulationStats[[]T]) {
		now := time.Now()
		if now.Sub(last) < refreshInterval && (stats.Generation%100 != 0 || prevFitness == stats.BestFitness) {
			return
		}
		last = now

		fmt.Printf("[%d]: distance: %.2f stddev: %.2f elapsed: %v\n", stats.Generation, stats.BestFitness, stats.StdDev, stats.Elapsed.Round(time.Millisecond))
		prevFitness = stats.BestFitness
	})
}
