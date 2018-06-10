package api

// An Observer monitors the evolution of a population.
//
// Once registered within the evolution engine, observers gets notified of every
// population update, that is, once an epoch is completed.
type Observer interface {

	// PopulationUpdate is called at every population update -once an epoch is
	// has been completed- with information and statistics about the current
	// population.
	PopulationUpdate(data *PopulationData)
}

// TODO: try to come up with a better and short name for PopulationUpdate
// and PopulationData maybe
type observerFunc struct{ f func(*PopulationData) }

// ObserverFunc returns a type satisfying the Observer interface, for which the
// PopulationUpdate method forwards to f.
func ObserverFunc(f func(*PopulationData)) Observer { return &observerFunc{f: f} }

func (obs *observerFunc) PopulationUpdate(data *PopulationData) { obs.f(data) }
