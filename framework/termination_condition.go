package framework

// TerminationCondition is the interface implemented by objects used to
// terminate evolutionary algorithms.
type TerminationCondition interface {

	// ShouldTerminate reports whether or not evolution should finish at the
	// current point.
	//
	// populationData is the information about the current state of evolution.
	// This may be used to determine whether evolution should continue or not.
	ShouldTerminate(*PopulationData) bool
}
