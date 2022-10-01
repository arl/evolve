package main

import (
	"image/color"
	"math/rand"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
)

func main() {
	myApp := app.New()
	w := myApp.NewWindow("TSP")

	tspw := newTSPWindow()
	tspw.buildUI(w)

	w.Resize(fyne.NewSize(800, 600))
	w.ShowAndRun()
}

func updatePath(w fyne.Window) {
	tick := time.NewTicker(500 * time.Millisecond)
	for {
		<-tick.C
		line := canvas.NewLine(color.White)
		width := 1 + rand.Intn(4)
		line.StrokeWidth = float32(width)
		w.SetContent(line)
	}
}
