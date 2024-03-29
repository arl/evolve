package evolve

import "testing"

// incrEvaluator breaks the rules for the caching evaluator in that it is not
// repeatable (it returns different values when invoked multiple times for the
// same candidate), but it allows us to see whether we are getting a cached
// value or a new value.
type incrEvaluator struct {
	natural bool
	count   int
}

func (ie *incrEvaluator) Fitness(cand int, pop []int) float64 {
	ie.count++
	return float64(ie.count)
}

func (ie *incrEvaluator) IsNatural() bool { return ie.natural }

func TestFitnessCacheMiss(t *testing.T) {
	eval := FitnessCache[int]{Wrapped: &incrEvaluator{natural: true}}
	fitness := eval.Fitness(101, nil)
	if fitness != 1 {
		t.Errorf("wrong fitness, want 1, got %v", fitness)
	}
	// Different candidate so shouldn't return a cached value.
	fitness = eval.Fitness(202, nil)

	if fitness != 2 {
		t.Errorf("wrong fitness, want 2, got %v", fitness)
	}
}

func TestFitnessCacheHit(t *testing.T) {
	eval := FitnessCache[int]{Wrapped: &incrEvaluator{natural: true}}
	fitness := eval.Fitness(101, nil)
	if fitness != 1 {
		t.Errorf("wrong fitness, want 1, got %v", fitness)
	}

	fitness = eval.Fitness(101, nil)
	// If the value is found in the cache it won't have changed. If it is
	// recalculated, it will have.
	if fitness != 1 {
		t.Errorf("wrong fitness, want 1 (cached), got %v", fitness)
	}
}

func TestFitnessCacheNaturalness(t *testing.T) {
	eval := FitnessCache[int]{Wrapped: &incrEvaluator{natural: true}}
	if !eval.IsNatural() {
		t.Errorf("fitness cache should be natural if wrapped is natural")
	}
	eval = FitnessCache[int]{Wrapped: &incrEvaluator{natural: false}}
	if eval.IsNatural() {
		t.Errorf("fitness cache should not be natural if wrapped is not natural")
	}
}
