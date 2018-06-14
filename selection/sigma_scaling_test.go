package selection

import (
	"testing"
)

func TestSigmaScalingNaturalFitness(t *testing.T) {
	testFitnessBasedSelection(t, SigmaScaling, fitnessBasedPopNatural, true)
}

func TestSigmaScalingNonNaturalFitness(t *testing.T) {
	testFitnessBasedSelection(t, SigmaScaling, fitnessBasedPopNonNatural, false)
}

// If all fitness scores are equal, standard deviation is zero.
func TestSigmaScalingNoVariance(t *testing.T) {
	testFitnessBasedSelection(t, SigmaScaling, fitnessBasedPopAllEqual, true)
	testFitnessBasedSelection(t, SigmaScaling, fitnessBasedPopAllEqual, false)
}
