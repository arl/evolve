package main

import (
	"math"
	"testing"
)

func TestRouteEvaluator(t *testing.T) {
	//
	// a----------b
	// |           \
	// |            \
	// d_____________c
	//
	a := point{0, 20}  // cities[0]
	b := point{20, 20} // cities[1]
	c := point{30, 0}  // cities[2]
	d := point{0, 0}   // cities[3]
	cities := []point{a, b, c, d}

	e := newRouteEvaluator(cities)

	tests := []struct {
		a, b int
		want int
	}{
		{a: 0, b: 1, want: 20},
		{a: 0, b: 3, want: 20},
		{a: 2, b: 3, want: 30},
		{a: 1, b: 2, want: int(math.Sqrt(20*20 + 10*10))},
	}

	var tot int
	for _, tt := range tests {
		ab := e.dists[tt.a][tt.b]
		ba := e.dists[tt.b][tt.a]

		if ab != ba {
			t.Errorf("got dists[a][b] != dists[b][a] (%v and %v)", ab, ba)
		}

		if ab != tt.want {
			t.Errorf("got dists[a][b] == %v, want %v", ab, tt.want)
		}
		tot += ab
	}

	if !t.Failed() {
		got := int(e.Fitness([]int{0, 1, 2, 3}, nil))
		if got != tot {
			t.Errorf("got total distance = %v, want %v", got, tot)
		}
	}
}
