package main

import (
	"evolve/example/tsp/internal/tsp"
	"fmt"
	"image"
	"image/color"
	"sync/atomic"
	"time"

	"github.com/arl/evolve"
	"github.com/arl/evolve/engine"

	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type (
	C = layout.Context
	D = layout.Dimensions
)

type gui struct {
	theme *material.Theme
	tspf  *tsp.File
}

func (g *gui) run(w *app.Window) error {
	var ops op.Ops

	zoomed := zoomable{}

	btn := startButton{theme: g.theme}

	pw := newPathWidget(g.tspf.Nodes)

	solutions := make(chan []int)
	var prev, paused time.Duration
	prevFitness := 0.0

	var observer = engine.ObserverFunc(func(stats *evolve.PopulationStats[[]int]) {
		// Handle paused UI. We can do this here since evolution observers are
		// all executed synchronously after each epoch, so blocking here means
		// blocking the whole evolution ^-^.
		bef := time.Now()
		for !btn.isRunning() {
			// UI is paused
			time.Sleep(100 * time.Millisecond)
		}
		paused += time.Since(bef)

		// In case of many consecutive improvements of the solution, we want
		// anyway to limit us to to drawning 30 fps.
		const refreshInterval = 1 * time.Second / 30
		if stats.Elapsed-prev < refreshInterval &&
			prevFitness == stats.BestFitness {
			return
		}
		prev = stats.Elapsed
		prevFitness = stats.BestFitness

		fmt.Printf("[%d]: distance: %.2f\n", stats.Generation, stats.BestFitness)
		solutions <- stats.Best
	})

	var bestPath []int
	var firstStart = true
	for {
		select {
		case path := <-solutions:
			bestPath = path
		case e := <-w.Events():
			switch e := e.(type) {
			case system.FrameEvent:
				gtx := layout.NewContext(&ops, e)

				if started := btn.handleClicked(); started && firstStart {
					// Start the TSP generic algorithm.
					cfg := config{cities: g.tspf.Nodes, maxgen: 0}
					go runTSP(cfg, observer)
					firstStart = false
				}

				layout.Flex{
					Axis:    layout.Vertical,
					Spacing: layout.SpaceBetween,
				}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						return layout.Flex{
							Axis:    layout.Horizontal,
							Spacing: layout.SpaceBetween,
						}.Layout(gtx,
							layout.Flexed(1, btn.Layout))
					}),
					layout.Flexed(1, func(gtx C) D {
						return zoomed.Layout(gtx, func(gtx C) D {
							return pw.layout(gtx, firstStart, bestPath)
						})
					}),
				)
				e.Frame(gtx.Ops)
			case system.DestroyEvent:
				return e.Err
			}
		}
	}
}

type pathWidget struct {
	citymax f32.Point
	cities  []tsp.Point2D
}

func newPathWidget(cities []tsp.Point2D) *pathWidget {
	max := func(a, b float32) float32 {
		if a > b {
			return a
		}
		return b
	}

	// Compute world bounds
	var citymax f32.Point
	for _, c := range cities {
		citymax.X = max(citymax.X, float32(c.X))
		citymax.Y = max(citymax.Y, float32(c.Y))
	}
	fmt.Println("world bounds", citymax)

	return &pathWidget{cities: cities, citymax: citymax}
}

var (
	pathColor = color.NRGBA{A: 255}
	dotColor  = color.NRGBA{R: 200, A: 255}
)

func (pw *pathWidget) layout(gtx C, onlyCities bool, sol []int) D {
	// Draw cities as red dots
	const cityRadius = 5
	for i := range pw.cities {
		city := pw.cities[i]
		circle := clip.Ellipse{
			Min: image.Pt(int(city.X-cityRadius), int(city.Y-cityRadius)),
			Max: image.Pt(int(city.X+cityRadius), int(city.Y+cityRadius)),
		}.Op(gtx.Ops)
		paint.FillShape(gtx.Ops, dotColor, circle)
	}

	if !onlyCities && len(sol) != 0 {
		// At start we may not have received the first solution yet.
		p := clip.Path{}
		p.Begin(gtx.Ops)
		pt := f32.Pt(float32(pw.cities[sol[0]].X), float32(pw.cities[sol[0]].Y))
		p.MoveTo(pt)
		for i := 1; i < len(sol); i++ {
			pt := f32.Pt(float32(pw.cities[sol[i]].X), float32(pw.cities[sol[i]].Y))
			p.LineTo(pt)
		}
		p.LineTo(pt)
		paint.FillShape(gtx.Ops, pathColor, clip.Stroke{Path: p.End(), Width: 1}.Op())
	}

	op.InvalidateOp{}.Add(gtx.Ops)
	return layout.Dimensions{Size: gtx.Constraints.Max}
}

// startButton is a single button used to start, pause and resume the
// evolutionnary algorithm.
type startButton struct {
	theme *material.Theme

	widget.Clickable
	started bool
	running atomic.Bool // running/paused
}

func (btn *startButton) isRunning() bool {
	return btn.running.Load()
}

func (btn *startButton) Layout(gtx C) D {
	txt := ""
	if !btn.started {
		txt = "Start"
	} else {
		if btn.running.Load() {
			txt = "Pause"
		} else {
			txt = "Resume"
		}
	}
	button := material.Button(btn.theme, &btn.Clickable, txt)
	return button.Layout(gtx)
}

// handleClicked switches the start button internal state, and return whether
// the button has been clicked for the first time.
func (btn *startButton) handleClicked() (firstClick bool) {
	if btn.Clicked() {
		if !btn.started {
			btn.started = true
			btn.running.Store(true)
			firstClick = true
		} else {
			btn.running.Store(!btn.running.Load())
		}
	}
	return
}
