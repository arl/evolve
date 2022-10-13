package main

import (
	"log"
	"net"
	"net/http"
	"os"
	"runtime"
	"runtime/pprof"

	"github.com/arl/statsviz"
)

// funcstack is a poor man's defer stack.
type funcstack struct {
	s []func()
}

// add and adde adds functions to the defer stack.
func (s *funcstack) add(f func())        { s.s = append(s.s, f) }
func (s *funcstack) adde(f func() error) { s.s = append(s.s, func() { _ = f() }) }

// run runs all functions in the stack, running the last inserted first.
func (s *funcstack) run() {
	for i := len(s.s) - 1; i >= 0; i-- {
		s.s[i]()
	}
}

func prof(stats, cpu, mem string, memrate int) (stop func()) {
	var _defer funcstack
	if cpu != "" {
		f, err := os.Create(cpu)
		if err != nil {
			log.Fatalf("creating cpu profile: %s", err)
			return
		}
		_defer.adde(f.Close)

		pprof.StartCPUProfile(f)
		_defer.add(pprof.StopCPUProfile)

		_defer.add(func() { log.Println("cpu profile created:", cpu) })
	}

	if mem != "" {
		if memrate != 0 {
			runtime.MemProfileRate = memrate
		}
		_defer.add(func() { writeHeapProfile(mem) })
	}

	if stats != "" {
		startStatsviz(stats)
	}

	return _defer.run
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
