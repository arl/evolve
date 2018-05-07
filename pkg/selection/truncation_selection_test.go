package selection

import (
	"fmt"
	"math"
	"testing"

	"github.com/aurelien-rainone/evolve/pkg/api"
	"github.com/stretchr/testify/assert"
)

// Unit test for truncation selection strategy ensures the 2 best candidates are
// selected.
func testTruncationSelection(t *testing.T, tpop testPopulation, natural bool) {
	ts := NewTruncation()
	errcheck(t, ts.SetRatio(0.5))
	testRandomBasedSelection(t, ts, tpop, natural, 2,
		func(selected []api.Candidate) error {
			if len(selected) != 2 {
				return fmt.Errorf("want len(selected) == 2, got %v", len(selected))
			}

			sstr := []string{}
			for _, c := range selected {
				sstr = append(sstr, c.(string))
			}

			assert.Contains(t, sstr, tpop[0].name, "best candidate not selected")
			assert.Contains(t, sstr, tpop[1].name, "second best candidate not selected")
			return nil
		})
}

func TestTruncationSelectionNatural(t *testing.T) {
	testTruncationSelection(t, randomBasedPopNatural, true)
}

func TestTruncationSelectionNonNatural(t *testing.T) {
	testTruncationSelection(t, randomBasedPopNonNatural, false)
}

func TestTruncationSelectionSetRatio(t *testing.T) {
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
		if got := NewTruncation().SetRatio(tt.ratio); got != tt.wantErr {
			t.Errorf("SetRatio(%v), got err = %v, wantErr = %v", tt.ratio, got, tt.wantErr)
		}
	}
}

func TestTruncationSelectionSetRatioRange(t *testing.T) {
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
		if got := NewTruncation().SetRatioRange(tt.min, tt.max); got != tt.wantErr {
			t.Errorf("SetRatioRange(%v, %v), got err = %v, wantErr = %v", tt.min, tt.max, got, tt.wantErr)
		}
	}
}
