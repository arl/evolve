package evolve

import "fmt"

// Condition is the interface that wraps the IsSatisfied method.
//
// IsSatisfied examines the current state of evolution and decides wether a
// predetermined condition is satisfied.
type Condition interface {
	fmt.Stringer

	// IsSatisfied examines the given population data and returns true if it
	// satisfies some predetermined condition, false otherwise.
	IsSatisfied(pdata *PopulationData) bool
}
