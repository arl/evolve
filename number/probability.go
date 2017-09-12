package number

import (
	"fmt"
	"math/rand"
)

// Probabilty represents a probability value (a number in the range 0..1 inclusive).
//
// A value of zero means that an event is guaranteed not to happen. A value of
// 1 means it is guaranteed to occur.
type Probability float64

var (
	// ProbabilityOne is a convenient constant representing a probability of
	// one, that is an event with a probability of one is a certainty.
	ProbabilityOne Probability

	// ProbabilityEven is a convenient constant representing a probability of
	// 0.5 (used to model an event that has a 50/50 chance of occurring).
	ProbabilityEven Probability

	// ProbabilityZero is a convenient constant representing a probability of
	// zero, that is an event that will never happen (it is an impossibility).
	ProbabilityZero Probability
)

func init() {
	ProbabilityOne, _ = NewProbability(1)
	ProbabilityEven, _ = NewProbability(0.5)
	ProbabilityZero, _ = NewProbability(0)
}

// NewProbability creates a new Probability value.
func NewProbability(probability float64) (Probability, error) {
	if probability < 0 || probability > 1 {
		return ProbabilityZero, fmt.Errorf("Probability must be in the range 0..1 inclusive, got%v", probability)
	}
	return Probability(probability), nil
}

// NewEvent generates an event according to the probability value.
//
// In other words NextEvent returns True with a probability of p, false with a
// probability of 1 - p.
func (p Probability) NextEvent(rng *rand.Rand) bool {
	// Don't bother generating an random value if the result is guaranteed.
	return p == 1 || rng.Float64() < float64(p)
}

// The Complement of a probability p is 1 - p. For example if p = 0.75, the
// complement is 0.25.
func (p Probability) Complement() Probability {
	return Probability(1 - p)
}
