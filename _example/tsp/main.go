package main

import (
	"flag"
	"log"
	"os"

	"evolve/example/tsp/internal/tsp"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

func main() {
	cpuprofile := flag.String("cpuprofile", "", "write cpu profile to `file`")
	memprofile := flag.String("memprofile", "", "write memory profile to `file`")
	memprofilerate := flag.Int("memprofilerate", 0, "set runtime.MemProfileRate to `rate`")
	statsvizAddr := flag.String("statsviz", "", "enable statsviz endpoint at `host:port`")
	nogui := flag.Bool("nogui", false, "disable gui, just starts the algorithm")
	maxgen := flag.Int("maxgen", -1, "maximum generation, -1, run forever")
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

	tspf, err := tsp.Load(f)
	if err != nil {
		log.Fatal(err)
	}

	if *nogui {
		runTSP(tspf.Nodes, *maxgen, printStatsToCli())
		return
	}

	app := app.New()
	w := app.NewWindow("TSP")

	tspw := newTSPWindow(tspf)
	tspw.buildUI(w)

	go func() {
		<-tspw.done
		w.Close()
	}()

	w.Resize(fyne.NewSize(800, 600))
	w.ShowAndRun()
}
