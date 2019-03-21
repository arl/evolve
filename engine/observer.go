package engine

import (
	"github.com/arl/evolve"
)

// An Observer monitors the evolution of a population.
//
// Once registered within the evolution engine, an observer gets notified, after
// each completed epoch, with the population statistics (i.e once for every
// completed epoch).
type Observer interface {
	Observe(*evolve.PopulationStats)
}

type observerFunc struct{ f func(*evolve.PopulationStats) }

// The ObserverFunc type is an adapter to allow the use of
// ordinary functions as evolution observers. If f is a function
// with the appropriate signature, ObserverFunc(f) is an
// Observer that calls f.
func ObserverFunc(f func(*evolve.PopulationStats)) Observer { return &observerFunc{f: f} }

func (obs *observerFunc) Observe(stats *evolve.PopulationStats) { obs.f(stats) }
