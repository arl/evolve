package selection

import (
	"errors"
	"fmt"
	"math/rand"

	"github.com/aurelien-rainone/evolve/base"
	"github.com/aurelien-rainone/evolve/number"
)

// TournamentSelection is a selection strategy that picks a pair of candidates
// at random and then selects the fitter of the two candidates with probability
// p, where p is the configured selection probability (therefore the probability
// of the less fit candidate being selected is 1 - p).
type TournamentSelection struct {
	selectionProbability number.ProbabilityGenerator
	description          string
}

// TournamentSelectionOption is the type of the functions used to set tournament
// selection options.
type TournamentSelectionOption func(*TournamentSelection) error

// WithConstantSelectionProbability sets up a constant probability that the
// fitter of two randomly chosen candidates will be selected.
func WithConstantSelectionProbability(selectionProbability number.Probability) TournamentSelectionOption {
	return func(ts *TournamentSelection) error {
		if selectionProbability <= 0.5 {
			return errors.New("selection threshold must be greater than 0.5")
		}
		ts.selectionProbability = number.NewConstantProbabilityGenerator(selectionProbability)
		ts.description = fmt.Sprintf("Tournament Selection (p = %v)", selectionProbability)
		return nil
	}
}

// WithVariableSelectionProbability sets up a variable probability that the
// fittest candidate is being selected in any given tournament.
//
// variable should be a probability generator that produce values in the range
// [0.5, 1]. These values are used as the probability of the fittest candidate
// being selected in any given tournament.
func WithVariableSelectionProbability(variable number.ProbabilityGenerator) TournamentSelectionOption {
	return func(ts *TournamentSelection) error {
		ts.selectionProbability = variable
		ts.description = "Tournament Selection"
		return nil
	}
}

// NewTournamentSelection creates a TournamentSelection configured with provided
// options.
func NewTournamentSelection(options ...TournamentSelectionOption) (*TournamentSelection, error) {
	// create with a selection probability of 0.5
	ts := &TournamentSelection{
		selectionProbability: number.NewConstantProbabilityGenerator(number.ProbabilityEven),
	}

	// set client options
	for _, option := range options {
		if err := option(ts); err != nil {
			return nil, fmt.Errorf("can't apply tournament selection options: %v", err)
		}
	}

	return ts, nil
}

// Select selects the specified number of candidates from the population.
func (ts *TournamentSelection) Select(population []*base.EvaluatedCandidate, naturalFitnessScores bool, selectionSize int, rng *rand.Rand) []base.Candidate {

	selection := make([]base.Candidate, selectionSize)
	for i := 0; i < selectionSize; i++ {
		// Pick two candidates at random.
		candidate1 := population[rng.Intn(len(population))]
		candidate2 := population[rng.Intn(len(population))]

		// Use a random value to decide wether to select the fitter individual or the weaker one.
		selectFitter := ts.selectionProbability.NextValue().NextEvent(rng)
		if selectFitter == naturalFitnessScores {

			// Select the fitter candidate.
			if candidate2.Fitness() > candidate1.Fitness() {
				selection[i] = candidate2.Candidate()
			} else {
				selection[i] = candidate1.Candidate()
			}

		} else {

			// Select the less fit candidate.
			if candidate2.Fitness() > candidate1.Fitness() {
				selection[i] = candidate1.Candidate()
			} else {
				selection[i] = candidate2.Candidate()
			}
		}
	}
	return selection
}

func (ts *TournamentSelection) String() string {
	return ts.description
}
