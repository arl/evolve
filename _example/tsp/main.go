package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"evolve/example/tsp/internal/tsp"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/widget/material"
)

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
	defer stopProf()

	// Fill config
	f, err := os.Open(*tspfname)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	fmt.Println("loaded", *tspfname, "successfully")

	tspf, err := tsp.Load(f)
	if err != nil {
		log.Fatal(err)
	}

	if *nogui {
		runTSP(config{cities: tspf.Nodes, maxgen: *maxgen}, printStatsToCli())
		return
	}

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
		os.Exit(0)
	}()

	app.Main()
}
