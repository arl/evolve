package xover

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
	"testing"

	"github.com/aurelien-rainone/evolve/pkg/api"
	"github.com/stretchr/testify/assert"
)

func TestCrossover_SetPoints(t *testing.T) {
	tests := []struct {
		npts    int
		wantErr error
	}{
		{npts: -1, wantErr: ErrInvalidXOverNumPoints},
		{npts: 0, wantErr: nil},
		{npts: 1, wantErr: nil},
		{npts: math.MaxInt32, wantErr: nil},
		{npts: math.MaxInt32 + 1, wantErr: ErrInvalidXOverNumPoints},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			op := New(nil)
			if err := op.SetPoints(tt.npts); err != tt.wantErr {
				t.Errorf("Crossover.SetPoints(%v) error = %v, wantErr '%v'", tt.npts, err, tt.wantErr)
			}
		})
	}
}

func TestCrossover_SetPointsRange(t *testing.T) {
	tests := []struct {
		min, max int
		wantErr  error
	}{
		{min: -1, max: 1, wantErr: ErrInvalidXOverNumPoints},
		{min: 1, max: 0, wantErr: ErrInvalidXOverNumPoints},
		{min: 0, max: 0, wantErr: nil},
		{min: 1, max: math.MaxInt32, wantErr: nil},
		{min: 1, max: math.MaxInt32 + 1, wantErr: ErrInvalidXOverNumPoints},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			op := New(nil)
			if err := op.SetPointsRange(tt.min, tt.max); err != tt.wantErr {
				t.Errorf("Crossover.SetPointsRange(%v, %v) error = %v, wantErr '%v'", tt.min, tt.max, err, tt.wantErr)
			}
		})
	}
}

func TestCrossover_SetProb(t *testing.T) {
	tests := []struct {
		prob    float64
		wantErr error
	}{
		{prob: -1.0, wantErr: ErrInvalidXOverProb},
		{prob: 0.0, wantErr: nil},
		{prob: 0.3333, wantErr: nil},
		{prob: 1.0, wantErr: nil},
		{prob: 1.1, wantErr: ErrInvalidXOverProb},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			op := New(nil)
			if err := op.SetProb(tt.prob); err != tt.wantErr {
				t.Errorf("Crossover.SetProb(%v) error = %v, wantErr %v", tt.prob, err, tt.wantErr)
			}
		})
	}
}

func TestCrossover_SetProbRange(t *testing.T) {
	tests := []struct {
		min, max float64
		wantErr  error
	}{
		{min: -1.0, max: 1.0, wantErr: ErrInvalidXOverProb},
		{min: 1.0, max: 0.0, wantErr: ErrInvalidXOverProb},
		{min: 0.0, max: 0.0, wantErr: nil},
		{min: 0.0, max: 1.0, wantErr: nil},
		{min: 1.0, max: 1.1, wantErr: ErrInvalidXOverProb},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			op := New(nil)
			if err := op.SetProbRange(tt.min, tt.max); err != tt.wantErr {
				t.Errorf("Crossover.SetProbRange(%v, %v) error = %v, wantErr '%v'", tt.min, tt.max, err, tt.wantErr)
			}
		})
	}
}

func sameStringPop(t *testing.T, a, b []api.Candidate) {
	t.Helper()

	s1 := make([]string, 0)
	s2 := make([]string, 0)
	for _, a := range a {
		s1 = append(s1, a.(string))
	}
	for _, b := range b {
		s2 = append(s2, b.(string))
	}
	sort.Strings(s1)
	sort.Strings(s2)
	assert.EqualValues(t, s1, s2)
}

func TestCrossover_Apply(t *testing.T) {
	rng := rand.New(rand.NewSource(99))
	pop := []api.Candidate{"abcde", "fghij", "klmno", "pqrst", "uvwxy"}

	t.Run("zero_crossover_points_is_noop", func(t *testing.T) {
		xover := New(StringMater{})
		assert.NoError(t, xover.SetPoints(0))
		got := xover.Apply(pop, rng)
		sameStringPop(t, pop, got)
	})

	t.Run("zero_crossover_probability_is_noop", func(t *testing.T) {
		xover := New(StringMater{})
		assert.NoError(t, xover.SetProb(0.0))
		got := xover.Apply(pop, rng)
		sameStringPop(t, pop, got)
	})
}
