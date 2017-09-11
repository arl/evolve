package number

import (
	"fmt"
	"math/rand"
)

/**
 * @param probability The probability value (a number in the range 0..1 inclusive).  A
 * value of zero means that an event is guaranteed not to happen.  A value of 1 means
 * it is guaranteed to occur.
 */
type Probability float64

const (
	/**
	 * Convenient constant representing a probability of one.  An event with
	 * a probability of one is a certainty.
	 * @see #ZERO
	 * @see #EVENS
	 */

	ProbabilityOne Probability = 1
	/**
	 * Convenient constant representing a probability of 0.5 (used to model
	 * an event that has a 50/50 chance of occurring).
	 * @see #ZERO
	 * @see #ONE
	 */

	ProbabilityEven = 0.5
	/**
	 * Convenient constant representing a probability of zero.  If an event has
	 * a probability of zero it will never happen (it is an impossibility).
	 * @see #ONE
	 * @see #EVENS
	 */

	ProbabilityZero = 0
)

func NewProbability(probability float64) (Probability, error) {
	if probability < 0 || probability > 1 {
		return ProbabilityZero, fmt.Errorf("Probability must be in the range 0..1 inclusive, got%v", probability)
	}
	return Probability(probability), nil
}

/**
 * Generates an event according the probability value {@literal p}.
 * @param rng A source of randomness for generating events.
 * @return True with a probability of {@literal p}, false with a probability of
 * {@literal 1 - p}.
 */
func (p Probability) NextEvent(rng *rand.Rand) bool {
	// Don't bother generating an random value if the result is guaranteed.
	return p == 1 || rng.Float64() < float64(p)
}

/**
 * The complement of a probability p is 1 - p.  If p = 0.75, then the complement is 0.25.
 * @return The complement of this probability.
 */
func (p Probability) Complement() Probability {
	return Probability(1 - p)
}
