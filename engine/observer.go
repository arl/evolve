package engine

import (
	"github.com/arl/evolve"
)

// An Observer monitors the evolution of a population.
//
// Once registered within the evolution engine, an observer gets notified, after
// each completed epoch, with the population statistics (i.e once for every
// completed epoch).
type Observer[T any] interface {
	Observe(*evolve.PopulationStats[T])
}

// The ObserverFunc type is an adapter to allow the use of
// ordinary functions as evolution observers. If f is a function
// with the appropriate signature, ObserverFunc(f) is an
// Observer that calls f.
func ObserverFunc[T any](f func(*evolve.PopulationStats[T])) Observer[T] {
	return &observerFunc[T]{f: f}
}

type observerFunc[T any] struct {
	f func(*evolve.PopulationStats[T])
}

func (obs *observerFunc[T]) Observe(stats *evolve.PopulationStats[T]) {
	obs.f(stats)
}
