package main

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

const (
	size   = 9    // dimension of puzzle square
	minval = 1    // min value of a cell
	maxval = size // max value of a cell
)

// sudoku is a potential solution for a Sudoku puzzle.
type sudoku [size][size]struct {
	val   int  // current value in this cell
	fixed bool // fixed cells are the one provided at the start
}

func (s *sudoku) String() string {
	buf := bytes.Buffer{}
	for i := range s {
		var row []string
		for j := range (*s)[i] {
			row = append(row, fmt.Sprintf("%d", (*s)[i][j].val))
		}
		fmt.Fprintln(&buf, strings.Join(row, " "))
	}
	return buf.String()
}

func newSudoku(strs []string) (*sudoku, error) {
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

// The fitness score for a potential Sudoku solution is the number of cells
// conflicting with other cells in the grid (i.e. if there are two 7s in the
// same column, both of these cells are conflicting). A lower score indicates a
// fitter individual.
type evaluator struct{}

func (ev evaluator) Fitness(cand interface{}, pop []interface{}) float64 { // golint: nolint
	// We can assume that there are no duplicates in any rows because the
	// candidate factory and evolutionary operators that we use do not permit
	// rows to contain duplicates.
	var fitness int
	sud := cand.(*sudoku)

	// Check columns for duplicates.
	values := make(map[int]struct{})
	for col := 0; col < size; col++ {
		for row := 0; row < size; row++ {
			values[sud[row][col].val] = struct{}{}
		}
		fitness += size - len(values)
		values = make(map[int]struct{})
	}

	// Check sub-grids for duplicates.
	for band := 0; band < size; band += 3 {
		for stack := 0; stack < size; stack += 3 {
			for row := band; row < band+3; row++ {
				for col := stack; col < stack+3; col++ {
					values[sud[row][col].val] = struct{}{}
				}
			}
			fitness += size - len(values)
			values = make(map[int]struct{})
		}
	}
	return float64(fitness)
}

func (evaluator) IsNatural() bool { return false }
