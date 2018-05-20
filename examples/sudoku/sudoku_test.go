package main

import (
	"testing"
)

func TestSudokuEvaluator(t *testing.T) {
	tests := []struct {
		cand []string // string representation
		want float64  // fitness
	}{
		{
			[]string{"4 1 5 2 7 9 3 6 8",
				"8 2 3 4 5 6 7 9 1",
				"6 7 9 1 3 8 2 4 5",
				"1 3 2 5 4 7 6 8 9",
				"5 4 6 8 9 2 1 3 7",
				"7 9 8 3 6 1 4 5 2",
				"2 5 1 6 8 3 9 7 4",
				"3 8 7 9 1 4 5 2 6",
				"9 6 4 7 2 5 8 1 3",
			},
			0.0,
		},
		{
			[]string{

				"4 1 5 2 7 9 3 8 6",
				"8 2 3 4 5 6 7 9 1",
				"6 7 9 1 3 8 2 4 5",
				"1 3 2 5 4 7 6 8 9",
				"5 4 6 8 9 2 1 3 7",
				"7 9 8 3 6 1 4 5 2",
				"2 5 1 6 8 3 9 7 4",
				"3 8 7 9 1 4 5 2 6",
				"9 6 4 7 2 5 8 1 3",
			},
			2.0,
		},
		{
			[]string{

				"4 1 5 2 7 9 3 8 6",
				"8 2 3 4 5 6 7 9 1",
				"6 7 9 1 3 8 2 4 5",
				"1 3 2 5 4 7 6 8 9",
				"5 4 6 8 9 2 1 3 7",
				"7 9 8 3 6 1 4 5 2",
				"2 5 6 1 8 3 9 7 4",
				"3 8 7 9 1 4 5 2 6",
				"9 6 4 7 2 5 8 1 3",
			},
			6.0,
		},
	}

	for _, tt := range tests {
		sud, err := sudokuFromStrings(tt.cand)
		if err != nil {
			t.Errorf("error creating sudoku from string: %v", err)
		}
		ev := evaluator{}
		if got := ev.Fitness(sud, nil); got != tt.want {
			t.Errorf("evaluator.Fitness() = %v, want %v", got, tt.want)
		}
	}
}
