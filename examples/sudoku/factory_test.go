package main

import (
	"math/rand"
	"testing"
	"time"
)

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
	population := factory.GenPopulation(20, rng)
	for _, iface := range population {
		sudo := iface.(*sudoku)

		// Check givens are correctly placed.
		if !sudo[2][8].fixed {
			t.Error("Cell (2, 8) should be fixed.")
		}
		if sudo[2][8].val != 5 {
			t.Error("Cell (2, 8) should contain 5.")
		}
		if !sudo[7][3].fixed {
			t.Error("Cell (7, 3) should be fixed.")
		}
		if sudo[7][3].val != 1 {
			t.Error("Cell (7, 3) should contain 1.")
		}
		if !sudo[3][4].fixed {
			t.Error("Cell (3, 4) should be fixed.")
		}
		if sudo[3][4].val != 2 {
			t.Error("Cell (3, 4) should contain 2.")
		}
		if !sudo[0][1].fixed {
			t.Error("Cell (0, 1) should be fixed.")
		}
		if sudo[0][1].val != 9 {
			t.Error("Cell (0, 1) should contain 9.")
		}
		if !sudo[8][8].fixed {
			t.Error("Cell (8, 8) should be fixed.")
		}
		if sudo[8][8].val != 9 {
			t.Error("Cell (8, 8) should contain 9.")
		}

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
			ErrPatternUnexpectedChar,
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
			ErrWrongNumberOfRows,
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
			ErrWrongNumberOfCols,
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
