package main

import "github.com/aurelien-rainone/evolve/framework"

type stringEvaluator struct {
	targetString string
}

func newStringEvaluator(targetString string) stringEvaluator {
	return stringEvaluator{
		targetString: targetString,
	}
}

// Assigns one "fitness point" for every character in the candidate string that
// doesn't match the corresponding position in the target string.
func (se stringEvaluator) Fitness(
	candidate framework.Candidate,
	population []framework.Candidate) float64 {

	var errors float64
	sc := candidate.(string)
	for i := range sc {
		if sc[i] != se.targetString[i] {
			errors++
		}
	}
	return errors
}

func (se stringEvaluator) IsNatural() bool {
	return false
}
