package main

import (
	"flag"
	"log"
	"os"

	"github.com/arl/evolve/pkg/tsp"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/widget/material"
)

var alg algorithm

func main() {
	cpuprofile := flag.String("cpuprofile", "", "write cpu profile to `file`")
	memprofile := flag.String("memprofile", "", "write memory profile to `file`")
	memprofilerate := flag.Int("memprofilerate", 0, "set runtime.MemProfileRate to `rate`")
	statsvizAddr := flag.String("statsviz", "", "enable statsviz endpoint at `host:port`")
	nogui := flag.Bool("nogui", false, "disable gui, just starts the algorithm")
	maxgen := flag.Int("maxgen", 0, "max number of generations to evolve. 0:forever")
	tspfname := flag.String("tspfile", "berlin52.tsp", "tspfile to load, by default, pre-load berlin52")
	flag.Parse()

	stopProf := prof(*statsvizAddr, *cpuprofile, *memprofile, *memprofilerate)

	// Fill config
	tspf, err := tsp.LoadFromFile(*tspfname)
	if err != nil {
		log.Fatal(err)
	}

	if *nogui {
		cliRun(tspf.Nodes, *maxgen)
		return
	}

	guiRun(tspf, *maxgen, stopProf)
}

func cliRun(cities []tsp.Point2D, maxgen int) {
	alg.cfg = config{cities: cities, maxgen: maxgen}
	if err := alg.setup(printStatsToCli[byte]()); err != nil {
		log.Fatalf("setup failed: %v", err)
	}
	alg.run()
}

// guiRun runs the algorithm in a gio user interface.
// beforeExit is called before os.Exit.
func guiRun(tspf *tsp.File, maxgen int, beforeExit func()) {
	go func() {
		theme := material.NewTheme(gofont.Collection())

		ui := newUI(theme, tspf)
		w := app.NewWindow(
			app.Title("evolve/TSP"),
			app.Size(1400, 900),
		)
		if err := ui.run(w); err != nil {
			log.Fatal(err)
		}
		if beforeExit != nil {
			beforeExit()
		}
		os.Exit(0)
	}()

	app.Main()
}
