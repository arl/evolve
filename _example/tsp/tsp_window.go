package main

import (
	"fmt"
	"image"
	"image/color"
	"math/rand"
	"os"
	"os/signal"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
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
	"github.com/fogleman/gg"
)

type tspWindow struct {
	running    bool
	maxw, maxh int // max cities coords
	cities     []point

	done chan struct{}

	path       *canvas.Image
	img        *image.RGBA
	generation *widget.Label
	distance   *widget.Label
	stddev     *widget.Label
	elapsed    *widget.Label
}

func newTSPWindow() *tspWindow {
	cities := berlin52
	maxw, maxh := worldBounds(cities)
	fmt.Println("world bounds", maxw, "x", maxh)

	return &tspWindow{
		cities: cities,
		maxw:   maxw,
		maxh:   maxh,
		done:   make(chan struct{}),
	}
}

func (w *tspWindow) buildUI(wnd fyne.Window) {
	// main vertical layout with:
	// - controls at the top
	// - path visualization and stats at the bottom

	var once sync.Once
	startButton := widget.NewButton("start", func() {
		once.Do(func() {
			runTSP(w.cities, w.updatePathAndStats())
			close(w.done)
		})
	})

	controls := container.New(layout.NewHBoxLayout(), startButton)

	w.img = image.NewRGBA(image.Rect(0, 0, w.maxw, w.maxh))
	w.path = canvas.NewImageFromImage(w.img)
	w.path.FillMode = canvas.ImageFillContain
	w.path.SetMinSize(fyne.NewSize(float32(800), float32(600)))

	w.generation = widget.NewLabel("generation: ")
	w.distance = widget.NewLabel("distance: ")
	w.stddev = widget.NewLabel("std dev: ")
	w.elapsed = widget.NewLabel("elapsed: ")

	stats := container.New(layout.NewVBoxLayout(), w.generation, w.distance, w.stddev, w.elapsed)
	pathAndStats := container.New(layout.NewHBoxLayout(), w.path, stats)

	content := container.New(layout.NewVBoxLayout(), controls, layout.NewSpacer(), pathAndStats)
	wnd.SetContent(content)
}

func (w *tspWindow) updatePathAndStats() engine.Observer[[]int] {
	start := time.Now()
	prevFitness := 0.0

	return engine.ObserverFunc[[]int](func(stats *evolve.PopulationStats[[]int]) {
		if stats.Generation%100 != 0 && prevFitness == stats.BestFitness {
			return
		}

		fmt.Printf("[%d]: distance: %v\n", stats.Generation, stats.BestFitness)
		prevFitness = stats.BestFitness

		w.generation.SetText(fmt.Sprintf("generation: %d", stats.Generation))
		w.distance.SetText(fmt.Sprintf("distance: %f", stats.BestFitness))
		w.stddev.SetText(fmt.Sprintf("std dev: %.3f", stats.StdDev))
		w.elapsed.SetText(fmt.Sprintf("elapsed: %s", time.Since(start).Round(time.Millisecond)))

		dc := gg.NewContextForImage(w.img)
		dc.SetColor(color.White)
		dc.Clear()
		dc.SetColor(color.Black)
		dc.MoveTo(float64(w.cities[stats.Best[0]].X), float64(w.cities[stats.Best[0]].Y))
		for i := 1; i < len(stats.Best); i++ {
			dc.LineTo(float64(w.cities[stats.Best[i]].X), float64(w.cities[stats.Best[i]].Y))
		}
		dc.SetLineWidth(2)
		dc.ClosePath()
		dc.Stroke()

		w.path.Image = dc.Image()
		canvas.Refresh(w.path)
	})
}

type point struct{ X, Y int }

func runTSP(cities []point, obs engine.Observer[[]int]) (*evolve.Population[[]int], error) {
	var pipeline operator.Pipeline[[]int]

	// Define the crossover operator.
	xover := xover.New[[]int](xover.PMX[int]{})
	xover.Points = generator.Const(2) // unused for cycle crossover
	xover.Probability = generator.Const(1.0)
	pipeline = append(pipeline, xover)

	// Define the mutation operator.
	rng := rand.New(mt19937.New(time.Now().UnixNano()))
	mut := &mutation.SliceOrder[int]{
		Count:       generator.NewPoisson[int](generator.Const(2.0), rng),
		Amount:      generator.NewPoisson[int](generator.Const(4.0), rng),
		Probability: generator.Const(2.0),
	}
	pipeline = append(pipeline, mut)

	indices := make([]int, len(cities))
	for i := 0; i < len(cities); i++ {
		indices[i] = i
	}

	eval := newRouteEvaluator(cities)

	generational := engine.Generational[[]int]{
		Operator:  pipeline,
		Evaluator: eval,
		// Selection: &selection.Tournament[[]int]{
		// 	Probability: generator.Const(0.7),
		// },
		Selection: &selection.RouletteWheel[[]int]{},
		// Selection: &selection.SigmaScaling[[]int]{
		// 	&selection.RouletteWheel[[]int]{},
		// },
		Elites: 2,
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

	pop, cond, err := eng.Evolve(150)
	fmt.Printf("TSP ended, reason: %v\n", cond)

	return pop, err
}
