package main

import (
	"math/rand"
	"testing"
	"time"
)

func checkCellVal(t *testing.T, s *sudoku, i, j, want int) {
	t.Helper()

	got := s[i][j].val
	if got != want {
		t.Errorf("Cell (%d, %d) value, want %v got %v", i, j, want, got)
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
func TestFactoryValidity(t *testing.T) {
	factory, err := newSudokuFactory([]string{
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
		t.Errorf("can't create factory: %v", err)
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	pop := factory.GenPopulation(20, rng)
	for _, iface := range pop {
		sudo := iface.(*sudoku)

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
		set := make(map[int]struct{})
		for i := 0; i < 9; i++ {
			row := sudo[i]
			for _, cell := range row {
				set[cell.val] = struct{}{}
			}
			if len(set) < 9 {
				t.Errorf("in\n%v\nrow %v contains duplicates", sudo, i)
			}
		}
	}
}

func TestFactoryInvalidPatterns(t *testing.T) {
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
			_, err := newSudokuFactory(tt.pattern)
			if err != tt.wantErr {
				t.Fatalf("newFactory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
