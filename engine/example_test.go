package engine

import (
	"github.com/arl/evolve"
	"github.com/arl/evolve/generator"

	"github.com/arl/evolve/pkg/bitstring"
)

// Create an engine evolving bit strings, in which the generator simply counts
// the number of ones.
// See full example in "evolve/example/bits/main.go"
func ExampleNew() {
	var epocher Generational

	gen := generator.Bitstring(2)

	eval := evolve.EvaluatorFunc(
		true, // natural fitness (higher is better)
		func(cand interface{}, pop []interface{}) float64 {
			// our evaluator counts the ones in the bitstring
			return float64(cand.(*bitstring.Bitstring).OnesCount())
		})

	New(gen, eval, &epocher)
}