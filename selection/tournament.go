package selection

import (
	"errors"
	"fmt"
	"math/rand"

	"github.com/arl/evolve"
)

// ErrInvalidTournamentProb is the error returned when trying to set an invalid
// tournament selection probability
var ErrInvalidTournamentProb = errors.New("crossover probability must be in the [0.0,1.0] range")

// Tournament is a selection strategy that picks a pair of candidates at random
// and then selects the fitter of the two candidates with probability p, where p
// is the configured selection probability (therefore the probability of the
// less fit candidate being selected is 1 - p).
type Tournament struct {
	prob             float64
	varprob          bool
	probmin, probmax float64
}

// NewTournament creates a TournamentSelection selection strategy where the
// probability of selecting the fitter of two randomly chosen candidates is set
// to 0.7.
func NewTournament() *Tournament {
	// create with a selection probability of 0.7
	return &Tournament{prob: 0.7, varprob: false, probmin: 0.7, probmax: 0.7}
}

// SetProb sets a constant probability that fitter of two randomly chosen
// candidates will be selected.
//
// The probability of selecting the fitter of two candidates must be greater
// than 0.5 to be useful (if it is not, there is no selection pressure, or the
// pressure is in favour of weaker candidates, which is counter-productive).
// If prob is not in the (0.5,1] range SetProb will return
// ErrInvalidTournamentProb
func (ts *Tournament) SetProb(prob float64) error {
	if prob <= 0.5 || prob > 1.0 {
		return ErrInvalidTournamentProb
	}
	ts.prob = prob
	ts.varprob = false
	return nil
}

// SetProbRange sets the range of possible tournament selection probabilities.
//
// The specific probability will be randomly chosen with the pseudo random
// number generator argument passed to Select, by linearly converting from
// (0.5,1) to [min,max).
//
// If min and max are not bounded by (0.5,1] SetProbRange will return
// ErrInvalidTournamentProb.
func (ts *Tournament) SetProbRange(min, max float64) error {
	if min > max || min < 0.5 || max > 1.0 {
		return ErrInvalidTournamentProb
	}
	ts.probmin = min
	ts.probmax = max
	ts.varprob = true
	return nil
}

// Select selects the specified number of candidates from the population.
func (ts *Tournament) Select(
	pop evolve.Population,
	natural bool,
	size int,
	rng *rand.Rand) []interface{} {

	sel := make([]interface{}, size)
	for i := 0; i < size; i++ {
		// Pick two candidates at random.
		cand1 := pop[rng.Intn(len(pop))]
		cand2 := pop[rng.Intn(len(pop))]

		// get a random value to decide wether to select the fitter individual
		// or the weaker one.
		prob := ts.prob
		if ts.varprob {
			prob = ts.probmin + (ts.probmax-ts.probmin)*rng.Float64()
		}

		if natural && rng.Float64() < prob { // Select the fitter candidate.
			if cand2.Fitness > cand1.Fitness {
				sel[i] = cand2.Candidate
			} else {
				sel[i] = cand1.Candidate
			}
		} else { // Select the less fit candidate.
			if cand2.Fitness > cand1.Fitness {
				sel[i] = cand1.Candidate
			} else {
				sel[i] = cand2.Candidate
			}
		}
	}
	return sel
}

func (ts *Tournament) String() string {
	s := "Tournament Selection (p = %v)"
	if ts.varprob {
		return fmt.Sprintf(s, fmt.Sprintf("[%v,%v]", ts.probmin, ts.probmax))
	}
	return fmt.Sprintf(s, ts.prob)
}
