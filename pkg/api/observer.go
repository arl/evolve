package api

// Observer is the interface that wraps the PopulationUpdate method. Observers
// give the opportunity to monitor the state of evolutionary algorithms.
//
// Depending on the parameters of the evolutionary program, an observer may
// be invoked dozens or hundreds of times a second, especially when the population
// size is small as this leads to shorter generations. The processing performed by an
// evolution observer should be reasonably short-lived so as to avoid slowing down
// the evolution.
type Observer interface {

	// PopulationUpdate is called when the state of the population has changed
	// (typically at the end of a generation) with data about the current
	// generation.
	PopulationUpdate(data *PopulationData)
}
