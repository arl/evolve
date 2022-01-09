package engine

import (
	"github.com/arl/bitstring"
	"github.com/arl/evolve"
	"github.com/arl/evolve/factory"
)

// Create an engine evolving bit strings, in which the generator simply counts
// the number of ones.
// See full example in "evolve/example/bits/main.go"
func ExampleNew() {
	eng := Engine[*bitstring.Bitstring]{
		Factory: factory.Bitstring(2),
		Epocher: &Generational[*bitstring.Bitstring]{},
		Evaluator: evolve.EvaluatorFunc(
			true, // natural fitness (higher is better)
			func(cand *bitstring.Bitstring, pop []*bitstring.Bitstring) float64 {
				// our evaluator counts the ones in the bitstring
				return float64(cand.OnesCount())
			}),
	}
	eng.Evolve(100)
}
