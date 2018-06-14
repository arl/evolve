package selection

import (
	"errors"
	"fmt"
	"math"
	"math/rand"

	"github.com/arl/evolve"
)

// ErrInvalidTruncRatio is the error returned when trying to set an invalid
// selection ratio for truncation selection
var ErrInvalidTruncRatio = errors.New("truncation selection ratio must be in the (0,1] range")

// Truncation implements the selection of n candidates from a population by
// simply selecting the n candidates with the highest fitness scores (the rest
// is discarded). The same candidate is never selected more than once.
type Truncation struct {
	ratio, minratio, maxratio float64
	varratio                  bool
}

// NewTruncation creates a TruncationSelection configured with the
// provided options.
//
// If no options are provided the selection ratio will vary uniformly between 0
// and 1.
func NewTruncation() *Truncation {
	return &Truncation{varratio: true, minratio: 0.5, maxratio: 1.0}
}

// SetRatio sets a constant selection ratio, that is the proportion of the
// highest ranked candidates to select from the population.
//
// If ratio is not in the (0,1] range SetRatio will return ErrInvalidTruncRatio
func (ts *Truncation) SetRatio(ratio float64) error {
	if ratio <= 0.0 || ratio > 1.0 {
		return ErrInvalidTruncRatio
	}
	ts.ratio = ratio
	ts.varratio = false
	return nil
}

// SetRatioRange sets the range of possible truncation selection ratio.
//
// The specific ratio will be randomly chosen with the pseudo random number
// generator argument of Select, by linearly converting from (0.5,1.0) to
// [min,max).
//
// If min and max are not bounded by [0,1] SetRatioRange will return
// ErrInvalidTruncRatio.
func (ts *Truncation) SetRatioRange(min, max float64) error {
	if min > max || min <= 0.0 || max > 1.0 {
		return ErrInvalidTruncRatio
	}
	ts.minratio = min
	ts.maxratio = max
	ts.varratio = true
	return nil
}

// Select selects the fittest candidates. If the selectionRatio results in
// fewer selected candidates than required, then these candidates are
// selected multiple times to make up the shortfall.
//
// pop is the population of evolved and evaluated candidates from which to
// select.
// natural indicates whether higher fitness values represent fitter individuals
// or not.
// size is the number of candidates to select from the evolved population.
//
// Returns the selected candidates.
func (ts *Truncation) Select(pop evolve.Population, natural bool, size int, rng *rand.Rand) []interface{} {

	sel := make([]interface{}, 0, size)

	// get a random value to decide wether to select the fitter individual
	// or the weaker one.
	ratio := ts.ratio
	if ts.varratio {
		ratio = ts.minratio + (ts.maxratio-ts.minratio)*rng.Float64()
	}

	eligible := int(math.Round(ratio * float64(len(pop))))
	if eligible > size {
		eligible = size
	}

	for {
		count := minint(eligible, size-len(sel))
		for i := 0; i < count; i++ {
			sel = append(sel, pop[i].Candidate)
		}
		if len(sel) >= size {
			break
		}
	}
	return sel
}

func (ts *Truncation) String() string {
	s := "Truncation Selection (%v%%)"
	if ts.varratio {
		return fmt.Sprintf(s, fmt.Sprintf("%5.2f-%5.2f", 100*ts.minratio, 100*ts.maxratio))
	}
	return fmt.Sprintf(s, fmt.Sprintf("%5.2f", 100*ts.ratio))
}

// minint returns the minimum of two int values.
func minint(a, b int) int {
	if a < b {
		return a
	}
	return b
}
