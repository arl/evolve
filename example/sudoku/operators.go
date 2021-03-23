package main

import (
	"math/rand"

	"github.com/arl/evolve/generator"
)

// A mater performs crossover between Sudoku grids by re-combining rows from
// parents to form new offspring. Rows are copied intact, only columns are
// disrupted by this cross-over.
type mater struct{}

// sudokuMater applies crossover vertically on sudoku puzzles square grids
func (m mater) Mate(parent1, parent2 interface{}, nxpts int64, rng *rand.Rand) []interface{} {
	var off1, off2 sudoku
	copy(off1[:], parent1.(*sudoku)[:])
	copy(off2[:], parent2.(*sudoku)[:])

	// Apply as many cross-overs as required.
	for i := int64(0); i < nxpts; i++ {
		// Cross-over index is always greater than zero and less than the length
		// of the parent so that we always pick a point that will result in a
		// meaningful cross-over.
		xidx := (1 + rng.Intn(size-1))
		for j := 0; j < xidx; j++ {
			off1[j], off2[j] = off2[j], off1[j]
		}
	}

	return []interface{}{&off1, &off2}
}

// rowMutation rows in a potential Sudoku solution by manipulating the order of
// non-fixed cells in much the same way as mutation.ListOrder operator does in
// the TSP example.
//
// Number is the number of mutation to apply to each row in a sudoku solution.
// Amount is the number of positions by which a list element will be displaced
// as a result of mutation
type rowMutation struct { // nolint: maligned
	Number generator.Int
	Amount generator.Int

	// These look-up tables keep track of which values are fixed in which
	// columns and sub-grids. Because the values are fixed, they are the same
	// for all potential solutions, so we cache the information here to minimise
	// the amount of processing that needs to be done for each mutation. There
	// is no need to worry about rows since the mutation ensures that rows are
	// always valid.
	fixedcols     [size][size]bool
	fixedsubgrids [size][size]bool
	cached        bool
}

// Apply applies the mutation operator to each entry in the list of selected
// candidates.
func (rm *rowMutation) Apply(sel []interface{}, rng *rand.Rand) []interface{} {
	if !rm.cached {
		rm.buildCache(sel[0].(*sudoku))
	}
	mutpop := make([]interface{}, len(sel))
	for i, cand := range sel {
		mutpop[i] = rm.mutate(cand.(*sudoku), rng)
	}
	return mutpop
}

func (rm *rowMutation) buildCache(sudo *sudoku) {
	for row := 0; row < size; row++ {
		for col := 0; col < size; col++ {
			if sudo[row][col].fixed {
				rm.fixedcols[col][sudo[row][col].val-1] = true
				rm.fixedsubgrids[toSubgrid(row, col)][sudo[row][col].val-1] = true
			}
		}
	}
	rm.cached = true
}

func (rm *rowMutation) mutate(sudo *sudoku, rng *rand.Rand) *sudoku {
	var newrows sudoku
	copy(newrows[:], sudo[:])

	// Find out the number of mutations for this mutation.
	nmut := rm.Number.Next()

	for nmut > 0 {
		row := rng.Intn(size)
		fromIndex := rng.Intn(size)

		// get/decide the amount for this mutation
		amount := int(rm.Amount.Next())
		toIndex := (fromIndex + amount) % size

		// Make sure we're not trying to mutate a 'given'.
		if !newrows[row][fromIndex].fixed && !newrows[row][toIndex].fixed &&
			// ...or trying to introduce a duplicate of a given value.
			(!rm.isAddConflict(sudo, row, fromIndex, toIndex) || rm.isRemoveConflict(sudo, row, fromIndex, toIndex)) {
			// Swap the randomly selected element with the one that is the
			// specified displacement distance away.
			newrows[row][fromIndex], newrows[row][toIndex] = newrows[row][toIndex], newrows[row][fromIndex]
			nmut--
		}
	}

	return &newrows
}

// Checks whether the proposed mutation would introduce a duplicate of a fixed
// value into a column or sub-grid.
func (rm *rowMutation) isAddConflict(sudo *sudoku, row, from, to int) bool {

	fromval, toval := sudo[row][from].val-1, sudo[row][to].val-1

	return rm.fixedcols[from][toval] ||
		rm.fixedcols[to][fromval] ||
		rm.fixedsubgrids[toSubgrid(row, from)][toval] ||
		rm.fixedsubgrids[toSubgrid(row, to)][fromval]
}

// Checks whether the proposed mutation would remove a duplicate of a fixed
// value from a column or sub-grid.
func (rm *rowMutation) isRemoveConflict(sudo *sudoku, row, from, to int) bool {

	fromval, toval := sudo[row][from].val-1, sudo[row][to].val-1

	return rm.fixedcols[from][fromval] ||
		rm.fixedcols[to][toval] ||
		rm.fixedsubgrids[toSubgrid(row, from)][fromval] ||
		rm.fixedsubgrids[toSubgrid(row, to)][toval]
}

// Returns the index of the sub-grid that the specified cells belongs to (a
// number between 0 for topleft and size-1 for bottomright
func toSubgrid(row, col int) int {
	band := row / 3
	stack := col / 3
	return band*3 + stack
}
