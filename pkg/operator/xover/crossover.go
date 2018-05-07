package xover

import (
	"errors"
	"math"
	"math/rand"

	"github.com/aurelien-rainone/evolve/pkg/api"
)

var (
	// ErrInvalidXOverNumPoints is the error returned when trying to set an
	// invalid number of crossover points
	ErrInvalidXOverNumPoints = errors.New("crossover points must be in the [0,MaxInt32] range")
	// ErrInvalidXOverProb is the error returned when trying to set an invalid
	// probability of crossover
	ErrInvalidXOverProb = errors.New("crossover probability must be in the [0.0,1.0] range")
)

// Mater is the interface implemented by objects defining the Mate function.
type Mater interface {

	// Mate performs crossover on a pair of parents to generate a pair of
	// offspring.
	//
	// parent1 and parent2 are the two individuals that provides the source
	// material for generating offspring.
	// TODO: should return 2 values of a slice of 2 values
	Mate(parent1, parent2 api.Candidate,
		numberOfCrossoverPoints int64,
		rng *rand.Rand) []api.Candidate
}

// Crossover implements a standard crossover operator.
//
// It supports all crossover processes that operate on a pair of parent
// candidates.
// Both the number of crossovers points and the crossover probability are
// configurable. Crossover is applied to a proportion of selected parent pairs,
// with the remainder copied unchanged into the output population. The size of
// this evolved proportion is controlled by the code crossoverProbability
// parameter.
type Crossover struct {
	Mater
	npts             int
	varnpts          bool
	nptsmin, nptsmax int
	prob             float64
	varprob          bool
	probmin, probmax float64
}

// New creates a Crossover operator based off the provided Mater.
//
// The returned Crossover performs a one point crossover with 1.0 (i.e always)
// probabilty.
func New(mater Mater) *Crossover {
	return &Crossover{
		npts: 1, varnpts: false, nptsmin: 1, nptsmax: 1,
		prob: 1.0, varprob: false, probmin: 1.0, probmax: 1.0,
		Mater: mater,
	}
}

// SetPoints sets the number of crossover points.
//
// If npts is not in the [0,MaxInt32] range SetPoints will return
// ErrInvalidXOverNumPoints.
func (op *Crossover) SetPoints(npts int) error {
	if npts < 0 || npts > math.MaxInt32 {
		return ErrInvalidXOverNumPoints
	}
	op.npts = npts
	op.varnpts = false
	return nil
}

// SetProb sets the crossover probability,
//
// If prob is not in the [0.1,1.0] range SetProb will return
// ErrInvalidXOverProb.
func (op *Crossover) SetProb(prob float64) error {
	if prob < 0.0 || prob > 1.0 {
		return ErrInvalidXOverProb
	}
	op.prob = prob
	op.varprob = false
	return nil
}

// SetPointsRange sets the range of possible crossover points.
//
// The specific number of crossover points will be randomly chosen with the
// pseudo random number generator argument of Apply, by linearly converting from
// [0,MaxInt32) to [min,max).
//
// If min and max are not bounded by [0,MaxInt32] SetPointsRange will
// return ErrInvalidXOverNumPoints.
func (op *Crossover) SetPointsRange(min, max int) error {
	if min > max || min < 0 || max > math.MaxInt32 {
		return ErrInvalidXOverNumPoints
	}
	op.nptsmin = min
	op.nptsmax = max
	op.varnpts = true
	return nil
}

// SetProbRange sets the range of possible crossover probabilities.
//
// The specific crossover probability will be randomly chosen with the pseudo
// random number generator argument of Apply, by linearly converting from
// [0.0,1.0) to [min,max).
//
// If min and max are not bounded by [0.0,1.0] SetProbRange will return
// ErrInvalidXOverProb.
func (op *Crossover) SetProbRange(min, max float64) error {
	if min > max || min < 0.0 || max > 1.0 {
		return ErrInvalidXOverProb
	}
	op.probmin = min
	op.probmax = max
	op.varprob = true
	return nil
}

// Apply applies the crossover operation to the selected candidates.
//
// Pairs of candidates are chosen randomly from the selected candidates and
// subjected to crossover to produce a pair of offspring candidates. The
// selected candidates, sel, are the evolved individuals that have survived to
// be eligible to reproduce.
//
// Returns the combined set of evolved offsprings generated by applying
// crossover to the selected candidates.
func (op *Crossover) Apply(sel []api.Candidate, rng *rand.Rand) []api.Candidate {
	// Shuffle the collection before applying each operation so that the
	// evolution is not influenced by any ordering artifacts from previous
	// operations.
	selcopy := make([]api.Candidate, len(sel))
	copy(selcopy, sel)
	api.ShuffleCandidates(selcopy, rng)

	res := make([]api.Candidate, 0, len(sel))
	for i := 0; i < len(selcopy); {
		p1 := selcopy[i]
		i++
		if i < len(selcopy) {
			p2 := selcopy[i]
			i++

			// get/decide a xover probability for this run
			prob := op.prob
			if op.varprob {
				prob = op.probmin + (op.probmax-op.probmin)*rng.Float64()
			}

			var npts int
			if rng.Float64() < prob {
				// we got a crossover to perform, get/decide the number of
				// crossover points
				if op.varnpts {
					npts = op.nptsmin + rng.Intn(op.nptsmax-op.nptsmin)
				} else {
					npts = op.npts
				}
			}
			if npts > 0 {
				res = append(res, op.Mate(p1, p2, int64(npts), rng)...)
			} else {
				// If there is no crossover to perform, just add the parents to the
				// results unaltered.
				res = append(res, p1, p2)
			}
		} else {
			// If we have an odd number of selected candidates, we can't pair up
			// the last one so just leave it unmodified.
			res = append(res, p1)
		}
	}
	return res
}
