package main

import "github.com/aurelien-rainone/evolve/pkg/api"

type stringEvaluator struct{ targetString string }

func newStringEvaluator(targetString string) stringEvaluator {
	return stringEvaluator{
		targetString: targetString,
	}
}

// Assigns one "fitness point" for every character in the candidate string that
// doesn't match the corresponding position in the target string.
func (se stringEvaluator) Fitness(
	cand api.Candidate,
	pop []api.Candidate) float64 {

	var errors float64
	sc := cand.(string)
	for i := range sc {
		if sc[i] != se.targetString[i] {
			errors++
		}
	}
	return errors
}

func (se stringEvaluator) IsNatural() bool { return false }
