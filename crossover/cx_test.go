package crossover

import (
	"math/rand"
	"testing"

	"github.com/arl/evolve"
	"github.com/arl/evolve/generator"
	"github.com/google/go-cmp/cmp"
)

func TestCX(t *testing.T) {
	rng := rand.New(rand.NewSource(99))
	xover := evolve.Crossover[[]int]{
		Probability: generator.Const(1.0),
		Mater:       CX[int]{},
	}

	tests := []struct {
		name       string
		p0, p1     []int
		off0, off1 []int
	}{
		{
			name: "4-elem cycle",
			p0:   []int{1, 2, 3, 4, 5, 6, 7},
			p1:   []int{7, 5, 1, 3, 2, 6, 4},
			off0: []int{1, 5, 3, 4, 2, 6, 7},
			off1: []int{7, 2, 1, 3, 5, 6, 4},
		},
		{
			name: "1-elem cycle",
			p0:   []int{1, 2, 3},
			p1:   []int{1, 3, 2},
			off0: []int{1, 3, 2},
			off1: []int{1, 2, 3},
		},
		{
			name: "single cycle",
			p0:   []int{1, 2, 3},
			p1:   []int{2, 3, 1},
			off0: []int{1, 2, 3},
			off1: []int{2, 3, 1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pop := evolve.NewPopulationOf([][]int{tt.p0, tt.p1}, nil)
			xover.Apply(pop, rng)
			off0, off1 := pop.Candidates[0], pop.Candidates[1]

			if !(cmp.Equal(tt.off0, off0) && cmp.Equal(tt.off1, off1)) &&
				!(cmp.Equal(tt.off0, off1) && cmp.Equal(tt.off1, off0)) {
				t.Errorf("unexpected offsprings\n%+v\n%+v\n\nwant\n\n%+v\n%+v\n", off0, off1, tt.off0, tt.off1)
			}
		})
	}
}
