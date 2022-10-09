package main

import (
	"flag"
	"log"
	"net"
	"net/http"
	"os"
	"runtime"
	"runtime/pprof"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/arl/statsviz"
)

func main() {
	cpuprofile := flag.String("cpuprofile", "", "write cpu profile to `file`")
	memprofile := flag.String("memprofile", "", "write memory profile to `file`")
	memprofilerate := flag.Int64("memprofilerate", 0, "set runtime.MemProfileRate to `rate`")
	statsvizAddr := flag.String("statsviz", "", "enable statsviz endpoint at `host:port`")
	flag.Parse()

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatalf("creating cpu profile: %s", err)
			return
		}
		defer f.Close()
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
		defer log.Println("cpu profile created:", cpuprofile)
	}

	if *memprofile != "" {
		if *memprofilerate != 0 {
			runtime.MemProfileRate = int(*memprofilerate)
		}
		defer writeHeapProfile(*memprofile)
	}

	if *statsvizAddr != "" {
		startStatsviz(*statsvizAddr)
	}

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

func startStatsviz(addr string) {
	_, _, err := net.SplitHostPort(addr)
	if err != nil {
		log.Fatalf("statsviz: invalid address: %v", err)
	}
	statsviz.RegisterDefault()
	log.Println("starting statsviz endpoint at address", addr)
	go func() {
		log.Println(http.ListenAndServe(addr, nil))
	}()
}

func writeHeapProfile(fname string) {
	f, err := os.Create(fname)
	if err != nil {
		log.Fatal("could not create memory profile: ", err)
	}
	defer f.Close()
	runtime.GC() // get up-to-date statistics
	if err := pprof.WriteHeapProfile(f); err != nil {
		log.Fatal("could not write memory profile: ", err)
	}
	log.Println("memory profile created:", fname)
}
