package main

import (
	"image/color"
	"sync/atomic"
	"time"

	"github.com/arl/evolve"
	"github.com/arl/evolve/engine"
	"github.com/arl/evolve/pkg/tsp"

	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/arl/gioexp/component/property"
)

var (
	bgColor    = color.NRGBA{R: 247, G: 231, B: 190, A: 255}
	panelColor = color.NRGBA{R: 218, G: 234, B: 240, A: 255}

	borderColor = color.NRGBA{R: 131, G: 140, B: 143, A: 255}
	borderWidth = unit.Dp(1)

	propertyColor = color.NRGBA{R: 216, G: 202, B: 227, A: 255}
)

var (
	entriesPanel = Panel{
		Axis: layout.Vertical,
		Size: unit.Dp(270),

		Background:  panelColor,
		Border:      borderColor,
		BorderWidth: borderWidth,
	}

	entriesHeaderPanel = Panel{
		Axis: layout.Horizontal,
		Size: unit.Dp(80),

		Background:  panelColor,
		Border:      borderColor,
		BorderWidth: borderWidth,
	}
)

type (
	C = layout.Context
	D = layout.Dimensions
)

// state holds the application state
type state struct {
	stats *evolve.PopulationStats[[]byte]
	tspf  *tsp.File
}

type UI struct {
	state state
	theme *material.Theme

	list *property.List

	startButton *startButton
	zoomable    *Zoomable
	pathWidget  *pathWidget[byte]
}

func newUI(theme *material.Theme, tspf *tsp.File) *UI {
	return &UI{
		theme: theme,
		state: state{
			tspf:  tspf,
			stats: &evolve.PopulationStats[[]byte]{},
		},
		list: property.NewList(),
	}
}

func (ui *UI) run(w *app.Window) error {
	ui.startButton = &startButton{}
	ui.pathWidget = newPathWidget[byte](ui.state.tspf.Nodes)

	gen := property.NewInt(0)
	gen.Editable = false
	ui.list.Add("Generation", gen)

	dist := property.NewFloat64(0)
	dist.Editable = false
	dist.SetFormat('f', 1)
	ui.list.Add("Distance", dist)

	stddev := property.NewFloat64(0)
	stddev.Editable = false
	ui.list.Add("Std dev", stddev)

	elapsed := property.NewString("")
	elapsed.Editable = false
	ui.list.Add("Elapsed", elapsed)

	solutions := make(chan *evolve.PopulationStats[[]byte])
	var prev, paused time.Duration
	prevFitness := 0.0

	observer := engine.ObserverFunc(func(stats *evolve.PopulationStats[[]byte]) {
		// Handle paused UI. We can do this here since evolution observers are
		// all executed synchronously after each epoch, so blocking here means
		// blocking the whole evolution ^-^.
		before := time.Now()
		for !ui.startButton.isRunning() {
			// UI is paused
			time.Sleep(100 * time.Millisecond)
		}
		paused += time.Since(before)

		// In case of many consecutive improvements of the solution, we want
		// anyway to limit us to to drawing 30 fps.
		const refreshInterval = 1 * time.Second / 30
		if stats.Elapsed-prev < refreshInterval &&
			prevFitness == stats.BestFitness {
			return
		}

		prev = stats.Elapsed
		prevFitness = stats.BestFitness

		solutions <- stats
	})

	var ops op.Ops
	starting := false
	for {
		select {
		case stats := <-solutions:
			// Substract paused time
			stats.Elapsed -= paused
			ui.state.stats = stats

			gen.SetValue(stats.Generation)
			dist.SetValue(stats.BestFitness)
			stddev.SetValue(stats.StdDev)
			elapsed.SetValue(ui.state.stats.Elapsed.Truncate(time.Millisecond).String())

		case e := <-w.Events():
			switch e := e.(type) {
			case system.FrameEvent:
				gtx := layout.NewContext(&ops, e)
				starting = ui.startButton.handleClicked()
				ui.Layout(gtx)
				e.Frame(gtx.Ops)
			case system.DestroyEvent:
				return e.Err
			}
		}

		if starting {
			alg.cfg = config{cities: ui.state.tspf.Nodes, maxgen: 0}
			alg.cfg.csvpath = "stats.csv"
			if err := alg.setup(observer); err != nil {
				return err
			}
			if alg.scsv != nil {
				defer alg.scsv.Close()
			}

			go alg.run()
		}
	}
}

func (ui *UI) Layout(gtx C) D {
	// TODO(arl) we probably can move all startButton logic into a method of UI
	// instead of a specific struct.
	drawPath := !ui.startButton.isStarted()

	gtx.Constraints.Min = gtx.Constraints.Max
	return layout.Flex{
		Axis: layout.Horizontal,
	}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return layout.Flex{
				Axis:    layout.Vertical,
				Spacing: layout.SpaceEnd,
			}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					gtx.Constraints.Max.X = 400
					return ui.list.Layout(ui.theme, gtx)
				}),
				layout.Rigid(func(gtx C) D {
					return ui.startButton.Layout(ui.theme, gtx)
				}),
			)
		}),
		layout.Flexed(1, func(gtx C) D {
			return ui.pathWidget.Layout(drawPath, ui.state.stats.Best, gtx)
		}),
	)
}

// startButton is a single button used to start, pause and resume the
// evolutionnary algorithm.
type startButton struct {
	widget.Clickable
	started bool
	running atomic.Bool // running/paused
}

// isStarted returns whether the start button has been clicked at least once.
func (btn *startButton) isStarted() bool {
	return btn.started
}

// isRunning returns whether the button is currently running (as opposed to paused).
func (btn *startButton) isRunning() bool {
	return btn.running.Load()
}

func (btn *startButton) Layout(theme *material.Theme, gtx C) D {
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
	button := material.Button(theme, &btn.Clickable, txt)
	button.Inset = layout.UniformInset(12)
	return button.Layout(gtx)
}

// handleClicked switches the start button internal state, and return whether
// the button has been clicked for the first time.
func (btn *startButton) handleClicked() (firstClick bool) {
	if btn.Clicked() {
		if !btn.started {
			btn.started = true
			btn.running.Store(true)
			return true
		}
		btn.running.Store(!btn.running.Load())
	}
	return false
}
