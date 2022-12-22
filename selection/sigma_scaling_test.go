package selection

import (
	"testing"
)

func TestSigmaScaling(t *testing.T) {
	t.Run("natural", testFitnessBasedSelection(&SigmaScaling[string]{}, fitnessBasedPopNatural, true))
	t.Run("non-natural", testFitnessBasedSelection(&SigmaScaling[string]{}, fitnessBasedPopNonNatural, false))

	// If all fitness scores are equal, standard deviation is zero.
	t.Run("no-variance/natural", testFitnessBasedSelection(&SigmaScaling[string]{}, fitnessBasedPopAllEqual, true))
	t.Run("no-variance/non-natural", testFitnessBasedSelection(&SigmaScaling[string]{}, fitnessBasedPopAllEqual, false))
}
