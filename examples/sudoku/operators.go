package main

import (
	"errors"
	"math"
	"math/rand"
)

// crossover performs crossover between Sudoku grids by re-combining rows from
// parents to form new offspring.  Rows are copied intact, only columns are
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

var errInvalidMutationCount = errors.New("mutation count must be greater than 0")
var errInvalidMutationAmount = errors.New("mutation amount must be greater than 0")

// rowMutation rows in a potential Sudoku solution by manipulating the order
// of non-fixed cells in much the same way as the ListOrderMutation
// operator does in the TSP example. TODO: add TSP example
//
// nmut represents the number of mutation to apply to each row in a sudoku
// solution candidate
// amnt (mutation amount) is the number of positions by which a list element
// will be displaced as a result of mutation
type rowMutation struct {
	nmut             int // number of mutations
	varnmut          bool
	nmutmin, nmutmax int

	amnt             int  // mutation amount
	varamnt          bool // variable amount
	amntmin, amntmax int  // min/max for variable amount

	// These look-up tables keep track of which values are fixed in which columns
	// and sub-grids.  Because the values are fixed, they are the same for all
	// potential solutions, so we cache the information here to minimise the amount
	// of processing that needs to be done for each mutation.  There is no need to
	// worry about rows since the mutation ensures that rows are always valid.
	fixedcols     [size][size]bool
	fixedsubgrids [size][size]bool
	cached        bool
}

func newRowMutation() *rowMutation {
	// default row mutation is 1 mutation per candidate
	return &rowMutation{
		nmut: 1, varnmut: false, nmutmin: 1, nmutmax: 1,
		amnt: 1, varamnt: false, amntmin: 1, amntmax: 1,
	}
}

func (rm *rowMutation) SetMutations(nmut int) error {
	if nmut < 1 || nmut > math.MaxInt32 {
		return errInvalidMutationCount
	}
	rm.nmut = nmut
	rm.varnmut = false
	return nil
}

func (rm *rowMutation) SetMutationsRange(min, max int) error {
	if min > max || min < 1 || max > math.MaxInt32 {
		return errInvalidMutationCount
	}
	rm.nmutmin = min
	rm.nmutmax = max
	rm.varnmut = true
	return nil
}

func (rm *rowMutation) SetAmount(amnt int) error {
	if amnt < 1 || amnt > math.MaxInt32 {
		return errInvalidMutationAmount
	}
	rm.amnt = amnt
	rm.varamnt = false
	return nil
}

func (rm *rowMutation) SetAmountRange(min, max int) error {
	if min > max || min < 1 || max > math.MaxInt32 {
		return errInvalidMutationAmount
	}
	rm.amntmin = min
	rm.amntmax = max
	rm.varamnt = true
	return nil
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

	// get/decide the number of mutations for this mutation
	nmut := rm.nmut
	if rm.varnmut {
		nmut = rm.nmutmin + rng.Intn(rm.nmutmax-rm.nmutmin)
	}
	//int mutationCount = Math.abs(mutationCountVariable.nextValue());
	for nmut > 0 {
		row := rng.Intn(size)
		fromIndex := rng.Intn(size)

		// get/decide the amount for this mutation
		amnt := rm.amnt
		if rm.varamnt {
			amnt = rm.amntmin + rng.Intn(rm.amntmax-rm.amntmin)
		}
		toIndex := (fromIndex + amnt) % size

		// Make sure we're not trying to mutate a 'given'.
		if !newrows[row][fromIndex].fixed && !newrows[row][toIndex].fixed &&
			// ...or trying to introduce a duplicate of a given value.
			(!rm.isAddConflict(sudo, row, fromIndex, toIndex) ||
				rm.isRemoveConflict(sudo, row, fromIndex, toIndex)) {
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
