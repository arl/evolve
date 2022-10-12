package main

import (
	"fmt"
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
	a := tsp.Point2D{X: 0, Y: 20}  // cities[0]
	b := tsp.Point2D{X: 20, Y: 20} // cities[1]
	c := tsp.Point2D{X: 30, Y: 0}  // cities[2]
	d := tsp.Point2D{X: 0, Y: 0}   // cities[3]
	cities := []tsp.Point2D{a, b, c, d}

	e := newRouteEvaluator(cities)

	tests := []struct {
		a, b int
		want float64
	}{
		{a: 0, b: 1, want: 20},
		{a: 0, b: 3, want: 20},
		{a: 2, b: 3, want: 30},
		{a: 1, b: 2, want: math.Sqrt(20*20 + 10*10)},
	}

	var tot float64
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
		got := e.Fitness([]int{0, 1, 2, 3}, nil)
		if got != tot {
			t.Errorf("got total distance = %v, want %v", got, tot)
		}
	}
}

func TestBerlin52Opt(t *testing.T) {
	opt := []int{1, 49, 32, 45, 19, 41, 8, 9, 10, 43, 33, 51, 11, 52, 14, 13, 47, 26, 27, 28, 12, 25, 4, 6, 15, 5, 24, 48, 38, 37, 40, 39, 36, 35, 34, 44, 46, 16, 29, 50, 20, 23, 30, 2, 7, 42, 21, 17, 3, 18, 31, 22}
	for i := range opt {
		opt[i] = opt[i] - 1
	}

	e := newRouteEvaluator(berlin52)
	fmt.Println(e.Fitness(opt, nil))
}
