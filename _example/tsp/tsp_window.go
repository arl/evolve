package main

/*
import (
	"fmt"
	"image"
	"image/color"
	"sync"
	"time"

	"evolve/example/tsp/internal/tsp"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/fogleman/gg"

	"github.com/arl/evolve"
	"github.com/arl/evolve/engine"
)

type tspWindow struct {
	running    bool
	maxw, maxh int // max cities coords
	cities     []tsp.Point2D

	done chan struct{}

	path       *canvas.Image
	img        *image.RGBA
	generation *widget.Label
	distance   *widget.Label
	stddev     *widget.Label
	elapsed    *widget.Label
}

func newTSPWindow(tspf *tsp.File) *tspWindow {
	cities := tspf.Nodes
	maxw, maxh := worldBounds(cities)
	fmt.Println("world bounds", maxw, "x", maxh)

	return &tspWindow{
		cities: cities,
		maxw:   int(maxw + 1),
		maxh:   int(maxh + 1),
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
			runTSP(w.cities, -1, w.updatePathAndStats())
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
	prev := time.Duration(0)
	prevFitness := 0.0
	const refreshInterval = 250 * time.Millisecond
	red := color.RGBA{255, 0, 0, 255}
	cityDiameter := 4.0

	return engine.ObserverFunc[[]int](func(stats *evolve.PopulationStats[[]int]) {
		if stats.Elapsed-prev < refreshInterval && (stats.Generation%100 != 0 || prevFitness == stats.BestFitness) {
			return
		}
		prev = stats.Elapsed

		fmt.Printf("[%d]: distance: %.2f\n", stats.Generation, stats.BestFitness)
		prevFitness = stats.BestFitness

		w.generation.SetText(fmt.Sprintf("generation: %d", stats.Generation))
		w.distance.SetText(fmt.Sprintf("distance: %.2f", stats.BestFitness))
		w.stddev.SetText(fmt.Sprintf("std dev: %.2f", stats.StdDev))
		w.elapsed.SetText(fmt.Sprintf("elapsed: %s", stats.Elapsed.Round(time.Millisecond)))

		dc := gg.NewContextForImage(w.img)
		dc.SetColor(color.White)
		dc.Clear()

		dc.SetColor(red)
		for i := 1; i < len(stats.Best); i++ {
			x := float64(w.cities[stats.Best[i]].X)
			y := float64(w.cities[stats.Best[i]].Y)
			dc.DrawPoint(x, y, cityDiameter)
		}
		dc.Stroke()

		dc.SetColor(color.Black)
		dc.MoveTo(float64(w.cities[stats.Best[0]].X), float64(w.cities[stats.Best[0]].Y))
		for i := 1; i < len(stats.Best); i++ {
			x := float64(w.cities[stats.Best[i]].X)
			y := float64(w.cities[stats.Best[i]].Y)
			dc.LineTo(x, y)
		}
		dc.SetLineWidth(1)
		dc.ClosePath()
		dc.Stroke()

		w.path.Image = dc.Image()
		canvas.Refresh(w.path)
	})
}
*/
