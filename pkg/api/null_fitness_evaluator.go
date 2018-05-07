package api

// NullFitnessEvaluator is a fitness evaluator that always return a score of 0.
//
// Fitness evaluation is not required for interactive selection, so this stub
// implementation is used to satisfy the framework requirements.
type NullFitnessEvaluator struct{}

// Fitness returns a score of zero, regardless of the candidate being
// evaluated.
func (NullFitnessEvaluator) Fitness(candidate Candidate, population []Candidate) float64 { return 0 }

// IsNatural always returns true. However, the return value of this method is
// irrelevant since no meaningful fitness scores are produced.
func (NullFitnessEvaluator) IsNatural() bool { return true }
