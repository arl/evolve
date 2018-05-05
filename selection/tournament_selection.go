package selection

import (
	"errors"
	"math/rand"

	"github.com/aurelien-rainone/evolve/framework"
)

// ErrInvalidTournamentProb is the error returned when trying to set an invalid
// tournament selection probability
var ErrInvalidTournamentProb = errors.New("crossover probability must be in the [0.0,1.0] range")

// TournamentSelection is a selection strategy that picks a pair of candidates
// at random and then selects the fitter of the two candidates with probability
// p, where p is the configured selection probability (therefore the probability
// of the less fit candidate being selected is 1 - p).
type TournamentSelection struct {
	prob             float64
	varprob          bool
	probmin, probmax float64
	description      string
}

// NewTournamentSelection creates a TournamentSelection selection strategy where
// the probability of selecting the fitter of two randomly chosen candidates is
// set to 0.7.
func NewTournamentSelection() TournamentSelection {
	// create with a selection probability of 0.7
	return TournamentSelection{
		prob: 0.7, varprob: false, probmin: 0.7, probmax: 0.7,
	}
}

// SetProb sets a constant probability that fitter of two randomly chosen
// candidates will be selected.
//
// The probability of selecting the fitter of two candidates must be greater
// than 0.5 to be useful (if it is not, there is no selection pressure, or the
// pressure is in favour of weaker candidates, which is counter-productive).
// If prob is not in the (0.5,1.0] range SetProb will return
// ErrInvalidTournamentProb
func (ts TournamentSelection) SetProb(prob float64) error {
	if prob < 0.5 || prob > 1.0 {
		return ErrInvalidTournamentProb
	}
	op.prob = prob
	op.varprob = false
	return nil
}

// SetProbRange sets the range of possible tournament selection probabilities.
//
// The specific probability will be randomly chosen with the pseudo random
// number generator argument of Apply, by linearly converting from (0.5,1.0) to
// [min,max).
//
// If min and max are not bounded by (0.5,1.0] SetProbRange will return
// ErrInvalidTournamentProb.
func (ts TournamentSelection) SetProbRange(min, max float64) error {
	if min > max || min < 0.5 || max > 1.0 {
		return ErrInvalidTournamentProb
	}
	op.probmin = min
	op.probmax = max
	op.varprob = true
	return nil
}

// Select selects the specified number of candidates from the population.
func (ts *TournamentSelection) Select(
	population framework.EvaluatedPopulation,
	naturalFitnessScores bool,
	selectionSize int,
	rng *rand.Rand) []framework.Candidate {

	selection := make([]framework.Candidate, selectionSize)
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

func (ts *TournamentSelection) String() string { return "Tournament Selection" }
