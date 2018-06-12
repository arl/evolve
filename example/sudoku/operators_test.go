package main

import (
	"math/rand"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/arl/evolve/pkg/mt19937"
)

func sudokuFromStrings(strs []string) (*sudoku, error) {
	s := &sudoku{}
	for i, row := range strs {
		vals := strings.Fields(row)
		for j, sval := range vals {
			val, err := strconv.ParseInt(sval, 10, 64)
			if err != nil {
				return nil, err
			}
			s[i][j].val = int(val)
		}
	}
	return s, nil
}

func TestSudokuMater(t *testing.T) {

	p1, err := sudokuFromStrings([]string{
		"1 1 1 1 1 1 1 1 1",
		"2 2 2 2 2 2 2 2 2",
		"3 3 3 3 3 3 3 3 3",
		"4 4 4 4 4 4 4 4 4",
		"5 5 5 5 5 5 5 5 5",
		"6 6 6 6 6 6 6 6 6",
		"7 7 7 7 7 7 7 7 7",
		"8 8 8 8 8 8 8 8 8",
		"9 9 9 9 9 9 9 9 9",
	})
	if err != nil {
		t.Errorf("error creating sudoku from string: %v", err)
	}

	p2, err := sudokuFromStrings([]string{
		"9 9 9 9 9 9 9 9 9",
		"8 8 8 8 8 8 8 8 8",
		"7 7 7 7 7 7 7 7 7",
		"6 6 6 6 6 6 6 6 6",
		"5 5 5 5 5 5 5 5 5",
		"4 4 4 4 4 4 4 4 4",
		"3 3 3 3 3 3 3 3 3",
		"2 2 2 2 2 2 2 2 2",
		"1 1 1 1 1 1 1 1 1",
	})
	if err != nil {
		t.Errorf("error creating sudoku from string: %v", err)
	}

	mater{}.Mate(p1, p2, 1, rand.New(mt19937.New(2)))
}

// Tests to ensure that rows are still valid after mutation.  Each row
// should contain each value 1-9 exactly once.
func TestRowMutationValidity(t *testing.T) { // nolint: gocyclo
	rmut := newRowMutation()
	err := rmut.SetMutations(8)
	if err != nil {
		t.Error(err)
	}
	err = rmut.SetAmount(1)
	if err != nil {
		t.Error(err)
	}
	sudo, err := sudokuFromStrings([]string{
		"1 2 8 5 4 3 9 6 7",
		"7 6 4 9 2 8 5 1 3",
		"3 9 5 7 6 1 2 4 8",
		"6 1 9 4 8 5 7 3 2",
		"5 8 3 6 7 2 1 9 4",
		"4 7 2 3 1 9 8 5 6",
		"8 5 1 2 3 6 4 7 9",
		"9 4 6 8 5 7 3 2 1",
		"2 3 7 1 9 4 6 8 5",
	})
	if err != nil {
		t.Error(err)
	}
	pop := []interface{}{sudo}

	counts := make(map[int]struct{})
	rng := rand.New(mt19937.New(time.Now().UnixNano()))

	for i := 0; i < 20; i++ {
		pop = rmut.Apply(pop, rng)
		if len(pop) != 1 {
			t.Errorf("population size should not be affected by mutation")
		}
		mutated := pop[0].(*sudoku)
		for j := 0; j < size; j++ {
			row := mutated[j]
			if len(row) != size {
				t.Errorf("row %v has an invalid length: want %v, got %v", j, size, len(row))
			}
			for _, cell := range row {
				if cell.val <= 0 || cell.val > size {
					t.Errorf("on row %v cell value is out of range, got %v", j, cell.val)
				}
				counts[cell.val] = struct{}{}
			}
			if len(counts) != size {
				t.Errorf("row %v contains some duplicated values", j)
			}
			// clear map
			counts = make(map[int]struct{})
		}
	}
}

//Check that the mutation never modifies the value of fixed cells.
func TestRowMutationFixedConstraints(t *testing.T) { // nolint: gocyclo
	rmut := newRowMutation()
	err := rmut.SetMutations(8)
	if err != nil {
		t.Error(err)
	}
	err = rmut.SetAmount(1)
	if err != nil {
		t.Error(err)
	}

	var sudo sudoku
	// One cell in each row is fixed (cell 1 in row 1, cell 2 in row 2, etc.)
	for row := 0; row < size; row++ {
		for col := 0; col < size; col++ {
			sudo[row][col].val = col + 1
			sudo[row][col].fixed = col == row
		}
	}
	rng := rand.New(mt19937.New(time.Now().UnixNano()))
	pop := []interface{}{&sudo}
	for i := 0; i < 100; i++ { // 100 generations of mutation.
		pop = rmut.Apply(pop, rng)
		mutated := pop[0].(*sudoku)
		for row := 0; row < size; row++ {
			for col := 0; col < size; col++ {
				if row == col {
					if !mutated[row][col].fixed {
						t.Errorf("fixed cell [%v][%v] has become unfixed", row, col)
					}
					if mutated[row][col].val != row+1 {
						t.Errorf("fixed cell [%v][%v] has changed value", row, col)
					}
				} else {
					if mutated[row][col].fixed {
						t.Errorf("unfixed cell [%v][%v] has become fixed", row, col)
					}
				}
			}
		}
	}
}

func TestRowMutationInvalid(t *testing.T) {
	if err := newRowMutation().SetMutations(0); err != errInvalidMutationCount {
		t.Errorf("SetMutations(0): want ErrInvalidMutationCount, got %v", err)
	}
	if err := newRowMutation().SetAmount(0); err != errInvalidMutationAmount {
		t.Errorf("SetAmount(0): want ErrInvalidMutationAmount, got %v", err)
	}
}
