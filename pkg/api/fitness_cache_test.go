package api

import "testing"

// incrEvaluator breaks the rules for the caching evaluator in that it is not
// repeatable (it returns different values when invoked multiple times for the
// same candidate), but it allows us to see whether we are getting a cached
// value or a new value.
type incrEvaluator struct {
	natural bool
	count   int
}

func (ie *incrEvaluator) Fitness(cand interface{}, pop []interface{}) float64 {
	ie.count++
	return float64(ie.count)
}

func (ie *incrEvaluator) IsNatural() bool { return ie.natural }

func TestFitnessCacheMiss(t *testing.T) {
	eval := FitnessCache{Wrapped: &incrEvaluator{natural: true}}
	fitness := eval.Fitness("Test1", nil)
	if fitness != 1 {
		t.Errorf("wrong fitness, want 1, got %v", fitness)
	}
	// Different candidate so shouldn't return a cached value.
	fitness = eval.Fitness("Test2", nil)

	if fitness != 2 {
		t.Errorf("wrong fitness, want 2, got %v", fitness)
	}
}

func TestFitnessCacheHit(t *testing.T) {
	eval := FitnessCache{Wrapped: &incrEvaluator{natural: true}}
	cand := "Test"
	fitness := eval.Fitness(cand, nil)
	if fitness != 1 {
		t.Errorf("wrong fitness, want 1, got %v", fitness)
	}

	fitness = eval.Fitness(cand, nil)
	// If the value is found in the cache it won't have changed. If it is
	// recalculated, it will have.
	if fitness != 1 {
		t.Errorf("wrong fitness, want 1 (cached), got %v", fitness)
	}
}

func TestFitnessCacheNaturalness(t *testing.T) {
	eval := FitnessCache{Wrapped: &incrEvaluator{natural: true}}
	if !eval.IsNatural() {
		t.Errorf("fitness cache should be natural if wrapped is natural")
	}
	eval = FitnessCache{Wrapped: &incrEvaluator{natural: false}}
	if eval.IsNatural() {
		t.Errorf("fitness cache should not be natural if wrapped is not natural")
	}
}
