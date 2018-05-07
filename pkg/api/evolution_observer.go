package api

// EvolutionObserver is a call-back interface so that programs can monitor the
// state of a long-running evolutionary algorithm.
//
// Depending on the parameters of the evolutionary program, an observer may
// be invoked dozens or hundreds of times a second, especially when the population
// size is small as this leads to shorter generations. The processing performed by an
// evolution observer should be reasonably short-lived so as to avoid slowing down
// the evolution.
type EvolutionObserver interface {
	// PopulationUpdate is invoked when the state of the population has changed
	// (typically at the end of a generation).
	//
	// param contains data statistics about the state of the current generation.
	PopulationUpdate(data *PopulationData)
}
