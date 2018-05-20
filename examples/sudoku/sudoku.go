package main

import (
	"bytes"
	"fmt"
	"log"
	"strings"
)

const size = 9 // dimension of puzzle square

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
		_, err := fmt.Fprintln(&buf, strings.Join(row, " "))
		if err != nil {
			log.Fatalf("can't write to buf: %v", err)
		}
	}
	return buf.String()
}

// The fitness score for a potential Sudoku solution is the number of cells
// conflicting with other cells in the grid (i.e. if there are two 7s in the
// same column, both of these cells are conflicting). A lower score indicates a
// fitter individual.
type evaluator struct{}

func (evaluator) Fitness(cand interface{}, pop []interface{}) float64 { // nolint: golint
	// We can assume that there are no duplicates in any rows because the
	// candidate factory and evolutionary operators that we use do not permit
	// rows to contain duplicates.
	var fitness int
	sudo := cand.(*sudoku)

	// Check columns for duplicates.
	values := make(map[int]struct{})
	for col := 0; col < size; col++ {
		for row := 0; row < size; row++ {
			values[sudo[row][col].val] = struct{}{}
		}
		fitness += size - len(values)
		values = make(map[int]struct{})
	}

	// Check sub-grids for duplicates.
	for band := 0; band < size; band += 3 {
		for stack := 0; stack < size; stack += 3 {
			for row := band; row < band+3; row++ {
				for col := stack; col < stack+3; col++ {
					values[sudo[row][col].val] = struct{}{}
				}
			}
			fitness += size - len(values)
			values = make(map[int]struct{})
		}
	}
	return float64(fitness)
}

func (evaluator) IsNatural() bool { return false }
