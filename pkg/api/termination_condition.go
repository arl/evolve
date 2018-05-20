package api

import "fmt"

// TerminationCondition is the interface implemented by objects used to
// terminate evolutionary algorithms.
//
// TODO: think about a clean way to refactor (and rename) the concept of
// termination condition. It seems it would be more go idiomatic to call it
// cancellation and use either Context, channels or both. ElpasedTime and
// UserAbort termination conditions could already be implemented with a Context.
// maybe a evolution Context implementation that proposes, th original Context
// inteface plus WithMaxGeneration(), WithFitness() etc. that is passed to
// evolution engine?
type TerminationCondition interface {
	fmt.Stringer

	// ShouldTerminate reports whether or not evolution should finish at the
	// current point.
	//
	// pdata is the information about the current state of evolution.  This may
	// be used to determine whether evolution should continue or not.
	ShouldTerminate(pdata *PopulationData) bool
}
