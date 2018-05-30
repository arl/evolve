package api

// Observer is the interface that wraps the PopulationUpdate method. Observers
// give the opportunity to monitor the state of evolutionary algorithms.
//
// Depending on the parameters of the evolutionary program, an observer may
// be invoked dozens or hundreds of times a second, especially when the population
// size is small as this leads to shorter generations. The processing performed by an
// evolution observer should be reasonably short-lived so as to avoid slowing down
// the evolution.
//
// XXX:  just thinking out loud, but why is this in api package, observer are
// only used in the context of the engine, called engine.Observer would also
// make for a clearer name
type Observer interface {

	// PopulationUpdate is called when the state of the population has changed
	// (typically at the end of a generation) with data about the current
	// generation.
	PopulationUpdate(data *PopulationData)
}

// TODO: try to come up with a better and short name for PopulationUpdate
// and PopulationData maybe
type observerFunc struct{ f func(*PopulationData) }

// ObserverFunc returns a type satisfying the Observer interface, for which the
// PopulationUpdate method forwards to f.
func ObserverFunc(f func(*PopulationData)) Observer { return &observerFunc{f: f} }

func (obs *observerFunc) PopulationUpdate(data *PopulationData) { obs.f(data) }
