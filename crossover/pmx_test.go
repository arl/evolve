package crossover

import (
	"math/rand"
	"testing"

	"github.com/arl/evolve"
	"github.com/arl/evolve/generator"
	"github.com/google/go-cmp/cmp"
	"golang.org/x/exp/constraints"
	"golang.org/x/exp/slices"
)

func TestPMX(t *testing.T) {
	rng := rand.New(rand.NewSource(99))

	xover := evolve.Crossover[[]int]{
		Probability: generator.Const(1.0),
		Mater:       PMX[int]{},
	}

	items := [][]int{
		{1, 2, 3, 4, 5, 6, 7, 8},
		{3, 7, 5, 1, 6, 8, 2, 4},
	}

	pop := evolve.NewPopulationOf(items, nil)

	// Perform multiple crossovers to check different crossover points.
	for i := 0; i < 50; i++ {
		xover.Apply(pop, rng)

		for _, ind := range pop.Candidates {
			for j := 1; j <= 8; j++ {
				if !slices.Contains(ind, j) {
					t.Errorf("offspring is missing element %d", j)
				}
			}
		}
	}
}

func Test_mapBasedPMX(t *testing.T) {
	tests := []struct {
		p1, p2             []int
		xp1, xp2           int
		wantOff1, wantOff2 []int
	}{
		{
			p1:       []int{0, 1, 2, 3, 4, 5},
			p2:       []int{3, 4, 5, 0, 1, 2},
			xp1:      0,
			xp2:      2,
			wantOff1: []int{3, 4, 2, 0, 1, 5},
			wantOff2: []int{0, 1, 5, 3, 4, 2},
		},
		{
			p1:       []int{0, 1, 2, 3, 4, 5},
			p2:       []int{3, 4, 5, 0, 1, 2},
			xp1:      4,
			xp2:      1,
			wantOff1: []int{3, 4, 5, 0, 1, 2},
			wantOff2: []int{0, 1, 2, 3, 4, 5},
		},
	}

	for _, tt := range tests {
		off1 := make([]int, len(tt.p1))
		off2 := make([]int, len(tt.p2))
		copy(off1, tt.p1)
		copy(off2, tt.p2)

		mapBasedPMX(tt.p1, tt.p2, off1, off2, tt.xp1, tt.xp2)

		if !cmp.Equal(tt.wantOff1, off1) {
			t.Errorf("off1 = %+v, want %+v", off1, tt.wantOff1)
		}
		if !cmp.Equal(tt.wantOff2, off2) {
			t.Errorf("off2 = %+v, want %+v", off2, tt.wantOff2)
		}
	}
}

func benchmarkPMX(seqlen int) func(*testing.B) {
	const nxpts = 2

	pmx := PMX[int]{}
	p1 := seq[int](seqlen)
	p2 := seq[int](seqlen)

	return func(b *testing.B) {
		rng := rand.New(rand.NewSource(99))

		b.ReportAllocs()
		b.ResetTimer()
		for n := 0; n < b.N; n++ {
			off1, off2 := pmx.Mate(p1, p2, rng)
			p1 = off1
			p2 = off2
		}
	}
}

func BenchmarkPMX(b *testing.B) {
	b.Run("seqlen=62", benchmarkPMX(62))
}

// seq returns a slice containing the sequence of consecutive numbers from 0 to n.
func seq[T constraints.Integer | constraints.Float](n int) []T {
	s := make([]T, n)
	for i := 0; i < n; i++ {
		s[i] = T(i)
	}
	return s
}
