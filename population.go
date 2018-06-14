package evolve

import (
	"fmt"
	"strings"
)

// Individual associates a candidate solution with its fitness score.
type Individual struct {
	Candidate interface{}
	Fitness   float64
}

// Population is a group of individual.
// TODO: check if and where we would benefit of having a slice of structs
// instead of pointers
type Population []*Individual

// Len is the number of elements in the collection.
func (s Population) Len() int { return len(s) }

// Less reports whether the element with
// index a should sort before the element with index b.
func (s Population) Less(i, j int) bool { return s[i].Fitness < s[j].Fitness }

// Swap swaps the elements with indexes i and j.
func (s Population) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func (s Population) String() string {
	reprs := make([]string, 0, len(s))
	for _, cand := range s {
		if cand != nil {
			reprs = append(reprs, fmt.Sprintf("%v|%v", cand.Candidate, cand.Fitness))
		} else {
			reprs = append(reprs, "nil")
		}
	}
	return fmt.Sprintf("{%s}", strings.Join(reprs, ", "))
}
