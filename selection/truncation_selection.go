package selection

import (
	"errors"
	"fmt"
	"math/rand"

	"github.com/aurelien-rainone/evolve/framework"
	"github.com/aurelien-rainone/evolve/number"
)

// TruncationSelection implements selection of n candidates from a population by
// simply selecting the n candidates with the highest fitness scores (the rest
// are discarded). A candidate is never selected more than once.
type TruncationSelection struct {
	//private static final DecimalFormat PERCENT_FORMAT = new DecimalFormat("#0.###%");
	selectionRatio number.Float64Generator

	description string
}

// TruncationSelectionOption is the type of functions used to specify options during
// the creation of TruncationSelection objects.
type TruncationSelectionOption func(*TruncationSelection) error

// NewTruncationSelection creates a TruncationSelection configured with the
// provided options.
//
// If no options are provided the selection ratio will vary uniformly between 0
// and 1.
func NewTruncationSelection(options ...TruncationSelectionOption) (*TruncationSelection, error) {
	sel := &TruncationSelection{
		selectionRatio: number.NewBoundedFloat64Generator(0, 1),
		description:    "Truncation Selection",
	}

	// set client options
	for _, option := range options {
		if err := option(sel); err != nil {
			return nil, fmt.Errorf("can't apply truncation selection options: %v", err)
		}
	}
	return sel, nil

}

// WithVariableSelectionRatio sets up a variable selection ratio provided by the
// specified number.Float64Generator.
//
// variable is a number generator that produce values in the range
// "0 < r < 1". These values are used to determine the proportion of the
// population that is retained in any given selection.
func WithVariableSelectionRatio(variable number.Float64Generator) TruncationSelectionOption {
	return func(sel *TruncationSelection) error {
		sel.selectionRatio = variable
		return nil
	}
}

// WithConstantSelectionRatio sets up the selection ration, that is the
// proportion of the highest ranked candidates to selection from the population.
//
// ratio must be positive and less than 1.
func WithConstantSelectionRatio(ratio float64) TruncationSelectionOption {
	return func(sel *TruncationSelection) error {
		if ratio <= 0 || ratio >= 1 {
			return errors.New("selection ratio must be positive and less than 1")
		}
		sel.selectionRatio = number.NewConstantFloat64Generator(ratio)
		sel.description = "Truncation Selection (" + fmt.Sprintf("%.3f", ratio) + ")"
		return nil
	}
}

// Select selects the fittest candidates. If the selectionRatio results in
// fewer selected candidates than required, then these candidates are
// selected multiple times to make up the shortfall.
//
// - population is the population of evolved and evaluated candidates
// from which to select.
// - naturalFitnessScores indicates whether higher fitness values represent fitter
// individuals or not.
// - selectionSize The number of candidates to select from the
// evolved population.
//
// Returns the selected candidates.
func (sel TruncationSelection) Select(
	population framework.EvaluatedPopulation,
	naturalFitnessScores bool,
	selectionSize int,
	rng *rand.Rand) []framework.Candidate {
	selection := make([]framework.Candidate, 0, selectionSize)

	ratio := sel.selectionRatio.NextValue()
	if ratio < 0 || ratio > 1 {
		panic(fmt.Sprintln("Selection ratio out-of-range:", ratio))
	}

	eligibleCount := round(ratio * float64(len(population)))
	if eligibleCount > selectionSize {
		eligibleCount = selectionSize
	}

	for {
		count := minint(eligibleCount, selectionSize-len(selection))
		for i := 0; i < count; i++ {
			selection = append(selection, population[i].Candidate())
		}
		if len(selection) >= selectionSize {
			break
		}
	}
	return selection
}

func (sel *TruncationSelection) String() string {
	return sel.description
}

// round rounds floats into integer numbers.
// FIXME: Remove this function when math.Round will exist (in Go 1.10)
func round(a float64) int {
	if a < 0 {
		return int(a - 0.5)
	}
	return int(a + 0.5)
}

// minint returns the minimum of two int values.
func minint(a, b int) int {
	if a < b {
		return a
	}
	return b
}
