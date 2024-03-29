package main

import (
	"math/rand"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/arl/evolve/generator"
	"github.com/arl/evolve/pkg/mt19937"

	"github.com/stretchr/testify/require"
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
	require.NoError(t, err)

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
	require.NoError(t, err)

	mater{}.Mate(p1, p2, 1, rand.New(mt19937.New(2)))
}

// Tests to ensure that rows are still valid after mutation.  Each row
// should contain each value 1-9 exactly once.
func TestRowMutationValidity(t *testing.T) {
	rng := rand.New(mt19937.New(time.Now().UnixNano()))

	rmut := &rowMutation{
		Number: generator.Const[uint](8),
		Amount: generator.Const[uint](1),
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
	require.NoError(t, err)

	pop := []*sudoku{sudo}

	counts := make(map[int]struct{})

	for i := 0; i < 20; i++ {
		pop = rmut.Apply(pop, rng)
		require.Len(t, pop, 1, "population size should not be affected by mutation")

		mutated := pop[0]
		for j := 0; j < size; j++ {
			row := mutated[j]
			require.Lenf(t, row, size, "row %v has an invalid length", j)

			for _, cell := range row {
				if cell.val <= 0 || cell.val > size {
					t.Errorf("on row %v cell value is out of range, got %v", j, cell.val)
				}
				counts[cell.val] = struct{}{}
			}
			require.Lenf(t, counts, size, "row %v contains some duplicated values", j)

			// Clear map
			counts = make(map[int]struct{})
		}
	}
}

// Check that the mutation never modifies the value of fixed cells.
func TestRowMutationFixedConstraints(t *testing.T) { // nolint: gocyclo
	rmut := &rowMutation{
		Number: generator.Const[uint](8),
		Amount: generator.Const[uint](1),
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
	pop := []*sudoku{&sudo}
	for i := 0; i < 100; i++ { // 100 generations of mutation.
		pop = rmut.Apply(pop, rng)
		mutated := pop[0]
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
