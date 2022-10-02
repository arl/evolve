package xover

import (
	"math/rand"
	"testing"

	"github.com/arl/evolve/generator"
	"github.com/stretchr/testify/assert"
)

func TestCX(t *testing.T) {
	rng := rand.New(rand.NewSource(99))

	xover := New[[]int](CX[int]{})
	xover.Probability = generator.Const(1.0)
	xover.Points = generator.Const(1) // unused

	tests := []struct {
		name       string
		p1, p2     []int
		off1, off2 []int
	}{
		{
			name: "4-elem cycle",
			p1:   []int{1, 2, 3, 4, 5, 6, 7},
			p2:   []int{7, 5, 1, 3, 2, 6, 4},
			off1: []int{1, 5, 3, 4, 2, 6, 7},
			off2: []int{7, 2, 1, 3, 5, 6, 4},
		},
		{
			name: "1-elem cycle",
			p1:   []int{1, 2, 3},
			p2:   []int{1, 3, 2},
			off1: []int{1, 3, 2},
			off2: []int{1, 2, 3},
		},
		{
			name: "single cycle",
			p1:   []int{1, 2, 3},
			p2:   []int{2, 3, 1},
			off1: []int{1, 2, 3},
			off2: []int{2, 3, 1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pop := [][]int{tt.p1, tt.p2}
			got := xover.Apply(pop, rng)
			assert.Equal(t, tt.off1, got[0], "unexpected offspring1")
			assert.Equal(t, tt.off2, got[1], "unexpected offspring2")
		})
	}
}

func TestCXDifferentLength(t *testing.T) {
	rng := rand.New(rand.NewSource(99))

	xover := New[[]int](CX[int]{})
	xover.Points = generator.Const(2)

	pop := make([][]int, 2)
	pop[0] = []int{1, 2, 3, 4, 5, 6, 7, 8}
	pop[1] = []int{3, 7, 5, 1}

	assert.Panics(t, func() {
		xover.Apply(pop, rng)
	})
}
