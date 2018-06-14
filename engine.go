package evolve

import (
	"math/rand"
)

// Epocher is the interface implemented by objects having an Epoch method.
type Epocher interface {

	// Epoch performs one epoch (i.e generation) of the evolutionary process.
	//
	// It takes as argument the population to evolve in that step, the elitism
	// count -that is how many of the fittest candidates are preserved and
	// directly inserted into the nexct generation, without selection- and a
	// source of randomess.
	//
	// It returns the next generation.
	Epoch(Population, int, *rand.Rand) Population
}
