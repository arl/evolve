package selection

import (
	"testing"
)

func TestSigmaScalingNaturalFitness(t *testing.T) {
	testFitnessBasedSelection(t, &SigmaScaling[string]{}, fitnessBasedPopNatural, true)
}

func TestSigmaScalingNonNaturalFitness(t *testing.T) {
	testFitnessBasedSelection(t, &SigmaScaling[string]{}, fitnessBasedPopNonNatural, false)
}

// If all fitness scores are equal, standard deviation is zero.
func TestSigmaScalingNoVariance(t *testing.T) {
	testFitnessBasedSelection(t, &SigmaScaling[string]{}, fitnessBasedPopAllEqual, true)
	testFitnessBasedSelection(t, &SigmaScaling[string]{}, fitnessBasedPopAllEqual, false)
}
