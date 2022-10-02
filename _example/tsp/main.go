package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

func main() {
	myApp := app.New()
	w := myApp.NewWindow("TSP")

	tspw := newTSPWindow()
	tspw.buildUI(w)

	go func() {
		<-tspw.done
		w.Close()
	}()

	w.Resize(fyne.NewSize(800, 600))
	w.ShowAndRun()
}
