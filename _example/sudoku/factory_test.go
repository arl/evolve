package main

import (
	"errors"
	"math/rand"
	"testing"
	"time"

	"github.com/arl/evolve"
	"github.com/arl/evolve/pkg/set"
)

func checkCellVal(t *testing.T, s *sudoku, i, j, want int) {
	t.Helper()

	if s[i][j].val != want {
		t.Errorf("Cell (%d, %d) value, want %v got %v", i, j, want, s[i][j].val)
	}
}

func checkCellFixed(t *testing.T, s *sudoku, i, j int) {
	t.Helper()

	if !s[i][j].fixed {
		t.Errorf("Cell (%d, %d) fixed, want true, got false", i, j)
	}
}

// Checks to make sure that the givens are correctly placed and that each row
// contains each value exactly once.
func TestGeneratorValidity(t *testing.T) {
	gen, err := newFactory([]string{
		".9.......",
		".........",
		"........5",
		"....2....",
		".........",
		".........",
		".........",
		"...1.....",
		"........9",
	})
	if err != nil {
		t.Fatal(err)
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	pop := evolve.GeneratePopulation[*sudoku](20, gen, nil, rng)
	for i := 0; i < pop.Len(); i++ {
		sudo := pop.Candidates[i]

		// Check givens are correctly placed.
		checkCellFixed(t, sudo, 2, 8)
		checkCellVal(t, sudo, 2, 8, 5)

		checkCellFixed(t, sudo, 7, 3)
		checkCellVal(t, sudo, 7, 3, 1)

		checkCellFixed(t, sudo, 3, 4)
		checkCellVal(t, sudo, 3, 4, 2)

		checkCellFixed(t, sudo, 0, 1)
		checkCellVal(t, sudo, 0, 1, 9)

		checkCellFixed(t, sudo, 8, 8)
		checkCellVal(t, sudo, 8, 8, 9)

		// Check that each row has no duplicates.
		counts := set.NewOf[int]()
		for i := 0; i < 9; i++ {
			row := sudo[i]
			for _, cell := range row {
				counts.Insert(cell.val)
			}
			if counts.Len() != 9 {
				t.Errorf("row %d has some duplicated values", 9)
			}
		}
	}
}

func TestGeneratorInvalidPatterns(t *testing.T) {
	tests := []struct {
		name    string
		pattern []string
		wantErr error
	}{
		{
			"invalid character",
			[]string{
				"....9....",
				"2..3.....",
				"........1",
				"....a....", // Invalid character on this line.
				"....4....",
				".........",
				".........",
				".........",
				".........",
			},
			errPatternUnexpectedChar,
		},
		{
			"wrong number of rows",
			[]string{
				"....9....",
				"2..3.....",
				"........1",
				".........",
				".........",
			},
			errWrongNumberOfRows,
		},
		{
			"wrong number of columns",
			[]string{
				"....9....",
				"2..3.....",
				"........1",
				".........",
				".........7",
				".........",
				".4.......6",
				"..1.3....",
				"........8",
			},
			errWrongNumberOfCols,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := newFactory(tt.pattern); !errors.Is(err, tt.wantErr) {
				t.Errorf("newFactory returned %v, want %v", err, tt.wantErr)
			}
		})
	}
}
