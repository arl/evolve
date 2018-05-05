package selection

import (
	"math"
	"math/rand"
	"testing"

	"github.com/aurelien-rainone/evolve/framework"
	"github.com/stretchr/testify/assert"
)

// Unit test for truncation selection strategy. Ensures the correct candidates
// are selected.

func TestTruncationSelectionNaturalFitness(t *testing.T) {
	rng := rand.New(rand.NewSource(99))

	ts := NewTruncationSelection()

	errcheck(t, ts.SetRatio(0.5))
	// Higher score is better.
	steve, _ := framework.NewEvaluatedCandidate("Steve", 10.0)
	mary, _ := framework.NewEvaluatedCandidate("Mary", 9.1)
	john, _ := framework.NewEvaluatedCandidate("John", 8.4)
	gary, _ := framework.NewEvaluatedCandidate("Gary", 6.2)

	pop := framework.EvaluatedPopulation{steve, mary, john, gary}

	selection := ts.Select(pop, true, 2, rng)

	assert.Len(t, selection, 2, "want selection size to be 2, got ", len(selection))
	assert.Contains(t, selection, steve.Candidate(), "best candidate not selected")
	assert.Contains(t, selection, mary.Candidate(), "second best candidate not selected")
}

func TestTruncationSelectionNonNaturalFitness(t *testing.T) {
	rng := rand.New(rand.NewSource(99))

	ts := NewTruncationSelection()
	errcheck(t, ts.SetRatio(0.5))

	// Lower score is better.
	gary, _ := framework.NewEvaluatedCandidate("Gary", 6.2)
	john, _ := framework.NewEvaluatedCandidate("John", 8.4)
	mary, _ := framework.NewEvaluatedCandidate("Mary", 9.1)
	steve, _ := framework.NewEvaluatedCandidate("Steve", 10.0)

	pop := framework.EvaluatedPopulation{gary, john, mary, steve}

	selection := ts.Select(pop, false, 2, rng)
	assert.Len(t, selection, 2, "want selection size to be 2, got ", len(selection))

	assert.Contains(t, selection, gary.Candidate(), "best candidate not selected")
	assert.Contains(t, selection, john.Candidate(), "second best candidate not selected")
}

func TestNewTruncationSelectionSetRatio(t *testing.T) {
	tests := []struct {
		ratio   float64
		wantErr error
	}{
		{ratio: -1, wantErr: ErrInvalidTruncRatio},
		{ratio: 0, wantErr: ErrInvalidTruncRatio},
		{ratio: 1.00001, wantErr: ErrInvalidTruncRatio},
		{ratio: math.SmallestNonzeroFloat64, wantErr: nil},
		{ratio: 0.5, wantErr: nil},
		{ratio: 1.0, wantErr: nil},
	}
	for _, tt := range tests {
		if got := NewTruncationSelection().SetRatio(tt.ratio); got != tt.wantErr {
			t.Errorf("SetRatio(%v), got err = %v, wantErr = %v", tt.ratio, got, tt.wantErr)
		}
	}
}

func TestNewTruncationSelectionSetRatioRange(t *testing.T) {
	tests := []struct {
		min, max float64
		wantErr  error
	}{
		{min: 0, max: 1, wantErr: ErrInvalidTruncRatio},
		{min: -1, max: 1, wantErr: ErrInvalidTruncRatio},
		{min: 0, max: 1.00001, wantErr: ErrInvalidTruncRatio},
		{min: 0.2, max: 0.1, wantErr: ErrInvalidTruncRatio},
		{min: 0.1, max: 0.2, wantErr: nil},
		{min: 0.1, max: 1, wantErr: nil},
		{min: 0.1, max: 1.01, wantErr: ErrInvalidTruncRatio},
	}
	for _, tt := range tests {
		if got := NewTruncationSelection().SetRatioRange(tt.min, tt.max); got != tt.wantErr {
			t.Errorf("SetRatioRange(%v, %v), got err = %v, wantErr = %v", tt.min, tt.max, got, tt.wantErr)
		}
	}
}
