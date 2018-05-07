package api

import "fmt"

// TerminationCondition is the interface implemented by objects used to
// terminate evolutionary algorithms.
type TerminationCondition interface {
	fmt.Stringer

	// ShouldTerminate reports whether or not evolution should finish at the
	// current point.
	//
	// populationData is the information about the current state of evolution.
	// This may be used to determine whether evolution should continue or not.
	ShouldTerminate(*PopulationData) bool
}
