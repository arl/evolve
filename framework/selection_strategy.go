package framework

import (
	"fmt"
	"math/rand"
)

// SelectionStrategy is the interface that implementd "natural" selection.
type SelectionStrategy interface {
	fmt.Stringer

	// Select selects the specified number of candidates from the population.
	//
	// Implementations may assume that the population is sorted in descending
	// order according to fitness (so the fittest individual is the first item
	// in the list).
	// NOTE: It is an error to call this method with an empty or null population.
	//
	// - population is the population from which to select.
	// naturalFitnessScores indicates whether higher fitness values represent
	// - fitter individuals or not.
	// - selectionSize is the number of individual selections to make (not
	// necessarily the number of distinct candidates to select, since the same
	// individual may potentially be selected more than once).
	//
	// Returns a slice containing the selected candidates. Some individual
	// candidates may potentially have been selected multiple times.
	Select(population []*EvaluatedCandidate, naturalFitnessScores bool, selectionSize int, rng *rand.Rand) []Candidate
}
