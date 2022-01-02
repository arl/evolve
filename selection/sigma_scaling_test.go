package selection

import (
	"testing"
)

func TestSigmaScalingNaturalFitness(t *testing.T) {
	testFitnessBasedSelection(t, NewSigmaScaling[string](nil), fitnessBasedPopNatural, true)
}

func TestSigmaScalingNonNaturalFitness(t *testing.T) {
	testFitnessBasedSelection(t, NewSigmaScaling[string](nil), fitnessBasedPopNonNatural, false)
}

// If all fitness scores are equal, standard deviation is zero.
func TestSigmaScalingNoVariance(t *testing.T) {
	testFitnessBasedSelection(t, NewSigmaScaling[string](nil), fitnessBasedPopAllEqual, true)
	testFitnessBasedSelection(t, NewSigmaScaling[string](nil), fitnessBasedPopAllEqual, false)
}
